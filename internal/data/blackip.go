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
)

type blackIpRepo struct {
	data *Data
}

func NewBlackIpRepo(data *Data) biz.BlackIpRepo {
	return &blackIpRepo{
		data: data,
	}
}

func (r *blackIpRepo) Get(id uint) (*biz.BlackIp, error) {
	// 优先从缓存获取
	db := r.data.db
	blackIp, err := r.GetFromCache(id)
	if err == nil && blackIp != nil {
		return blackIp, nil
	}
	blackIp = &biz.BlackIp{
		Id: id,
	}
	err = db.Model(&biz.BlackIp{}).First(blackIp).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("blackIpRepo|Get:%v", err)
	}
	return blackIp, nil
}

func (r *blackIpRepo) GetByIP(ip string) (*biz.BlackIp, error) {
	db := r.data.db
	blackIP := &biz.BlackIp{
		Ip: ip,
	}
	err := db.Model(&biz.BlackIp{}).Where("ip = ?", ip).First(blackIP).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("blackIpRepo|GetByIp:%v", err)
	}
	return blackIP, nil
}

func (r *blackIpRepo) GetByIPWithCache(ip string) (*biz.BlackIp, error) {
	db := r.data.db
	// 优先从缓存获取
	blackIp, err := r.GetByCache(ip)
	// 从缓存获取到IP
	if err == nil && blackIp != nil {
		return blackIp, nil
	}
	// 缓存中没有获取到ip
	blackIP := &biz.BlackIp{
		Ip: ip,
	}
	err = db.Model(&biz.BlackIp{}).Where("ip = ?", ip).First(blackIP).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		return nil, fmt.Errorf("blackIpRepo|GetByIp:%v", err)
	}
	// 数据库中正确读到数据，设置到缓存中
	if err = r.SetByCache(blackIP); err != nil {
		return nil, fmt.Errorf("blackIpRepo|SetByCache:%v", err)
	}
	return blackIP, nil
}

func (r *blackIpRepo) GetAll() ([]*biz.BlackIp, error) {
	db := r.data.db
	var BlackIps []*biz.BlackIp
	err := db.Model(&biz.BlackIp{}).Where("").Order("sys_updated desc").Find(&BlackIps).Error
	if err != nil {
		return nil, fmt.Errorf("blackIpRepo|GetAll:%v", err)
	}
	return BlackIps, nil
}

func (r *blackIpRepo) CountAll() (int64, error) {
	db := r.data.db
	var num int64
	err := db.Model(&biz.BlackIp{}).Count(&num).Error
	if err != nil {
		return 0, fmt.Errorf("blackIpRepo|CountAll:%v", err)
	}
	return num, nil
}

func (r *blackIpRepo) Create(blackIp *biz.BlackIp) error {
	db := r.data.db
	err := db.Model(blackIp).Create(blackIp).Error
	if err != nil {
		return fmt.Errorf("blackIpRepo|Create:%v", err)
	}
	return nil
}

func (r *blackIpRepo) Delete(id uint) error {
	db := r.data.db
	blackIp := &biz.BlackIp{Id: id}
	if err := db.Model(blackIp).Delete(blackIp).Error; err != nil {
		return fmt.Errorf("blackIpRepo|Delete:%v")
	}
	return nil
}

func (r *blackIpRepo) Update(ip string, blackIp *biz.BlackIp, cols ...string) error {
	db := r.data.db
	if err := r.UpdateByCache(&biz.BlackIp{Ip: ip}); err != nil {
		return fmt.Errorf("blackIpRepo|UpdateWithCache:%v", err)
	}
	var err error
	if len(cols) == 0 {
		err = db.Model(blackIp).Where("ip=?", ip).Updates(blackIp).Error
	} else {
		err = db.Model(blackIp).Where("ip=?", ip).Select(cols).Updates(blackIp).Error
	}
	if err != nil {
		return fmt.Errorf("blackIpRepo|Update:%v", err)
	}
	return nil
}

