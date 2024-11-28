package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/pkg/middlewares/log"
	"gorm.io/gorm"
	"strconv"
)

// biz.Coupon 优惠券表

type couponRepo struct {
	data *Data
}

func NewCouponRepo(data *Data) biz.CouponRepo {
	return &couponRepo{
		data: data,
	}
}

func (r *couponRepo) Get(id uint) (*biz.Coupon, error) {
	db := r.data.db
	// 优先从缓存获取
	coupon, err := r.GetFromCache(id)
	if err == nil && coupon != nil {
		return coupon, nil
	}
	coupon = &biz.Coupon{
		Id: id,
	}
	err = db.Model(&biz.Coupon{}).First(coupon).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("couponRepo|Get:%v", err)
	}
	return coupon, nil
}

func (r *couponRepo) GetAll() ([]*biz.Coupon, error) {
	db := r.data.db
	var coupons []*biz.Coupon
	err := db.Model(&biz.Coupon{}).Order("sys_updated desc").Find(&coupons).Error
	if err != nil {
		return nil, fmt.Errorf("couponRepo|GetAll:%v", err)
	}
	return coupons, nil
}

func (r *couponRepo) GetCouponListByPrizeID(prizeID uint) ([]*biz.Coupon, error) {
	db := r.data.db
	var coupons []*biz.Coupon
	err := db.Model(&biz.Coupon{}).Where("prize_id=?", prizeID).Order("id desc").Find(&coupons).Error
	if err != nil {
		return nil, fmt.Errorf("couponRepo|GetAll:%v", err)
	}
	return coupons, nil
}

func (r *couponRepo) CountAll() (int64, error) {
	db := r.data.db
	var num int64
	err := db.Model(&biz.Coupon{}).Count(&num).Error
	if err != nil {
		return 0, fmt.Errorf("couponRepo|CountAll:%v", err)
	}
	return num, nil
}

func (r *couponRepo) Create(coupon *biz.Coupon) error {
	db := r.data.db
	err := db.Model(&biz.Coupon{}).Create(coupon).Error
	if err != nil {
		return fmt.Errorf("couponRepo|Create:%v", err)
	}
	return nil
}

func (r *couponRepo) Delete(id uint) error {
	db := r.data.db
	coupon := &biz.Coupon{Id: id}
	if err := db.Model(&biz.Coupon{}).Delete(coupon).Error; err != nil {
		return fmt.Errorf("couponRepo|Delete:%v", err)
	}
	return nil
}

func (r *couponRepo) DeleteAllWithCache() error {
	db := r.data.db
	couponList := make([]biz.Coupon, 0)
	if err := db.Model(&biz.Coupon{}).Select("prize_id").Distinct().Find(&couponList).Error; err != nil {
		log.Errorf("couponRepo|DeleteAllWithCache:%v", err)
		return fmt.Errorf("couponRepo|DeleteAllWithCache:%v", err)
	}
	if err := db.Exec("DELETE FROM t_coupon").Error; err != nil {
		log.Errorf("couponRepo|DeleteAllWithCache:%v", err)
		return fmt.Errorf("couponRepo|DeleteAllWithCache:%v", err)
	}
	log.Infof("couponRepo|DeleteAllWithCache|couponList=%v", couponList)
	for _, coupon := range couponList {
		key := fmt.Sprintf(constant.PrizeCouponCacheKey+"%d", coupon.PrizeId)
		if err := r.data.cache.Delete(context.Background(), key); err != nil {
			log.Errorf("couponRepo|DeleteAllWithCache|redis delete:%v", err)
			return fmt.Errorf("couponRepo|DeleteAllWithCache|redis delete:%v", err)
		}
	}
	return nil
}

func (r *couponRepo) Update(coupon *biz.Coupon, cols ...string) error {
	db := r.data.db
	var err error
	if len(cols) == 0 {
		err = db.Model(coupon).Updates(coupon).Error
	} else {
		err = db.Model(coupon).Select(cols).Updates(coupon).Error
	}
	if err != nil {
		return fmt.Errorf("couponRepo|Update:%v", err)
	}
	return nil
}

func (r *couponRepo) UpdateByCode(code string, coupon *biz.Coupon, cols ...string) error {
	db := r.data.db
	var err error
	if len(cols) == 0 {
		err = db.Model(coupon).Where("code = ?", code).Updates(coupon).Error
	} else {
		err = db.Model(coupon).Where("code = ?", code).Select(cols).Updates(coupon).Error
	}
	if err != nil {
		return fmt.Errorf("couponRepo|Update:%v", err)
	}
	return nil
}

