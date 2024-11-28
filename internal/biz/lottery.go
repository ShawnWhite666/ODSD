package biz

import (
	"context"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/pkg/middlewares/lock"
	"github.com/BitofferHub/pkg/middlewares/log"
	"strconv"
	"strings"
	"time"
)

type LotteryCase struct {
	prizeRepo     PrizeRepo
	couponRepo    CouponRepo
	blackUserRepo BlackUserRepo
	blackIpRepo   BlackIpRepo
	resultRepo    ResultRepo
	tm            Transaction
}

func NewLotteryCase(pr PrizeRepo, cr CouponRepo, bur BlackUserRepo,
	bir BlackIpRepo, result ResultRepo, tm Transaction) *LotteryCase {
	return &LotteryCase{
		prizeRepo:     pr,
		couponRepo:    cr,
		blackUserRepo: bur,
		blackIpRepo:   bir,
		resultRepo:    result,
		tm:            tm,
	}
}

func (l *LotteryCase) GetPrize(ctx context.Context, prizeCode int) (*LotteryPrize, error) {
	var prize *LotteryPrize
	lotteryPrizeList, err := l.GetAllUsefulPrizes(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|ToLotteryPrize:%v", err)
		return nil, err
	}
	for _, lotteryPrize := range lotteryPrizeList {
		if lotteryPrize.PrizeCodeLow <= prizeCode &&
			lotteryPrize.PrizeCodeHigh >= prizeCode {
			// 中奖编码区间满足条件，说明可以中奖
			if lotteryPrize.PrizeType < constant.PrizeTypeEntitySmall { // 如果非实物奖直接发，实物奖需要看是不是在黑名单外
				prize = lotteryPrize
				break
			}
		}
	}
	return prize, nil
}

// GetPrizeWithCache 获取中奖的奖品类型
func (l *LotteryCase) GetPrizeWithCache(ctx context.Context, prizeCode int) (*LotteryPrize, error) {
	var prize *LotteryPrize
	lotteryPrizeList, err := l.GetAllUsefulPrizesWithCache(ctx)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|ToLotteryPrize:%v", err)
		return nil, err
	}
	for _, lotteryPrize := range lotteryPrizeList {
		if lotteryPrize.PrizeCodeLow <= prizeCode &&
			lotteryPrize.PrizeCodeHigh >= prizeCode {
			// 中奖编码区间满足条件，说明可以中奖
			// if lotteryPrize.PrizeType < constant.PrizeTypeEntitySmall { //如果非实物奖直接发，实物奖需要看是不是在黑名单外
			prize = lotteryPrize
			break
			// }
		}
	}
	return prize, nil
}

// GiveOutPrize 发奖，奖品数量减1
func (l *LotteryCase) GiveOutPrize(ctx context.Context, prizeID int) (bool, error) {
	// 该类奖品的库存数量减1
	ok, err := l.prizeRepo.DecrLeftNum(prizeID, 1)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GiveOutPrize err:%v", err)
		return false, fmt.Errorf("LotteryCase|GiveOutPrize:%v", err)
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

// GiveOutPrizeWithCache 发奖，奖品数量减1,并且同步更新缓存
func (l *LotteryCase) GiveOutPrizeWithCache(ctx context.Context, prizeID int) (bool, error) {
	// 该类奖品的库存数量减1
	ok, err := l.prizeRepo.DecrLeftNum(prizeID, 1)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GiveOutPrize err:%v", err)
		return false, fmt.Errorf("LotteryCase|GiveOutPrize:%v", err)
	}
	if !ok {
		return false, nil
	}
	// 扣减库存成功
	if err = l.prizeRepo.UpdateByCache(&Prize{Id: uint(prizeID)}); err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GiveOutPrize|UpdateByCache err:%v", err)
		return false, fmt.Errorf("LotteryCase|GiveOutPrize|UpdateByCache:%v", err)
	}

	return true, nil
}

func (l *LotteryCase) GiveOutPrizeWithPool(ctx context.Context, prizeID int) (bool, error) {
	cnt, err := l.prizeRepo.DecrLeftNumByPool(prizeID)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GiveOutPrizeWithPool err:%v", err)
	}
	// 扣减完之后剩余奖品池中该奖品的数量小与0，所以当前时段该奖品不足了，不能发奖
	if cnt < 0 {
		return false, nil
	}
	// 奖品池成功之后再周数据库发奖逻辑
	return l.GiveOutPrize(ctx, prizeID)
}