func (r *blackIpRepo) UpdateWithCache(ip string, blackIp *biz.BlackIp, cols ...string) error {
	db := r.data.db
	if err := r.UpdateByCache(&biz.BlackIp{Ip: ip}); err != nil {
		return fmt.Errorf("blackIpRepo|UpdateWithCache:%v", err)
	}
	var err error
	if len(cols) == 0 {
		err = db.Model(blackIp).Where("ip=?", ip).Updates(blackIp).Error
	} else {
		err = db.Model(blackIp).Where("ip=?", ip).Select(cols).Updates(blackIp).Error
	}
	if err != nil {
		return fmt.Errorf("blackIpRepo|Update:%v", err)
	}
	return nil
}

// GetFromCache 根据id从缓存获取奖品
func (r *blackIpRepo) GetFromCache(id uint) (*biz.BlackIp, error) {
	redisCli := r.data.cache
	idStr := strconv.FormatUint(uint64(id), 10)
	ret, exist, err := redisCli.Get(context.Background(), idStr)
	if err != nil {
		log.Errorf("blackIpRepo|GetFromCache:" + err.Error())
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	blackIp := biz.BlackIp{}
	if err = json.Unmarshal([]byte(ret), &blackIp); err != nil {
		return nil, fmt.Errorf("blackIpRepo|GetFromCache|json.Unmarshal:%v", err)
	}

	return &blackIp, nil
}

func (s *blackIpRepo) SetByCache(blackIp *biz.BlackIp) error {
	if blackIp == nil || blackIp.Ip == "" {
		return fmt.Errorf("blackIpRepo|SetByCache invalid user")
	}
	redisCli := s.data.cache
	key := fmt.Sprintf(constant.IpCacheKeyPrefix+"%s", blackIp.Ip)
	valueMap := make(map[string]string)
	valueMap["Id"] = strconv.Itoa(int(blackIp.Id))
	valueMap["BlackTime"] = utils.FormatFromUnixTime(blackIp.BlackTime.Unix())
	valueMap["SysCreated"] = utils.FormatFromUnixTime(blackIp.SysCreated.Unix())
	valueMap["SysUpdated"] = utils.FormatFromUnixTime(blackIp.SysUpdated.Unix())
	valueMap["Ip"] = blackIp.Ip
	_, err := redisCli.HMSet(context.Background(), key, valueMap)
	if err != nil {
		log.Errorf("blackUserRepo|SetByCache invalid user")
	}
	return nil
}

func (s *blackIpRepo) GetByCache(ip string) (*biz.BlackIp, error) {
	redisCli := s.data.cache
	key := fmt.Sprintf(constant.IpCacheKeyPrefix+"%s", ip)
	valueMap, err := redisCli.HGetAll(context.Background(), key)
	if err != nil {
		return nil, fmt.Errorf("blackIpRepo|GetByCache:%v", err)
	}
	idStr := valueMap["Id"]
	id, _ := strconv.Atoi(idStr)
	blackIp := &biz.BlackIp{
		Id: uint(id),
		Ip: ip,
	}
	blackTime, err := utils.ParseTime(valueMap["BlackTime"])
	if err != nil {
		return nil, fmt.Errorf("blackIpRepo|GetByCache:%v", err)
	}
	blackIp.BlackTime = blackTime
	sysCreated, err := utils.ParseTime(valueMap["SysCreated"])
	if err != nil {
		return nil, fmt.Errorf("blackIpRepo|GetByCache:%v", err)
	}
	blackIp.SysCreated = &sysCreated
	sysUpdated, err := utils.ParseTime(valueMap["SysUpdated"])
	if err != nil {
		return nil, fmt.Errorf("blackIpRepo|GetByCache:%v", err)
	}
	blackIp.SysUpdated = &sysUpdated
	return blackIp, nil
}

func (r *blackIpRepo) UpdateByCache(blackIp *biz.BlackIp) error {
	redisCli := r.data.cache
	if blackIp == nil || blackIp.Ip == "" {
		return fmt.Errorf("blackIpRepo|UpdateByCache invalid blackUser")
	}
	key := fmt.Sprintf(constant.UserCacheKeyPrefix+"%s", blackIp.Ip)
	if err := redisCli.Delete(context.Background(), key); err != nil {
		return fmt.Errorf("blackIpRepo|UpdateByCache:%v", err)
	}
	return nil
}
