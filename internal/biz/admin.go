package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/lotterysvr/internal/utils"
	"github.com/BitofferHub/pkg/middlewares/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type AdminCase struct {
	couponRepo       CouponRepo
	prizeRepo        PrizeRepo
	lotteryTimesRepo LotteryTimesRepo
	resultRepo       ResultRepo
}

func NewAdminCase(pr PrizeRepo, cr CouponRepo, lr LotteryTimesRepo, rp ResultRepo) *AdminCase {
	return &AdminCase{
		couponRepo:       cr,
		prizeRepo:        pr,
		lotteryTimesRepo: lr,
		resultRepo:       rp,
	}
}

// GetPrizeList 获取db奖品列表
func (a *AdminCase) GetPrizeList(ctx context.Context) ([]*Prize, error) {
	//log.InfoContextf(ctx, "GetPrizeList!!!!!")
	list, err := a.prizeRepo.GetAll()
	if err != nil {
		log.ErrorContextf(ctx, "prizeCase|GetPrizeList err:%v", err)
		return nil, fmt.Errorf("prizeCase|GetPrizeList: %v", err)
	}
	return list, nil
}

// GetPrizeListWithCache 获取db奖品列表
func (a *AdminCase) GetPrizeListWithCache(ctx context.Context) ([]*Prize, error) {
	list, err := a.prizeRepo.GetAllWithCache()
	if err != nil {
		log.ErrorContextf(ctx, "prizeCase|GetPrizeList err:%v", err)
		return nil, fmt.Errorf("prizeCase|GetPrizeList: %v", err)
	}
	return list, nil
}

// GetViewPrizeList 获取奖品列表,这个方法用于管理后台使用，因为管理后台不需要高性能，所以不走缓存
func (a *AdminCase) GetViewPrizeList(ctx context.Context) ([]*ViewPrize, error) {
	list, err := a.prizeRepo.GetAll()
	if err != nil {
		log.ErrorContextf(ctx, "prizeCase|GetPrizeList err:%v", err)
		return nil, fmt.Errorf("prizeCase|GetPrizeList: %v", err)
	}
	prizeList := make([]*ViewPrize, 0)
	for _, prize := range list {
		if prize.SysStatus != constant.PrizeStatusNormal {
			continue
		}
		num, err := a.prizeRepo.GetPrizePoolNum(prize.Id)
		if err != nil {
			return nil, fmt.Errorf("prizeCase|GetPrizeList: %v", err)
		}
		title := fmt.Sprintf("【%d】%s", num, prize.Title)
		prizeList = append(prizeList, &ViewPrize{
			Id:        prize.Id,
			Title:     title,
			Img:       prize.Img,
			PrizeNum:  prize.PrizeNum,
			LeftNum:   prize.LeftNum,
			PrizeType: prize.PrizeType,
		})

	}
	return prizeList, nil
}

// GetViewPrizeListWithCache 获取奖品列表,优先从缓存获取
func (a *AdminCase) GetViewPrizeListWithCache(ctx context.Context) ([]*ViewPrize, error) {
	list, err := a.prizeRepo.GetAllWithCache()
	if err != nil {
		log.ErrorContextf(ctx, "prizeCase|GetPrizeList err:%v", err)
		return nil, fmt.Errorf("prizeCase|GetPrizeList: %v", err)
	}
	prizeList := make([]*ViewPrize, 0)
	for _, prize := range list {
		if prize.SysStatus != constant.PrizeStatusNormal {
			continue
		}
		prizeList = append(prizeList, &ViewPrize{
			Id:        prize.Id,
			Title:     prize.Title,
			Img:       prize.Img,
			PrizeNum:  prize.PrizeNum,
			LeftNum:   prize.LeftNum,
			PrizeType: prize.PrizeType,
		})
	}
	return prizeList, nil
}

// GetPrize 获取某个奖品
func (a *AdminCase) GetPrize(ctx context.Context, id uint) (*ViewPrize, error) {
	prizeModel, err := a.prizeRepo.Get(id)
	if err != nil {
		log.ErrorContextf(ctx, "prizeCase|GetPrize:%v", err)
		return nil, fmt.Errorf("prizeCase|GetPrize:%v", err)
	}
	prize := &ViewPrize{
		Id:        prizeModel.Id,
		Title:     prizeModel.Title,
		Img:       prizeModel.Img,
		PrizeNum:  prizeModel.PrizeNum,
		LeftNum:   prizeModel.LeftNum,
		PrizeType: prizeModel.PrizeType,
	}
	return prize, nil
}

