package data

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/pkg/middlewares/log"
	"gorm.io/gorm"
	"strconv"
)

type resultRepo struct {
	data *Data
}

func NewResultRepo(data *Data) biz.ResultRepo {
	return &resultRepo{
		data: data,
	}
}

func (r *resultRepo) Get(id uint) (*biz.Result, error) {
	db := r.data.db
	// 优先从缓存获取
	result, err := r.GetFromCache(id)
	if err == nil && result != nil {
		return result, nil
	}
	result = &biz.Result{
		Id: id,
	}
	err = db.Model(&biz.Result{}).First(result).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("resultRepo|Get:%v", err)
	}
	return result, nil
}

func (r *resultRepo) GetAll() ([]*biz.Result, error) {
	db := r.data.db
	var results []*biz.Result
	err := db.Model(&biz.Result{}).Where("").Order("sys_updated desc").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("resultRepo|GetAll:%v", err)
	}
	return results, nil
}

func (r *resultRepo) CountAll() (int64, error) {
	db := r.data.db
	var num int64
	err := db.Model(&biz.Result{}).Count(&num).Error
	if err != nil {
		return 0, fmt.Errorf("resultRepo|CountAll:%v", err)
	}
	return num, nil
}

func (r *resultRepo) Create(result *biz.Result) error {
	db := r.data.db
	err := db.Model(&biz.Result{}).Create(result).Error
	if err != nil {
		return fmt.Errorf("resultRepo|Create:%v", err)
	}
	return nil
}

func (r *resultRepo) Delete(id uint) error {
	db := r.data.db
	result := &biz.Result{Id: id}
	if err := db.Model(&biz.Result{}).Delete(result).Error; err != nil {
		return fmt.Errorf("resultRepo|Delete:%v")
	}
	return nil
}

func (r *resultRepo) DeleteAll() error {
	db := r.data.db
	if err := db.Exec("DELETE FROM t_result").Error; err != nil {
		log.Errorf("resultRepo|DeleteAll:%v", err)
		return fmt.Errorf("resultRepo|DeleteAll:%v", err)
	}
	return nil
}

func (r *resultRepo) Update(result *biz.Result, cols ...string) error {
	db := r.data.db
	var err error
	if len(cols) == 0 {
		err = db.Model(result).Updates(result).Error
	} else {
		err = db.Model(result).Select(cols).Updates(result).Error
	}
	if err != nil {
		return fmt.Errorf("resultRepo|Update:%v", err)
	}
	return nil
}

// GetFromCache 根据id从缓存获取奖品
func (r *resultRepo) GetFromCache(id uint) (*biz.Result, error) {
	redisCli := r.data.cache
	idStr := strconv.FormatUint(uint64(id), 10)
	ret, exist, err := redisCli.Get(context.Background(), idStr)
	if err != nil {
		log.Errorf("resultRepo|GetFromCache:" + err.Error())
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	result := biz.Result{}
	json.Unmarshal([]byte(ret), &biz.Result{})

	return &result, nil
}
