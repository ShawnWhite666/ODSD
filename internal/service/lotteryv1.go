package service

import (
	"context"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/lotterysvr/internal/utils"
	"github.com/BitofferHub/pkg/middlewares/lock"
	"github.com/BitofferHub/pkg/middlewares/log"
	pb "github.com/BitofferHub/proto_center/api/lotterysvr/v1"
)

func (l *LotteryService) LotteryV1(ctx context.Context, req *pb.LotteryReq) (*pb.LotteryRsp, error) {
	rsp := &pb.LotteryRsp{
		CommonRsp: &pb.CommonRspInfo{
			Code:   int32(Success),
			Msg:    GetErrMsg(Success),
			UserId: req.UserId,
		},
	}
	defer func() {
		// 通过对应的Code，获取Msg
		rsp.CommonRsp.Msg = GetErrMsg(ErrCode(rsp.CommonRsp.Code))
	}()
	var (
		ok  bool
		err error
	)
	// 1. 根据token解析出用户信息
	//jwtClaims, err := utils.ParseJwtToken(req.Token, constant.SecretKey)
	//if err != nil || jwtClaims == nil {
	//	rsp.CommonRsp.Code = int32(ErrJwtParse)
	//	log.Errorf("jwt parse err, token=%s,user_id=%s\n", req.Token, req.UserId)
	//	return nil, fmt.Errorf("LotteryV1|jwt parse err")
	//}
	//log.Infof("LotteryV1|req====%+v", req)
	userID := uint(req.UserId)
	log.Infof("LotteryV1|user_id=%d", userID)
	lockKey := fmt.Sprintf(constant.LotteryLockKeyPrefix+"%d", userID)
	lock1 := lock.NewRedisLock(lockKey, lock.WithExpireSeconds(5), lock.WithWatchDogMode())

	// 1. 用户抽奖分布式锁定,防止同一个用户同一时间抽奖抽奖多次
	if err := lock1.Lock(ctx); err != nil {
		rsp.CommonRsp.Code = int32(ErrInternalServer)
		log.ErrorContextf(ctx, "LotteryHandler|Process:%v", err)
		return nil, fmt.Errorf("LotteryV1|lock err")
	}
	defer lock1.Unlock(ctx)
	// 2. 验证用户今日抽奖次数
	ok, err = l.limitCase.CheckUserDayLotteryTimes(ctx, userID)
	if err != nil {
		rsp.CommonRsp.Code = int32(ErrInternalServer)
		log.ErrorContextf(ctx, "LotteryHandler|CheckUserDayLotteryTimes:%v", err)
		return nil, fmt.Errorf("LotteryV1|CheckUserDayLotteryTimes err")
	}
	if !ok {
		rsp.CommonRsp.Code = int32(ErrUserLimitInvalid)
		//log.InfoContextf(ctx, "LotteryHandler|CheckUserDayLotteryTimes:%v", err)
		return rsp, nil
	}

	// 3. 验证当天IP参与的抽奖次数
	ipDayLotteryTimes := l.limitCase.CheckIPLimit(ctx, req.Ip)
	if ipDayLotteryTimes > constant.IpLimitMax {
		rsp.CommonRsp.Code = int32(ErrIPLimitInvalid)
		//log.InfoContextf(ctx, "LotteryHandler|CheckUserDayLotteryTimes:%v", err)
		return rsp, nil
	}

	// 4. 验证IP是否在ip黑名单
	ok, blackIpInfo, err := l.limitCase.CheckBlackIP(ctx, req.Ip)
	if err != nil {
		rsp.CommonRsp.Code = int32(ErrInternalServer)
		log.ErrorContextf(ctx, "LotteryHandler|CheckBlackIP:%v", err)
		return nil, fmt.Errorf("LotteryV1|CheckBlackIP err")
	}
	// ip黑明单生效
	if !ok {
		rsp.CommonRsp.Code = int32(ErrBlackedIP)
		//log.InfoContextf(ctx, "LotteryHandler|CheckBlackIP blackIpInfo is %+v\n", blackIpInfo)
		return rsp, nil
	}

	// 5. 验证用户是否在黑明单中
	ok, blackUserInfo, err := l.limitCase.CheckBlackUser(ctx, userID)
	if err != nil {
		rsp.CommonRsp.Code = int32(ErrInternalServer)
		log.ErrorContextf(ctx, "LotteryHandler|CheckBlackUser:%v", err)
		return nil, fmt.Errorf("LotteryV1|CheckBlackIP err")
	}
	// 用户黑明单生效
	if !ok {
		rsp.CommonRsp.Code = int32(ErrBlackedUser)
		log.ErrorContextf(ctx, "LotteryHandler|CheckBlackUser blackUserInfo is %v\n", blackUserInfo)
		return rsp, nil
	}

	// 6. 中奖逻辑实现
	prizeCode := utils.Random(constant.PrizeCodeMax)
	log.InfoContextf(ctx, "LotteryHandlerV1|prizeCode=%d\n", prizeCode)
	prize, err := l.lotteryCase.GetPrize(ctx, prizeCode)
	if err != nil {
		rsp.CommonRsp.Code = int32(ErrInternalServer)
		log.ErrorContextf(ctx, "LotteryHandler|CheckBlackUser:%v", err)
		return nil, fmt.Errorf("LotteryV1|GetPrize err")
	}
	if prize == nil || prize.PrizeNum < 0 || (prize.PrizeNum > 0 && prize.LeftNum <= 0) {
		rsp.CommonRsp.Code = int32(ErrNotWon)
		return rsp, nil
	}

	// 7. 有剩余奖品发放
	if prize.PrizeNum > 0 {
		ok, err = l.lotteryCase.GiveOutPrize(ctx, int(prize.Id))
		if err != nil {
			rsp.CommonRsp.Code = int32(ErrInternalServer)
			log.ErrorContextf(ctx, "LotteryHandler|GiveOutPrize:%v", err)
			return nil, fmt.Errorf("LotteryV1|GiveOutPrize err")
		}
		// 奖品不足，发放失败
		if !ok {
			rsp.CommonRsp.Code = int32(ErrPrizeNotEnough)
			//log.InfoContextf(ctx, "LotteryHandler|GiveOutPrize:%v", err)
			return rsp, nil
		}
	}

	/***如果中奖记录重要的的话，可以考虑用事务将下面逻辑包裹*****/
	// 8. 发优惠券
	if prize.PrizeType == constant.PrizeTypeCouponDiff {
		code, err := l.lotteryCase.PrizeCouponDiff(ctx, int(prize.Id))
		if err != nil {
			rsp.CommonRsp.Code = int32(ErrInternalServer)
			//log.InfoContextf(ctx, "LotteryHandler|PrizeCouponDiff:%v", err)
			return nil, fmt.Errorf("LotteryV1|PrizeCouponDiff err")
		}
		if code == "" {
			rsp.CommonRsp.Code = int32(ErrNotWon)
			//log.InfoContextf(ctx, "LotteryHandler|PrizeCouponDiff coupon left is nil")
			return rsp, nil
		}
		prize.CouponCode = code
	}
	rsp.PrizeInfo = &pb.LotteryPrizeInfo{
		Id:            uint32(prize.Id),
		Title:         prize.Title,
		PrizeNum:      int32(prize.PrizeNum),
		LeftNum:       int32(prize.LeftNum),
		PrizeCodeLow:  int32(prize.PrizeCodeLow),
		PrizeCodeHigh: int32(prize.PrizeCodeHigh),
		Img:           prize.Img,
		DisplayOrder:  uint32(prize.DisplayOrder),
		PrizeType:     uint32(prize.PrizeType),
		PrizeProfile:  prize.PrizeProfile,
		CouponCode:    prize.CouponCode,
	}

	// 9 记录中奖纪录
	if err := l.lotteryCase.LotteryResult(ctx, prize, userID, req.UserName, req.Ip, prizeCode); err != nil {
		rsp.CommonRsp.Code = int32(ErrInternalServer)
		log.ErrorContextf(ctx, "LotteryHandler|PrizeCouponDiff:%v", err)
		return nil, fmt.Errorf("LotteryV1|LotteryResult err")
	}

	// 10. 如果中了实物大奖，需要把ip和用户置于黑明单中一段时间，防止同一个用户频繁中大奖
	if prize.PrizeType == constant.PrizeTypeEntityLarge {
		lotteryUserInfo := biz.LotteryUserInfo{
			UserID:   userID,
			UserName: req.UserName,
			IP:       req.Ip,
		}
		log.InfoContextf(ctx, "LotteryV1|user_id=%d", userID)
		if err := l.lotteryCase.PrizeLargeBlackLimit(ctx, blackUserInfo, blackIpInfo, &lotteryUserInfo); err != nil {
			rsp.CommonRsp.Code = int32(ErrInternalServer)
			log.ErrorContextf(ctx, "LotteryHandler|PrizeLargeBlackLimit:%v", err)
			return nil, fmt.Errorf("LotteryV1|PrizeLargeBlackLimit err")
		}
	}

	return rsp, nil
}