// AddPrize 新增奖品
func (a *AdminCase) AddPrize(ctx context.Context, viewPrize *ViewPrize) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("AddPrize panic%v\n", err)
		}
	}()
	prize := Prize{
		Title:        viewPrize.Title,
		PrizeNum:     viewPrize.PrizeNum,
		LeftNum:      viewPrize.PrizeNum,
		PrizeCode:    viewPrize.PrizeCode,
		PrizeTime:    viewPrize.PrizeTime,
		Img:          viewPrize.Img,
		DisplayOrder: viewPrize.DisplayOrder,
		PrizeType:    viewPrize.PrizeType,
		BeginTime:    viewPrize.BeginTime,
		EndTime:      viewPrize.EndTime,
		PrizePlan:    viewPrize.PrizePlan,
		SysStatus:    1,
	}
	// 因为奖品是全量string缓存，新增奖品之后缓存有变动，所有要更新
	if err := a.prizeRepo.Create(&prize); err != nil {
		log.ErrorContextf(ctx, "adminCase|AddPrize err:%v", err)
		return fmt.Errorf("adminCase|AddPrize:%v", err)
	}
	return nil
}

// AddPrizeList 新增奖品列表
func (a *AdminCase) AddPrizeList(ctx context.Context, viewPrizeList []*ViewPrize) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("AddPrize panic%v\n", err)
		}
	}()
	prizeList := make([]Prize, 0)
	for _, viewPrize := range viewPrizeList {
		prize := Prize{
			Title:        viewPrize.Title,
			PrizeNum:     viewPrize.PrizeNum,
			LeftNum:      viewPrize.PrizeNum,
			PrizeCode:    viewPrize.PrizeCode,
			PrizeTime:    viewPrize.PrizeTime,
			Img:          viewPrize.Img,
			DisplayOrder: viewPrize.DisplayOrder,
			PrizeType:    viewPrize.PrizeType,
			BeginTime:    viewPrize.BeginTime,
			EndTime:      viewPrize.EndTime,
			PrizePlan:    viewPrize.PrizePlan,
			SysStatus:    1,
		}
		prizeList = append(prizeList, prize)
	}
	// 因为奖品是全量string缓存，新增奖品之后缓存有变动，所有要更新
	if err := a.prizeRepo.CreateInBatches(prizeList); err != nil {
		log.ErrorContextf(ctx, "adminCase|AddPrizeList err:%v", err)
		return fmt.Errorf("adminCase|AddPrizeList:%v", err)
	}
	return nil
}

func (a *AdminCase) ClearPrize(ctx context.Context) error {
	if err := a.prizeRepo.DeleteAll(); err != nil {
		log.ErrorContextf(ctx, "adminCase|ClearPrize err:%v", err)
		return fmt.Errorf("adminCase|ClearPrize:%v", err)
	}
	return nil
}

func (a *AdminCase) ClearCoupon(ctx context.Context) error {
	if err := a.couponRepo.DeleteAllWithCache(); err != nil {
		log.ErrorContextf(ctx, "adminCase|ClearCoupon err:%v", err)
		return fmt.Errorf("adminCase|ClearCoupon:%v", err)
	}
	return nil
}

// AddPrizeWithPool 带奖品池的新增奖品实现
func (a *AdminCase) AddPrizeWithPool(ctx context.Context, viewPrize *ViewPrize) error {
	prize := Prize{
		Title:        viewPrize.Title,
		PrizeNum:     viewPrize.PrizeNum,
		LeftNum:      viewPrize.PrizeNum,
		PrizeCode:    viewPrize.PrizeCode,
		PrizeTime:    viewPrize.PrizeTime,
		Img:          viewPrize.Img,
		DisplayOrder: viewPrize.DisplayOrder,
		PrizeType:    viewPrize.PrizeType,
		BeginTime:    viewPrize.BeginTime,
		EndTime:      viewPrize.EndTime,
		PrizePlan:    viewPrize.PrizePlan,
		SysStatus:    1,
		//SysUpdated:   time.Now(),
	}
	// 因为奖品是全量string缓存，新增奖品之后缓存有变动，所有要更新
	if err := a.prizeRepo.CreateWithCache(&prize); err != nil {
		log.Errorf("adminCase|AddPrize err:%v", err)
		return fmt.Errorf("adminCase|AddPrize:%v", err)
	}
	if err := a.ResetPrizePlan(ctx, &prize); err != nil {
		log.Errorf("adminCase|AddPrize ResetPrizePlan prize err:%v", err)
		return fmt.Errorf("adminCase|AddPrize ResetPrizePlan prize err:%v", err)
	}
	return nil
}