// GetAllUsefulPrizes 获取所有可用奖品
func (l *LotteryCase) GetAllUsefulPrizes(ctx context.Context) ([]*LotteryPrize, error) {
	list, err := l.prizeRepo.GetAllUsefulPrizeList()
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GetAllUsefulPrizes:%v", err)
		return nil, fmt.Errorf("LotteryCase|GetAllUsefulPrizes:%v", err)
	}
	if len(list) == 0 {
		return nil, nil
	}
	lotteryPrizeList := make([]*LotteryPrize, 0)
	for _, prize := range list {
		codes := strings.Split(prize.PrizeCode, "-")
		if len(codes) == 2 {
			// 设置了获奖编码范围 a-b 才可以进行抽奖
			codeA := codes[0]
			codeB := codes[1]
			low, err1 := strconv.Atoi(codeA)
			high, err2 := strconv.Atoi(codeB)
			if err1 == nil && err2 == nil && high >= low && low >= 0 && high < constant.PrizeCodeMax {
				lotteryPrize := &LotteryPrize{
					Id:            prize.Id,
					Title:         prize.Title,
					PrizeNum:      prize.PrizeNum,
					LeftNum:       prize.LeftNum,
					PrizeCodeLow:  low,
					PrizeCodeHigh: high,
					Img:           prize.Img,
					DisplayOrder:  prize.DisplayOrder,
					PrizeType:     prize.PrizeType,
					PrizeProfile:  prize.PrizeProfile,
				}
				lotteryPrizeList = append(lotteryPrizeList, lotteryPrize)
			}
		}
	}
	return lotteryPrizeList, nil
}

func (l *LotteryCase) GetAllUsefulPrizesWithCache(ctx context.Context) ([]*LotteryPrize, error) {
	// 筛选出符合条件的奖品列表
	list, err := l.prizeRepo.GetAllUsefulPrizeListWithCache()
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GetAllUsefulPrizes:%v", err)
		return nil, fmt.Errorf("LotteryCase|GetAllUsefulPrizes:%v", err)
	}
	if len(list) == 0 {
		return nil, nil
	}
	// 对db的prize做一个类型转换，转化为LotteryPrize
	lotteryPrizeList := make([]*LotteryPrize, 0)
	for _, prize := range list {
		codes := strings.Split(prize.PrizeCode, "-")
		if len(codes) == 2 {
			// 设置了获奖编码范围 a-b 才可以进行抽奖
			codeA := codes[0]
			codeB := codes[1]
			low, err1 := strconv.Atoi(codeA)
			high, err2 := strconv.Atoi(codeB)
			if err1 == nil && err2 == nil && high >= low && low >= 0 && high < constant.PrizeCodeMax {
				lotteryPrize := &LotteryPrize{
					Id:            prize.Id,
					Title:         prize.Title,
					PrizeNum:      prize.PrizeNum,
					LeftNum:       prize.LeftNum,
					PrizeCodeLow:  low,
					PrizeCodeHigh: high,
					Img:           prize.Img,
					DisplayOrder:  prize.DisplayOrder,
					PrizeType:     prize.PrizeType,
					PrizeProfile:  prize.PrizeProfile,
				}
				lotteryPrizeList = append(lotteryPrizeList, lotteryPrize)
			}
		}
	}
	return lotteryPrizeList, nil
}

// PrizeCouponDiff 发放不同编码的优惠券
func (l *LotteryCase) PrizeCouponDiff(ctx context.Context, prizeID int) (string, error) {
	// 分布式锁保证查询和更新操作的原子性，并且保证每个连续操作串行执行
	// 因为需要更新数据的信息，所以要select，单纯用条件update只会返回受影响的记录数，不会返回具体信息，就拿不到优惠券的编码，所以需要两个操作，先select，再update
	key := fmt.Sprint(0 - prizeID - constant.CouponDiffLockLimit)
	lock1 := lock.NewRedisLock(key, lock.WithExpireSeconds(5), lock.WithWatchDogMode())
	if err := lock1.Lock(ctx); err != nil {
		log.ErrorContextf(ctx, "LotteryCase|PrizeCouponDiff:%v", err)
		return "", fmt.Errorf("LotteryCase|PrizeCouponDiff:%v", err)
	}
	defer lock1.Unlock(ctx)
	// 查询
	couponID := 0
	coupon, err := l.couponRepo.GetGetNextUsefulCoupon(prizeID, couponID)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|PrizeCouponDiff:%v\n", err)
		return "", err
	}
	if coupon == nil {
		//log.InfoContextf(ctx, "LotteryCase|PrizeCouponDiff: coupon is nil")
		return "", nil
	}
	// 更新
	coupon.SysStatus = 2
	if err := l.couponRepo.Update(coupon, "sys_status"); err != nil {
		log.ErrorContextf(ctx, "LotteryCase|PrizeCouponDiff:%v\n", err)
		return "", err
	}
	return coupon.Code, nil
}

