package biz

import "time"

type Coupon struct {
	Id         uint       `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	PrizeId    uint       `gorm:"column:prize_id;type:int(10) unsigned;default:0;comment:奖品ID，关联lt_prize表;NOT NULL" json:"prize_id"`
	Code       string     `gorm:"column:code;type:varchar(255);comment:虚拟券编码;NOT NULL" json:"code"`
	SysCreated *time.Time `gorm:"autoCreateTime;column:sys_created;type:datetime;default null;comment:创建时间;NOT NULL" json:"sys_created"`
	SysUpdated *time.Time `gorm:"autoUpdateTime;column:sys_updated;type:datetime;default null;comment:更新时间;NOT NULL" json:"sys_updated"`
	SysStatus  uint       `gorm:"column:sys_status;type:smallint(5) unsigned;default:0;comment:状态，1正常，2作废，2已发放;NOT NULL" json:"sys_status"`
}

func (c *Coupon) TableName() string {
	return "t_coupon"
}

type CouponRepo interface {
	Get(id uint) (*Coupon, error)
	GetAll() ([]*Coupon, error)
	GetCouponListByPrizeID(prizeID uint) ([]*Coupon, error)
	CountAll() (int64, error)
	Create(coupon *Coupon) error
	Delete(id uint) error
	DeleteAllWithCache() error
	Update(coupon *Coupon, cols ...string) error
	UpdateByCode(code string, coupon *Coupon, cols ...string) error
	GetFromCache(id uint) (*Coupon, error)
	GetGetNextUsefulCoupon(prizeID, couponID int) (*Coupon, error)
	ImportCacheCoupon(prizeID uint, code string) (bool, error)
	ReSetCacheCoupon(prizeID uint) (int64, int64, error)
	GetCacheCouponNum(prizeID uint) (int64, int64, error)
	GetNextUsefulCouponFromCache(prizeID int) (string, error)
}