// AddPrizeWithCache 带缓存优化的新增奖品
func (a *AdminCase) AddPrizeWithCache(ctx context.Context, viewPrize *ViewPrize) error {
	prize := Prize{
		Title:        viewPrize.Title,
		PrizeNum:     viewPrize.PrizeNum,
		LeftNum:      viewPrize.PrizeNum,
		PrizeCode:    viewPrize.PrizeCode,
		PrizeTime:    viewPrize.PrizeTime,
		Img:          viewPrize.Img,
		DisplayOrder: viewPrize.DisplayOrder,
		PrizeType:    viewPrize.PrizeType,
		BeginTime:    viewPrize.BeginTime,
		EndTime:      viewPrize.EndTime,
		PrizePlan:    viewPrize.PrizePlan,
		SysStatus:    1,
		//SysUpdated:   time.Now(),
	}
	// 因为奖品是全量string缓存，新增奖品之后缓存有变动，所有要更新
	if err := a.prizeRepo.CreateWithCache(&prize); err != nil {
		log.ErrorContextf(ctx, "adminCase|AddPrize err:%v", err)
		return fmt.Errorf("adminCase|AddPrize:%v", err)
	}
	return nil
}

func (a *AdminCase) UpdateDbPrizeWithCache(ctx context.Context, prize *Prize, cols ...string) error {
	if err := a.prizeRepo.UpdateWithCache(prize, cols...); err != nil {
		log.ErrorContextf(ctx, "UpdateDbPrizeWithCache|%v", err)
		return fmt.Errorf("UpdateDbPrizeWithCache|%v", err)
	}
	return nil
}

func (a *AdminCase) UpdateDbPrize(ctx context.Context, db *gorm.DB, prize *Prize, cols ...string) error {
	if err := a.prizeRepo.Update(prize, cols...); err != nil {
		log.ErrorContextf(ctx, "UpdateDbPrize|%v", err)
		return fmt.Errorf("UpdateDbPrize|%v", err)
	}
	return nil
}

func (a *AdminCase) UpdatePrize(ctx context.Context, viewPrize *ViewPrize) error {
	if viewPrize == nil || viewPrize.Id <= 0 {
		log.Errorf("adminCase|UpdatePrize invalid prize err:%v", viewPrize)
		return fmt.Errorf("adminCase|UpdatePrize invalid prize")
	}
	prize := Prize{
		Title:        viewPrize.Title,
		PrizeNum:     viewPrize.PrizeNum,
		LeftNum:      viewPrize.LeftNum,
		PrizeCode:    viewPrize.PrizeCode,
		PrizeTime:    viewPrize.PrizeTime,
		Img:          viewPrize.Img,
		DisplayOrder: viewPrize.DisplayOrder,
		PrizeType:    viewPrize.PrizeType,
		BeginTime:    viewPrize.BeginTime,
		EndTime:      viewPrize.EndTime,
		PrizePlan:    viewPrize.PrizePlan,
		SysStatus:    viewPrize.SysStatus,
	}
	oldPrize, err := a.prizeRepo.Get(viewPrize.Id)
	if err != nil {
		log.ErrorContextf(ctx, "adminCase|UpdatePrize get old prize err:%v", err)
		return fmt.Errorf("adminCase|UpdatePrize:%v", err)
	}
	if oldPrize == nil {
		log.ErrorContextf(ctx, "adminCase|UpdatePrize prize not exists with id: %d", viewPrize.Id)
		return fmt.Errorf("adminCase|UpdatePrize prize not exists with id: %d", viewPrize.Id)
	}
	// 奖品数量发生了改变
	if prize.PrizeNum != oldPrize.PrizeNum {
		if prize.PrizeNum <= 0 {
			prize.PrizeNum = 0
		}
		if prize.LeftNum <= 0 {
			prize.LeftNum = 0
		}
	}
	if a.prizeRepo.Update(&prize, "title", "prize_num", "left_num", "prize_code", "prize_time", "img",
		"display_order", "prize_type", "begin_time", "end_time", "prize_plan"); err != nil {
		log.Errorf("adminCase|UpdatePrize Update prize err:%v", err)
		return fmt.Errorf("adminCase|UpdatePrize Update prize:%v", err)
	}
	return nil
}

