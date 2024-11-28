package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/lotterysvr/internal/utils"
	"github.com/BitofferHub/pkg/middlewares/log"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type prizeRepo struct {
	data *Data
}

func NewPrizeRepo(data *Data) biz.PrizeRepo {
	return &prizeRepo{
		data: data,
	}
}

func (r *prizeRepo) Get(id uint) (*biz.Prize, error) {
	db := r.data.db
	prize := &biz.Prize{
		Id: id,
	}
	err := db.Model(&biz.Prize{}).First(prize).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("prizeRepo|Get:%v", err)
	}
	return prize, nil
}

func (r *prizeRepo) GetWithCache(id uint) (*biz.Prize, error) {
	prizeList, err := r.GetAllWithCache()
	if err != nil {
		return nil, fmt.Errorf("prizeRepo|GetWithCache:%v", err)
	}
	for _, prize := range prizeList {
		if prize.Id == id {
			return prize, nil
		}
	}
	return nil, nil
}

func (r *prizeRepo) GetAll() ([]*biz.Prize, error) {
	db := r.data.db
	var prizes []*biz.Prize
	err := db.Model(&biz.Prize{}).Find(&prizes).Error
	if err != nil {
		return nil, fmt.Errorf("prizeRepo|GetAll:%v", err)
	}
	return prizes, nil
}

func (r *prizeRepo) GetAllWithCache() ([]*biz.Prize, error) {
	prizeList, err := r.GetAllByCache()
	if err != nil {
		return nil, fmt.Errorf("prizeRepo|GetAllWithCache:%v", err)
	}
	if prizeList == nil {
		// 缓存没查到，从db获取
		prizeList, err = r.GetAll()
		if err != nil {
			return nil, fmt.Errorf("prizeRepo|GetAllWithCache:%v", err)
		}
		// 将数据更新到缓存中
		if err = r.SetAllByCache(prizeList); err != nil {
			return nil, fmt.Errorf("prizeRepo|GetAllWithCache:%v", err)
		}
	}
	return prizeList, nil
}

func (r *prizeRepo) CountAll() (int64, error) {
	db := r.data.db
	var num int64
	err := db.Model(&biz.Prize{}).Count(&num).Error
	if err != nil {
		return 0, fmt.Errorf("prizeRepo|CountAll:%v", err)
	}
	return num, nil
}

func (r *prizeRepo) CountAllWithCache() (int64, error) {
	prizeList, err := r.GetAllWithCache()
	if err != nil {
		return 0, fmt.Errorf("prizeRepo|CountAllWithCache:%v", err)
	}
	return int64(len(prizeList)), nil
}

func (r *prizeRepo) Create(prize *biz.Prize) error {
	db := r.data.db
	err := db.Model(&biz.Prize{}).Create(prize).Error
	if err != nil {
		return fmt.Errorf("prizeRepo|Create:%v", err)
	}
	return nil
}

func (r *prizeRepo) CreateInBatches(prizeList []biz.Prize) error {
	db := r.data.db
	err := db.Model(&biz.Prize{}).Create(&prizeList).Error
	if err != nil {
		return fmt.Errorf("prizeRepo|Create:%v", err)
	}
	return nil
}

func (r *prizeRepo) CreateWithCache(prize *biz.Prize) error {
	if err := r.UpdateByCache(prize); err != nil {
		return fmt.Errorf("prizeRepo|CreateWithCache:%v", err)
	}
	return r.Create(prize)
}

func (r *prizeRepo) Delete(id uint) error {
	db := r.data.db
	prize := &biz.Prize{Id: id}
	if err := db.Model(&biz.Prize{}).Delete(prize).Error; err != nil {
		return fmt.Errorf("prizeRepo|Delete:%v", err)
	}
	return nil
}

func (r *prizeRepo) DeleteAll() error {
	db := r.data.db
	if err := db.Exec("DELETE FROM t_prize").Error; err != nil {
		return fmt.Errorf("prizeRepo|DeleteAll:%v", err)
	}
	return nil
}

func (r *prizeRepo) DeleteWithCache(id uint) error {
	prize := &biz.Prize{
		Id: id,
	}
	if err := r.UpdateByCache(prize); err != nil {
		return fmt.Errorf("prizeRepo|DeleteWithCache:%v", err)
	}
	return r.Delete(id)
}