// GetFromCache 根据id从缓存获取奖品
func (r *couponRepo) GetFromCache(id uint) (*biz.Coupon, error) {
	redisCli := r.data.cache
	idStr := strconv.FormatUint(uint64(id), 10)
	ret, exist, err := redisCli.Get(context.Background(), idStr)
	if err != nil {
		log.Errorf("couponRepo|GetFromCache:" + err.Error())
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	coupon := biz.Coupon{}
	json.Unmarshal([]byte(ret), &biz.Coupon{})

	return &coupon, nil
}

// GetGetNextUsefulCoupon 获取下一个可用编码的优惠券
func (r *couponRepo) GetGetNextUsefulCoupon(prizeID, couponID int) (*biz.Coupon, error) {
	db := r.data.db
	coupon := &biz.Coupon{}
	err := db.Model(coupon).Where("prize_id=?", prizeID).Where("id > ?", couponID).
		Where("sys_status = ?", 1).First(coupon).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("couponRepo|GetGetNextUsefulCoupon err:%v", err)
	}
	return coupon, nil
}

// ImportCacheCoupon 往缓存导入优惠券
func (r *couponRepo) ImportCacheCoupon(prizeID uint, code string) (bool, error) {
	redisCli := r.data.cache
	key := fmt.Sprintf(constant.PrizeCouponCacheKey+"%d", prizeID)
	cnt, err := redisCli.SAdd(context.Background(), key, code)
	if err != nil {
		return false, fmt.Errorf("couponRepo|ImportCacheCoupon:%v", err)
	}
	if cnt == 0 {
		return false, nil
	}
	return true, nil
}

// ReSetCacheCoupon 根据库存优惠券重置优惠券缓存
func (r *couponRepo) ReSetCacheCoupon(prizeID uint) (int64, int64, error) {
	redisCli := r.data.cache
	var successNum, failureNum int64 = 0, 0
	couponList, err := r.GetCouponListByPrizeID(prizeID)
	if err != nil {
		return 0, 0, fmt.Errorf("couponRepo")
	}
	if couponList == nil || len(couponList) == 0 {
		return 0, 0, nil
	}
	key := fmt.Sprintf(constant.PrizeCouponCacheKey+"%d", prizeID)
	// 这里先用临时keu统计，在原key上统计的话，因为db里的数量可能变化，没有同步到缓存中，比如db里面减少了10条数据，如果在原key上增加，那么缓存就会多处10条数据，所以根据db全部统计完了之后，在覆盖
	tmpKey := "tmp_" + key
	for _, coupon := range couponList {
		code := coupon.Code
		if coupon.SysStatus == 1 {
			cnt, err := redisCli.SAdd(context.Background(), tmpKey, code)
			if err != nil {
				return 0, 0, fmt.Errorf("couponRepo|ReSetCacheCoupon:%v", err)
			}
			if cnt <= 0 {
				failureNum++
			} else {
				successNum++
			}
		}
	}
	_, err = redisCli.Rename(context.Background(), tmpKey, key)
	if err != nil {
		return 0, 0, fmt.Errorf("couponRepo|ReSetCacheCoupon:%v", err)
	}
	return successNum, failureNum, nil
}

// GetCacheCouponNum 获取缓存中的剩余优惠券数量以及数据库中的剩余优惠券数量
func (r *couponRepo) GetCacheCouponNum(prizeID uint) (int64, int64, error) {
	redisCli := r.data.cache
	var dbNum, cacheNum int64 = 0, 0
	couponList, err := r.GetCouponListByPrizeID(prizeID)
	if err != nil {
		return 0, 0, fmt.Errorf("couponRepo|GetCacheCouponNum:%v", err)
	}
	if couponList == nil {
		return 0, 0, nil
	}
	for _, coupon := range couponList {
		if coupon.SysStatus == 1 {
			dbNum++
		}
	}
	key := fmt.Sprintf(constant.PrizeCouponCacheKey+"%d", prizeID)
	cacheNum, err = redisCli.SCard(context.Background(), key)
	if err != nil {
		return 0, 0, fmt.Errorf("couponRepo|GetCacheCouponNum:%v", err)
	}
	return dbNum, cacheNum, nil
}

// GetNextUsefulCouponFromCache 从缓存中拿出一个可用优惠券
func (r *couponRepo) GetNextUsefulCouponFromCache(prizeID int) (string, error) {
	redisCli := r.data.cache
	key := fmt.Sprintf(constant.PrizeCouponCacheKey+"%d", prizeID)
	code, err := redisCli.SPop(context.Background(), key)
	if err != nil {
		if err.Error() == "redis: nil" {
			//log.Infof("coupon not left")
			return "", nil
		}
		return "", fmt.Errorf("lotteryService|PrizeCouponDiffByCache:%v", err)
	}
	if code == "" {
		//log.Infof("lotteryService|PrizeCouponDiffByCache code is nil with prize_id=%d", prizeID)
		return "", nil
	}
	return code, nil
}