func (a *AdminCase) UpdatePrizeWithPool(ctx context.Context, viewPrize *ViewPrize) error {
	if viewPrize == nil || viewPrize.Id <= 0 {
		log.Errorf("adminCase|UpdatePrize invalid prize err:%v", viewPrize)
		return fmt.Errorf("adminCase|UpdatePrize invalid prize")
	}
	prize := Prize{
		Title:        viewPrize.Title,
		PrizeNum:     viewPrize.PrizeNum,
		LeftNum:      viewPrize.LeftNum,
		PrizeCode:    viewPrize.PrizeCode,
		PrizeTime:    viewPrize.PrizeTime,
		Img:          viewPrize.Img,
		DisplayOrder: viewPrize.DisplayOrder,
		PrizeType:    viewPrize.PrizeType,
		BeginTime:    viewPrize.BeginTime,
		EndTime:      viewPrize.EndTime,
		PrizePlan:    viewPrize.PrizePlan,
		SysStatus:    viewPrize.SysStatus,
	}
	oldPrize, err := a.prizeRepo.Get(viewPrize.Id)
	if err != nil {
		log.Errorf("adminCase|UpdatePrize get old prize err:%v", err)
		return fmt.Errorf("adminCase|UpdatePrize:%v", err)
	}
	if oldPrize == nil {
		log.Errorf("adminCase|UpdatePrize prize not exists with id: %d", viewPrize.Id)
		return fmt.Errorf("adminCase|UpdatePrize prize not exists with id: %d", viewPrize.Id)
	}
	// 奖品数量发生了改变
	if prize.PrizeNum != oldPrize.PrizeNum {
		if prize.PrizeNum <= 0 {
			prize.PrizeNum = 0
		}
		if prize.LeftNum <= 0 {
			prize.LeftNum = 0
		}
		if err := a.ResetPrizePlan(ctx, &prize); err != nil {
			log.Errorf("adminCase|UpdatePrize ResetPrizePlan prize err:%v", err)
			return fmt.Errorf("adminCase|UpdatePrize ResetPrizePlan prize err:%v", err)
		}
	}
	if prize.PrizeTime != oldPrize.PrizeTime {
		if err := a.ResetPrizePlan(ctx, &prize); err != nil {
			log.Errorf("adminCase|UpdatePrize ResetPrizePlan prize err:%v", err)
			return fmt.Errorf("adminCase|UpdatePrize ResetPrizePlan prize err:%v", err)
		}
	}
	if a.prizeRepo.Update(&prize, "title", "prize_num", "left_num", "prize_code", "prize_time", "img",
		"display_order", "prize_type", "begin_time", "end_time", "prize_plan"); err != nil {
		log.Errorf("adminCase|UpdatePrize Update prize err:%v", err)
		return fmt.Errorf("adminCase|UpdatePrize Update prize:%v", err)
	}
	return nil
}

// GetCouponList 获取优惠券列表,库存优惠券数量和缓存优惠券数量，当这两个数量不一致的时候，需要重置缓存优惠券数量
func (a *AdminCase) GetCouponList(ctx context.Context, prizeID uint) ([]*ViewCouponInfo, int64, int64, error) {
	var (
		viewCouponList []*ViewCouponInfo
		couponList     []*Coupon
		err            error
		dbNum          int64
		cacheNum       int64
	)
	if prizeID > 0 {
		couponList, err = a.couponRepo.GetCouponListByPrizeID(prizeID)
		if err != nil {
			log.ErrorContextf(ctx, "AdminCase|GetCouponListByPrizeID|%v", err)
			return nil, 0, 0, fmt.Errorf("adminCase|GetCouponList invalid prize_id:%d", prizeID)
		}
		dbNum, cacheNum, err = a.couponRepo.GetCacheCouponNum(prizeID)
		if err != nil {
			log.ErrorContextf(ctx, "AdminCase|GetCacheCouponNum|%v", err)
			return nil, 0, 0, fmt.Errorf("adminCase|GetCouponList invalid prize_id:%d", prizeID)
		}
	} else {
		couponList, err = a.couponRepo.GetAll()
		if err != nil {
			log.ErrorContextf(ctx, "AdminCase|couponRepo.GetAll|%v", err)
			return nil, 0, 0, fmt.Errorf("adminCase|GetCouponList invalid prize_id:%d", prizeID)
		}
	}
	for _, coupon := range couponList {
		viewCouponList = append(viewCouponList, &ViewCouponInfo{
			Id:      coupon.Id,
			PrizeId: coupon.PrizeId,
			Code:    coupon.Code,
			//SysCreated: coupon.SysCreated,
			//SysUpdated: coupon.SysUpdated,
			SysStatus: coupon.SysStatus,
		})
	}
	return viewCouponList, dbNum, cacheNum, nil
}