func (r *prizeRepo) Update(prize *biz.Prize, cols ...string) error {
	db := r.data.db
	var err error
	if len(cols) == 0 {
		err = db.Model(prize).Updates(prize).Error
	} else {
		err = db.Model(prize).Select(cols).Updates(prize).Error
	}
	if err != nil {
		return fmt.Errorf("prizeRepo|Update:%v", err)
	}
	return nil
}

func (r *prizeRepo) UpdateWithCache(prize *biz.Prize, cols ...string) error {
	if err := r.UpdateByCache(prize); err != nil {
		return fmt.Errorf("prizeRepo|UpdateWithCache:%v", err)
	}
	return r.Update(prize, cols...)
}

// GetFromCache 根据id从缓存获取奖品
func (r *prizeRepo) GetFromCache(id uint) (*biz.Prize, error) {
	redisCli := r.data.cache
	idStr := strconv.FormatUint(uint64(id), 10)
	ret, exist, err := redisCli.Get(context.Background(), idStr)
	if err != nil {
		log.Errorf("prizeRepo|GetFromCache:" + err.Error())
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	prize := biz.Prize{}
	json.Unmarshal([]byte(ret), &biz.Prize{})

	return &prize, nil
}

func (r *prizeRepo) GetAllUsefulPrizeList() ([]*biz.Prize, error) {
	db := r.data.db
	now := time.Now()
	list := make([]*biz.Prize, 0)
	err := db.Model(&biz.Prize{}).Where("begin_time<=?", now).Where("end_time >= ?", now).
		Where("prize_num>?", 0).Where("sys_status=?", 1).Order("sys_updated desc").
		Order("display_order asc").Find(&list).Error
	if err != nil {
		return nil, fmt.Errorf("prizeRepo|GetAllUsefulPrizeList:%v", err)
	}
	return list, nil
}

// GetAllUsefulPrizeListWithCache 筛选出符合条件的奖品列表
func (r *prizeRepo) GetAllUsefulPrizeListWithCache() ([]*biz.Prize, error) {
	// 优先从缓存取，缓存没取到，从db取
	prizeList, err := r.GetAllWithCache()
	if err != nil {
		return nil, fmt.Errorf("prizeRepo|GetAllUsefulPrizeListWithCache:%v", err)
	}
	now := time.Now()
	dataList := make([]*biz.Prize, 0)
	for _, prize := range prizeList {
		if prize.Id > 0 && prize.SysStatus == 1 && prize.PrizeNum > 0 &&
			prize.BeginTime.Before(now) && prize.EndTime.After(now) {
			dataList = append(dataList, prize)
		}
	}
	return dataList, nil
}

func (r *prizeRepo) DecrLeftNum(id int, num int) (bool, error) {
	db := r.data.db
	//log.Infof("id: %d, num: %d\n", id, num)
	res := db.Model(&biz.Prize{}).Where("id = ? and left_num >= ?", id, num).UpdateColumn("left_num", gorm.Expr("left_num - ?", num))
	if res.Error != nil {
		return false, fmt.Errorf("prizeRepo|DecrLeftNum:%v", res.Error)
	}
	if res.RowsAffected <= 0 {
		return false, nil
	}
	return true, nil
}

// DecrLeftNumByPool 奖品缓冲池 对应奖品数量递减
func (r *prizeRepo) DecrLeftNumByPool(prizeID int) (int64, error) {
	redisCli := r.data.cache
	key := constant.PrizePoolCacheKey
	field := strconv.Itoa(prizeID)
	cnt, err := redisCli.HIncrBy(context.Background(), key, field, -1)
	if err != nil {
		return -1, fmt.Errorf("prizeRepo|DecrLeftNumByPool:%v", err)
	}
	return cnt, nil
}

func (r *prizeRepo) IncrLeftNum(id int, column string, num int) error {
	db := r.data.db
	if err := db.Model(&biz.Prize{}).Where("id = ?", id).
		Update(column, gorm.Expr(column+" + ？", num)).Error; err != nil {
		return fmt.Errorf("prizeRepo|IncrLeftNum err: %v", err)
	}
	return nil
}

// SetAllByCache 全量数据保存到redis中
func (r *prizeRepo) SetAllByCache(prizeList []*biz.Prize) error {
	redisCli := r.data.cache
	value := ""
	if len(prizeList) > 0 {
		prizeMapList := make([]map[string]interface{}, len(prizeList))
		for i := 0; i < len(prizeList); i++ {
			prize := prizeList[i]
			prizeMap := make(map[string]interface{})
			prizeMap["Id"] = prize.Id
			prizeMap["Title"] = prize.Title
			prizeMap["PrizeNum"] = prize.PrizeNum
			prizeMap["LeftNum"] = prize.LeftNum
			prizeMap["PrizeCode"] = prize.PrizeCode
			prizeMap["PrizeTime"] = prize.PrizeTime
			prizeMap["Img"] = prize.Img
			prizeMap["DisplayOrder"] = prize.DisplayOrder
			prizeMap["PrizeType"] = prize.PrizeType
			prizeMap["PrizeProfile"] = prize.PrizeProfile
			prizeMap["BeginTime"] = utils.FormatFromUnixTime(prize.BeginTime.Unix())
			prizeMap["EndTime"] = utils.FormatFromUnixTime(prize.EndTime.Unix())
			prizeMap["PrizePlan"] = prize.PrizePlan
			prizeMap["PrizeBegin"] = utils.FormatFromUnixTime(prize.PrizeBegin.Unix())
			prizeMap["PrizeEnd"] = utils.FormatFromUnixTime(prize.PrizeEnd.Unix())
			prizeMap["SysStatus"] = prize.SysStatus
			prizeMap["SysCreated"] = utils.FormatFromUnixTime(prize.SysCreated.Unix())
			prizeMap["SysUpdated"] = utils.FormatFromUnixTime(prize.SysUpdated.Unix())
			prizeMap["SysIp"] = prize.SysIp
			prizeMapList[i] = prizeMap
		}
		bytes, err := json.Marshal(prizeMapList)
		if err != nil {
			log.Errorf("SetAllByCache|marshal err:%v", err)
			return fmt.Errorf("SetAllByCache|marshal err:%v", err)
		}
		value = string(bytes)
	}
	if err := redisCli.Set(context.Background(), constant.AllPrizeCacheKey, value, time.Second*time.Duration(constant.AllPrizeCacheTime)); err != nil {
		log.Errorf("SetAllByCache|set cache err:%v", err)
		return fmt.Errorf("SetAllByCache|set cache err:%v", err)
	}
	return nil
}

// GetAllByCache 从缓存中获取所有的奖品信息
func (r *prizeRepo) GetAllByCache() ([]*biz.Prize, error) {
	redisCli := r.data.cache
	valueStr, ok, err := redisCli.Get(context.Background(), constant.AllPrizeCacheKey)
	if err != nil {
		return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
	}
	// 缓存中没数据
	if !ok {
		return nil, nil
	}
	str := utils.GetString(valueStr, "")
	if str == "" {
		return nil, nil
	}
	// 将json数据反序列化
	prizeMapList := []map[string]interface{}{}
	err = json.Unmarshal([]byte(str), &prizeMapList)
	if err != nil {
		log.Errorf("prizeRepo|GetAllByCache:%v", err)
		return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
	}
	prizeList := make([]*biz.Prize, len(prizeMapList))
	for i := 0; i < len(prizeMapList); i++ {
		prizeMap := prizeMapList[i]
		id := utils.GetInt64FromMap(prizeMap, "Id", 0)
		if id <= 0 {
			prizeList[i] = &biz.Prize{}
			continue
		}
		prizeBegin, err := utils.ParseTime(utils.GetStringFromMap(prizeMap, "PrizeBegin", ""))
		if err != nil {
			log.Errorf("prizeRepo|GetAllByCache ParseTime PrizeBegin err:%v", err)
			return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
		}
		prizeEnd, err := utils.ParseTime(utils.GetStringFromMap(prizeMap, "PrizeEnd", ""))
		if err != nil {
			log.Errorf("prizeRepo|GetAllByCache ParseTime PrizeEnd err:%v", err)
			return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
		}
		beginTime, err := utils.ParseTime(utils.GetStringFromMap(prizeMap, "BeginTime", ""))
		if err != nil {
			log.Errorf("prizeRepo|GetAllByCache ParseTime BeginTime err:%v", err)
			return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
		}
		endTime, err := utils.ParseTime(utils.GetStringFromMap(prizeMap, "EndTime", ""))
		if err != nil {
			log.Errorf("prizeRepo|GetAllByCache ParseTime EndTime err:%v", err)
			return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
		}
		sysCreated, err := utils.ParseTime(utils.GetStringFromMap(prizeMap, "SysCreated", ""))
		if err != nil {
			log.Errorf("prizeRepo|GetAllByCache ParseTime SysCreated err:%v", err)
			return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
		}
		sysUpdated, err := utils.ParseTime(utils.GetStringFromMap(prizeMap, "SysUpdated", ""))
		if err != nil {
			log.Errorf("prizeRepo|GetAllByCache ParseTime SysUpdated err:%v", err)
			return nil, fmt.Errorf("prizeRepo|GetAllByCache:%v", err)
		}
		prize := &biz.Prize{
			Id:           uint(id),
			Title:        utils.GetStringFromMap(prizeMap, "Title", ""),
			PrizeNum:     int(utils.GetInt64FromMap(prizeMap, "PrizeNum", 0)),
			LeftNum:      int(utils.GetInt64FromMap(prizeMap, "LeftNum", 0)),
			PrizeCode:    utils.GetStringFromMap(prizeMap, "PrizeCode", ""),
			PrizeTime:    uint(utils.GetInt64FromMap(prizeMap, "PrizeTime", 0)),
			Img:          utils.GetStringFromMap(prizeMap, "Img", ""),
			DisplayOrder: uint(utils.GetInt64FromMap(prizeMap, "DisplayOrder", 0)),
			PrizeType:    uint(utils.GetInt64FromMap(prizeMap, "PrizeType", 0)),
			PrizeProfile: utils.GetStringFromMap(prizeMap, "PrizeProfile", ""),
			BeginTime:    beginTime,
			EndTime:      endTime,
			// PrizeData:    comm.GetStringFromMap(data, "PrizeData", ""),
			PrizeBegin: prizeBegin,
			PrizeEnd:   prizeEnd,
			SysStatus:  uint(utils.GetInt64FromMap(prizeMap, "SysStatus", 0)),
			SysCreated: &sysCreated,
			SysUpdated: &sysUpdated,
			SysIp:      utils.GetStringFromMap(prizeMap, "SysIp", ""),
		}
		prizeList[i] = prize
	}
	return prizeList, nil
}

// UpdateByCache 数据更新，需要更新缓存，直接清空缓存数据
func (r *prizeRepo) UpdateByCache(prize *biz.Prize) error {
	if prize == nil || prize.Id <= 0 {
		return nil
	}
	redisCli := r.data.cache
	if err := redisCli.Delete(context.Background(), constant.AllPrizeCacheKey); err != nil {
		return fmt.Errorf("prizeRepo|UpdateByCache err:%v", err)
	}
	return nil
}

// GetPrizePoolNum 获取奖品缓冲池中获取数据
func (r *prizeRepo) GetPrizePoolNum(prizeID uint) (int, error) {
	redisCli := r.data.cache
	key := constant.PrizePoolCacheKey
	field := strconv.Itoa(int(prizeID))
	res, err := redisCli.HGet(context.Background(), key, field)
	if err != nil {
		return 0, fmt.Errorf("prizeRepo|GetPrizePoolNum:%v", err)
	}
	num, err := strconv.Atoi(res)
	if err != nil {
		return 0, fmt.Errorf("prizeRepo|GetPrizePoolNum:%v", err)
	}
	return num, nil
}

func (r *prizeRepo) SetPrizePoolNum(key string, prizeID uint, num int) error {
	idStr := strconv.Itoa(int(prizeID))
	_, err := r.data.cache.HSet(context.Background(), key, idStr, strconv.Itoa(num))
	if err != nil {
		return fmt.Errorf("adminCase|setPrizePool:%v", err)
	}
	return nil
}

func (r *prizeRepo) IncrPrizePoolNum(key string, prizeID uint, num int) (int, error) {
	idStr := strconv.Itoa(int(prizeID))
	cnt, err := r.data.cache.HIncrBy(context.Background(), key, idStr, int64(num))
	if err != nil {
		log.Errorf("incrPrizePool err:%v", err)
		return 0, fmt.Errorf("incrPrizePool err:%v", err)
	}
	if int(cnt) < num {
		//log.Infof("incrPrizePool twice,num=%d,cnt=%d", num, int(cnt))
		left := num - int(cnt)
		// 数量不等，存在没有成功的情况，补偿一次
		cnt, err = r.data.cache.HIncrBy(context.Background(), key, idStr, int64(left))
		if err != nil {
			log.Errorf("incrPrizePool twice err:%v", err)
			return 0, fmt.Errorf("incrPrizePool err:%v", err)
		}
	}
	return int(cnt), nil
}