// PrizeCouponDiffWithCache 带缓存的优惠券发奖，从缓存中拿出一个优惠券,要用缓存的话，需要项目启动的时候将优惠券导入到缓存
func (l *LotteryCase) PrizeCouponDiffWithCache(ctx context.Context, prizeID int) (string, error) {
	code, err := l.couponRepo.GetNextUsefulCouponFromCache(prizeID)
	if err != nil {
		return "", fmt.Errorf("LotteryCase|PrizeCouponDiffByCache:%v", err)
	}
	if code == "" {
		log.InfoContextf(ctx, "LotteryCase|PrizeCouponDiffByCache code is nil with prize_id=%d", prizeID)
		return "", nil
	}
	coupon := Coupon{
		Code:      code,
		SysStatus: 2,
	}
	if err = l.couponRepo.UpdateByCode(code, &coupon, "sys_status"); err != nil {
		return "", fmt.Errorf("LotteryCase|PrizeCouponDiffByCache:%v", err)
	}
	return code, nil
}

func (l *LotteryCase) PrizeLargeBlackLimit(ctx context.Context, blackUser *BlackUser,
	blackIp *BlackIp, lotteryUserInfo *LotteryUserInfo) error {
	now := time.Now()
	blackTime := constant.DefaultBlackTime
	// 用户黑明单限制
	if blackUser == nil || blackUser.UserId <= 0 {
		blackUserInfo := &BlackUser{
			UserId:    lotteryUserInfo.UserID,
			UserName:  lotteryUserInfo.UserName,
			BlackTime: now.Add(time.Second * time.Duration(blackTime)),
			// SysCreated: time.Time{},
			// SysUpdated: time.Time{},
			SysIp: lotteryUserInfo.IP,
		}
		if err := l.blackUserRepo.Create(blackUserInfo); err != nil {
			log.ErrorContextf(ctx, "LotteryCase|PrizeLargeBlackLimit:%v", err)
			return fmt.Errorf("LotteryCase|PrizeLargeBlackLimit:%v", err)
		}
	} else {
		blackUserInfo := &BlackUser{
			UserId:    lotteryUserInfo.UserID,
			BlackTime: now.Add(time.Second * time.Duration(blackTime)),
		}
		if err := l.blackUserRepo.Update(lotteryUserInfo.UserID, blackUserInfo, "black_time"); err != nil {
			log.ErrorContextf(ctx, "LotteryCase|PrizeLargeBlackLimit:%v", err)
			return fmt.Errorf("LotteryCase|PrizeLargeBlackLimit:%v", err)
		}
	}
	// ip黑明但限制
	if blackIp == nil || blackIp.Ip == "" {
		blackIPInfo := &BlackIp{
			Ip:        lotteryUserInfo.IP,
			BlackTime: now.Add(time.Second * time.Duration(blackTime)),
			// SysCreated: time.Time{},
			// SysUpdated: time.Time{},
		}
		if err := l.blackIpRepo.Create(blackIPInfo); err != nil {
			log.ErrorContextf(ctx, "LotteryCase|PrizeLargeBlackLimit:%v", err)
			return fmt.Errorf("LotteryCase|PrizeLargeBlackLimit:%v", err)
		}
	} else {
		blackIPInfo := &BlackIp{
			Ip:        lotteryUserInfo.IP,
			BlackTime: now.Add(time.Second * time.Duration(blackTime)),
			// SysUpdated: time.Time{},
		}
		if err := l.blackIpRepo.Update(lotteryUserInfo.IP, blackIPInfo, "black_time"); err != nil {
			log.ErrorContextf(ctx, "LotteryCase|PrizeLargeBlackLimit:%v", err)
			return fmt.Errorf("LotteryCase|PrizeLargeBlackLimit:%v", err)
		}
	}
	return nil
}

func (l *LotteryCase) GetPrizeNumWithPool(ctx context.Context, prizeID uint) (int, error) {

	num, err := l.prizeRepo.GetPrizePoolNum(prizeID)
	if err != nil {
		log.ErrorContextf(ctx, "LotteryCase|GetPrizeNumWithPool err: %v", err)
		return 0, fmt.Errorf("LotteryCase|GetPrizeNumWithPool:%v", err)
	}
	return num, nil
}

func (l *LotteryCase) LotteryResult(ctx context.Context, prize *LotteryPrize, uid uint, userName, ip string, prizeCode int) error {
	result := Result{
		PrizeId:   prize.Id,
		PrizeName: prize.Title,
		PrizeType: prize.PrizeType,
		UserId:    uid,
		UserName:  userName,
		PrizeCode: uint(prizeCode),
		PrizeData: prize.PrizeProfile,
		// SysCreated: time.Now(),
		SysIp:     ip,
		SysStatus: 1,
	}

	if err := l.resultRepo.Create(&result); err != nil {
		log.ErrorContextf(ctx, "resultService|LotteryResult:%v", err)
		return fmt.Errorf("resultService|LotteryResult:%v", err)
	}
	return nil
}