// ImportCoupon 导入优惠券
func (a *AdminCase) ImportCoupon(ctx context.Context, prizeID uint, codes string) (int, int, error) {
	if prizeID <= 0 {
		return 0, 0, fmt.Errorf("adminCase|ImportCoupon invalid prizeID:%d", prizeID)
	}
	prize, err := a.prizeRepo.Get(prizeID)
	if err != nil {
		return 0, 0, fmt.Errorf("adminCase|ImportCoupon invalid prizeID:%d", prizeID)
	}
	if prize == nil || prize.PrizeType != constant.PrizeTypeCouponDiff {
		//log.InfoContextf(ctx, "adminCase|ImportCoupon invalid prize type:%d with prize_id %d", prize.PrizeType, prizeID)
		return 0, 0, fmt.Errorf("adminCase|ImportCoupon prize_type is not coupon with prize_id %d", prizeID)
	}
	var (
		successNum int
		failNum    int
	)
	codeList := strings.Split(codes, "\n")
	for _, code := range codeList {
		code = strings.TrimSpace(code)
		coupon := &Coupon{
			PrizeId: prizeID,
			Code:    code,
			//SysCreated: time.Now(),
			SysStatus: 1,
		}
		if err = a.couponRepo.Create(coupon); err != nil {
			failNum++
		} else {
			successNum++
		}
	}
	return successNum, failNum, nil
}

// ImportCouponWithCache 导入优惠券
func (a *AdminCase) ImportCouponWithCache(ctx context.Context, prizeID uint, codes string) (int, int, error) {
	if prizeID <= 0 {
		return 0, 0, fmt.Errorf("adminCase|ImportCoupon invalid prizeID:%d", prizeID)
	}
	prize, err := a.prizeRepo.GetWithCache(prizeID)
	if err != nil {
		return 0, 0, fmt.Errorf("adminCase|ImportCoupon invalid prizeID:%d", prizeID)
	}
	if prize == nil || prize.PrizeType != constant.PrizeTypeCouponDiff {
		log.InfoContextf(ctx, "adminCase|ImportCoupon invalid prize type:%d with prize_id %d", prize.PrizeType, prizeID)
		return 0, 0, fmt.Errorf("adminCase|ImportCoupon prize_type is not coupon with prize_id %d", prizeID)
	}
	var (
		successNum int
		failNum    int
	)
	codeList := strings.Split(codes, "\n")
	for _, code := range codeList {
		code = strings.TrimSpace(code)
		coupon := &Coupon{
			PrizeId: prizeID,
			Code:    code,
			//SysCreated: time.Now(),
			SysStatus: 1,
		}
		if err = a.couponRepo.Create(coupon); err != nil {
			failNum++
		} else {
			// db导入成功之后，再导入缓存
			ok, err := a.couponRepo.ImportCacheCoupon(prizeID, code)
			if err != nil {
				return 0, 0, fmt.Errorf("adminCase|ImportCoupon prize_type is not coupon with prize_id %d", prizeID)
			}
			if !ok {
				failNum++
			} else {
				successNum++
			}
		}
	}
	return successNum, failNum, nil
}

// ReCacheCoupon 根据数据库重置某种奖品的优惠券数据到缓存中
func (a *AdminCase) ReCacheCoupon(ctx context.Context, prizeID uint) (int64, int64, error) {
	if prizeID <= 0 {
		return 0, 0, fmt.Errorf("adminCase|ReCacheCoupon invalid prizeID:%d", prizeID)
	}
	successNum, failureNum, err := a.couponRepo.ReSetCacheCoupon(prizeID)
	if err != nil {
		log.ErrorContextf(ctx, "AdminCase|ReSetCacheCoupon|%v", err)
		return 0, 0, fmt.Errorf("adminCase|ReCacheCoupon:%v", err)
	}
	return successNum, failureNum, nil
}

// ResetPrizePlan 重置某种奖品的发奖计划
func (a *AdminCase) ResetPrizePlan(ctx context.Context, prize *Prize) error {
	if prize == nil || prize.Id < 1 {
		return fmt.Errorf("limitCase|ResetGiftPrizePlan invalid prize")
	}
	now := time.Now()
	// 奖品状态不对，不能发奖
	if prize.SysStatus == 2 ||
		prize.BeginTime.After(now) || // 还未开始
		prize.EndTime.Before(now) || // 已经结束
		prize.LeftNum <= 0 ||
		prize.PrizeNum <= 0 {
		if prize.PrizePlan != "" {
			// 在重置的时候，如果发现原来奖品的发奖计划不为空，需要清空发奖计划
			a.clearPrizePlan(ctx, prize)
		}
		//log.InfoContextf(ctx, "prize can not be given out")
		return nil
	}
	// PrizeTime, 发奖周期，这类奖品需要在多少天内发完
	prizePlanDays := int(prize.PrizeTime)
	if prizePlanDays <= 0 {
		a.setPrizePool(ctx, prize.Id, prize.LeftNum)
		//log.InfoContextf(ctx, "adminCase|ResetGiftPrizePlan|prizePlanDays <= 0")
		return nil
	}
	// 对于设置发奖周期的奖品重新计算出来合适的奖品发放节奏
	// 奖品池的剩余数先设置为空
	a.setPrizePool(ctx, prize.Id, 0)
	// 发奖周期中的每天的发奖概率一样，一天内24小时，每个小时的概率是不一样的，每个小时内的每一分钟的概率一样
	prizeNum := prize.PrizeNum
	// 先计算每天至少发多少奖
	avgPrizeNum := prizeNum / prizePlanDays

	// 每天可以分配到的奖品数量
	dayPrizeNumMap := make(map[int]int)
	// 发奖周期大雨1天，并且平均每天发的奖品书大于等于1
	if prizePlanDays > 0 && avgPrizeNum >= 1 {
		for day := 0; day < prizePlanDays; day++ {
			dayPrizeNumMap[day] = avgPrizeNum
		}
	}
	// 剩下的奖品一个一个的随机分配到任意哪天
	prizeNum -= prizePlanDays * avgPrizeNum
	for prizeNum > 0 {
		prizeNum--
		day := utils.Random(prizePlanDays)
		dayPrizeNumMap[day] += 1
	}
	// 发奖map：map[int]map[int][60]int
	//map[天]map[小时][60]奖品数量：后一个map表示value是一个60大小的数组，表示一个小时中每分钟要发的奖品数量
	prizePlanMap := make(map[int]map[int][60]int)
	//log.Infof("prize_id = %d\ndayPrizeNumMap = %+v", prize.Id, dayPrizeNumMap)
	for day, num := range dayPrizeNumMap {
		//计算一天的发奖计划
		dayPrizePlan := a.prizePlanOneDay(num)
		prizePlanMap[day] = dayPrizePlan
	}
	//log.Infof("prize_id = %d\nprizePlanMap = %+v", prize.Id, prizePlanMap)
	// 格式化 dayPrizePlan数据，序列化成为一个[时间:数量]二元组的数组
	planList, err := a.formatPrizePlan(now, prizePlanDays, prizePlanMap)
	if err != nil {
		log.ErrorContextf(ctx, "limitCase|ResetPrizePlan|formatPrizePlan err:", err)
		return fmt.Errorf("limitCase|ResetGiftPrizePlan:%v", err)
	}
	bytes, err := json.Marshal(planList)
	if err != nil {
		log.ErrorContextf(ctx, "limitCase|ResetPrizePlan|planList json marshal error=", err)
		return fmt.Errorf("limitCase|ResetGiftPrizePlan:%v", err)
	}
	// 保存奖品的分布计划数据
	info := &Prize{
		Id:         prize.Id,
		LeftNum:    prize.PrizeNum,
		PrizePlan:  string(bytes),
		PrizeBegin: now,
		PrizeEnd:   now.Add(time.Second * time.Duration(86400*prizePlanDays)),
	}
	err = a.prizeRepo.UpdateWithCache(info, "prize_plan", "prize_begin", "prize_end")
	if err != nil {
		log.ErrorContextf(ctx, "limitCase|ResetPrizePlan|prizeRepo.Update err:", err)
		return fmt.Errorf("limitCase|ResetPrizePlan:%v", err)
	}
	return nil
}

// clearPrizeData 清空奖品的发放计划
func (a *AdminCase) clearPrizePlan(ctx context.Context, prize *Prize) error {
	info := &Prize{
		Id:        prize.Id,
		PrizePlan: "",
	}
	err := a.prizeRepo.UpdateWithCache(info, "prize_plan")
	if err != nil {
		log.ErrorContextf(ctx, "limitCase|clearPrizePlan|prizeRepo.Update err", err)
		return fmt.Errorf("limitCase|clearPrizePlan:%v", err)
	}
	//奖品池也设为0
	if err = a.setPrizePool(ctx, prize.Id, 0); err != nil {
		return fmt.Errorf("limitCase|clearPrizePlan:%v", err)
	}
	return nil
}

// setGiftPool 设置奖品池中某种奖品的数量
func (a *AdminCase) setPrizePool(ctx context.Context, id uint, num int) error {
	key := constant.PrizePoolCacheKey
	if err := a.prizeRepo.SetPrizePoolNum(key, id, num); err != nil {
		log.ErrorContextf(ctx, "AdminCase|setPrizePool|%v", err)
		return fmt.Errorf("AdminCase|setPrizePool|%v", err)
	}
	return nil
}

// prizePlanOneDay 计算一天内具体到每小时每分钟应该发出的奖品，map[int][60]int： map[hour][minute]num
func (a *AdminCase) prizePlanOneDay(num int) map[int][60]int {
	resultMap := make(map[int][60]int)
	hourPrizeNumList := [24]int{} // 长度为24的数组表示1天中每个小时对应的奖品数
	// 计算一天中的24个小时，每个小时应该发出的奖品数，为什么是100，100表示每一天的权重百分比
	if num > 100 {
		for _, h := range DayPrizeWeights {
			hourPrizeNumList[h]++
		}
		for h := 0; h < 24; h++ {
			d := hourPrizeNumList[h]
			n := num * d / 100 // d / 100 每个小时所占的奖品数量概率
			hourPrizeNumList[h] = n
			num -= n
		}
	}
	//log.Infof("num = %d", num)
	for num > 0 {
		num--
		// 随机将这个奖品分配到某一个小时上
		hourIndex := utils.Random(100)
		//log.Infof("hourIndex = %d", hourIndex)
		h := DayPrizeWeights[hourIndex]
		hourPrizeNumList[h]++
	}
	//log.Infof("hourPrizeNumList = %v", hourPrizeNumList)
	// 将每个小时内的奖品数量分配到60分钟
	for h, hourPrizenum := range hourPrizeNumList {
		if hourPrizenum <= 0 {
			continue
		}
		minutePrizeNumList := [60]int{}
		if hourPrizenum >= 60 {
			avgMinutePrizeNum := hourPrizenum / 60
			for i := 0; i < 60; i++ {
				minutePrizeNumList[i] = avgMinutePrizeNum
			}
			hourPrizenum -= avgMinutePrizeNum * 60
		}
		for hourPrizenum > 0 {
			hourPrizenum--
			// 随机将这个奖品分配到某一分钟上
			m := utils.Random(60)
			//log.Infof("minuteIndex = %d", m)
			minutePrizeNumList[m]++
		}
		//log.Infof("minutePrizeNumList = %v", minutePrizeNumList)
		resultMap[h] = minutePrizeNumList
	}
	//log.Infof("resultMap=%v", resultMap)
	//log.Infof("-----------------------------------------------------------")
	return resultMap
}

// 将prizeData格式化成具体到一个时间（分钟）的奖品数量
// 结构为： [day][hour][minute]num
// result: [][时间, 数量]
func (a *AdminCase) formatPrizePlan(now time.Time, prizePlanDays int, prizePlan map[int]map[int][60]int) ([]*TimePrizeInfo, error) {
	result := make([]*TimePrizeInfo, 0)
	nowHour := now.Hour()
	for i := 0; i < prizePlanDays; i++ {
		dayPrizePlanMap, ok := prizePlan[i]
		if !ok {
			continue
		}
		dayTimeStamp := int(now.Unix()) + i*86400 // dayTimeStamp 为发奖周期中的每一天对应当前时间的时刻
		for h := 0; h < 24; h++ {
			hourPrizePlanMap, ok := dayPrizePlanMap[(h+nowHour)%24]
			if !ok {
				continue
			}
			hourTimeStamp := dayTimeStamp + h*3600 // hourTimeStamp 为发奖周期中的每一天中每个小时对应的时刻
			for m := 0; m < 60; m++ {
				num := hourPrizePlanMap[m]
				if num <= 0 {
					continue
				}
				// 找到特定一个时间的计划数据
				minuteTimeStamp := hourTimeStamp + m*60 // minuteTimeStamp 为发奖周期中的每一分钟对应的时刻
				result = append(result, &TimePrizeInfo{
					Time: utils.FormatFromUnixTime(int64(minuteTimeStamp)),
					Num:  num,
				})
			}
		}
	}
	return result, nil
}

// ResetAllPrizePlan 重置所有奖品的发奖计划
func (a *AdminCase) ResetAllPrizePlan() {
	//log.Infof("Resetting all prizes!!!!!")
	prizeList, err := a.GetPrizeList(context.Background())
	if err != nil {
		log.Errorf("ResetAllPrizePlan err:%v", err)
	}
	now := time.Now()
	for _, prize := range prizeList {
		if prize.PrizeTime > 0 && (prize.PrizePlan == "" || prize.PrizeEnd.Before(now)) {
			// ResetPrizePlan只会更新db的数据
			if err = a.ResetPrizePlan(context.Background(), prize); err != nil {
				log.Errorf("ResetAllPrizePlan err:%v", err)
			}
			// 通过读取缓存将db的数据同步到缓存中
			_, err = a.GetPrizeListWithCache(context.Background())
			if err != nil {
				log.Errorf("ResetAllPrizePlan err:%v", err)
			}
		}
	}
}

func (a *AdminCase) FillAllPrizePool() {
	log.Infof("FillAllPrizePool!!!!")
	totalNum, err := a.fillPrizePool()
	if err != nil {
		log.Errorf("FillAllPrizePool err:%v", err)
	}
	log.Infof("FillAllPrizePool with num:%d", totalNum)
}

func (a *AdminCase) fillPrizePool() (int, error) {
	totalNum := 0
	prizeList, err := a.GetPrizeList(context.Background())
	now := time.Now()
	if err != nil {
		log.Errorf("FillPrizePool err:%v", err)
		return 0, fmt.Errorf("FillPrizePool|GetPrizeList:%v", err)
	}
	if prizeList == nil || len(prizeList) == 0 {
		return 0, nil
	}
	for _, prize := range prizeList {
		if prize.SysStatus != 1 {
			continue
		}
		if prize.PrizeNum <= 0 {
			continue
		}
		if prize.BeginTime.After(now) || prize.EndTime.Before(now) {
			continue
		}
		// 发奖计划数据不正确
		if len(prize.PrizePlan) <= 7 {
			continue
		}
		prizePlanList := []*TimePrizeInfo{}
		if err = json.Unmarshal([]byte(prize.PrizePlan), &prizePlanList); err != nil {
			log.Errorf("FillPrizePool|Unmarshal TimePrizeInfo err:%v", err)
			return 0, fmt.Errorf("FillPrizePool|Unmarshal TimePrizeInfo:%v", err)
		}
		index := 0
		prizeNum := 0
		for i, prizePlanInfo := range prizePlanList {
			t, err := utils.ParseTime(prizePlanInfo.Time)
			if err != nil {
				log.Errorf("FillPrizePool|ParseTime err:%v", err)
				return 0, fmt.Errorf("FillPrizePool|ParseTime:%v", err)
			}
			if t.After(now) {
				break
			}
			//log.Infof("fillPrizePool|prize_id=%d\n", prize.Id, prizePlanInfo)
			// 该类奖品中，之前没有发放的奖品数量都要放入奖品池
			prizeNum += prizePlanInfo.Num
			index = i + 1
		}
		if prizeNum > 0 {
			a.incrPrizePool(prize.Id, prizeNum)
			totalNum += prizeNum
		}
		// 更新发奖计划
		if index > 0 {
			if index < len(prizePlanList) {
				prizePlanList = prizePlanList[index:]
			} else {
				prizePlanList = make([]*TimePrizeInfo, 0)
			}
			// 将新的发奖计划更新到数据库
			bytes, err := json.Marshal(&prizePlanList)
			if err != nil {
				log.Errorf("FillPrizePool|Marshal err:%v", err)
				return 0, fmt.Errorf("FillPrizePool|Marshal:%v", err)
			}
			updatePrize := &Prize{
				Id:        prize.Id,
				PrizePlan: string(bytes),
			}
			if err = a.UpdateDbPrizeWithCache(context.Background(), updatePrize, "prize_plan"); err != nil {
				log.Errorf("FillPrizePool|UpdateDbPrizeWithCache err:%v", err)
				return 0, fmt.Errorf("FillPrizePool|UpdateDbPrizeWithCache:%v", err)
			}
		}
		if totalNum > 0 {
			// totalNum>0,说明有奖品被填充到奖品池中，有奖品的发奖计划发生了改变，需要将更新后的数据加载到缓存中
			_, err = a.GetPrizeListWithCache(context.Background())
			if err != nil {
				log.Errorf("FillPrizePool|GetPrizeListWithCache err:%v", err)
				return 0, fmt.Errorf("FillPrizePool|GetPrizeListWithCache:%v", err)
			}
		}
	}
	return totalNum, nil
}

// incrPrizePool 根据计划数据，往奖品池增加奖品数量
func (a *AdminCase) incrPrizePool(prizeID uint, num int) (int, error) {
	key := constant.PrizePoolCacheKey
	cnt, err := a.prizeRepo.IncrPrizePoolNum(key, prizeID, num)
	if err != nil {
		log.Errorf("AdminCase|incrPrizePool|%v", err)
		return 0, fmt.Errorf("AdminCase|incrPrizePool|%v", err)
	}
	log.Infof("AdminCase|incrPrizePool|prize_id=%d, num=%d", prizeID, num)
	return cnt, nil
}

func (a *AdminCase) ClearLotteryTimes(ctx context.Context) error {
	if err := a.lotteryTimesRepo.DeleteAll(); err != nil {
		log.ErrorContextf(ctx, "adminCase|ClearCoupon err:%v", err)
		return fmt.Errorf("adminCase|ClearCoupon:%v", err)
	}
	return nil
}

func (a *AdminCase) ClearResult(ctx context.Context) error {
	if err := a.resultRepo.DeleteAll(); err != nil {
		log.ErrorContextf(ctx, "adminCase|ClearCoupon err:%v", err)
		return fmt.Errorf("adminCase|ClearCoupon:%v", err)
	}
	return nil
}
