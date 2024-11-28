package biz

import "time"

// BlackUser 用户黑明单表
type BlackUser struct {
	Id         uint       `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	UserId     uint       `gorm:"column:user_id;type:int(10) unsigned;default:0;comment:用户ID;NOT NULL" json:"user_id"`
	UserName   string     `gorm:"column:user_name;type:varchar(50);comment:用户名;NOT NULL" json:"user_name"`
	BlackTime  time.Time  `gorm:"column:black_time;type:datetime;default:1000-01-01 00:00:00;comment:黑名单限制到期时间;NOT NULL" json:"black_time"`
	RealName   string     `gorm:"column:real_name;type:varchar(50);comment:真是姓名;NOT NULL" json:"real_name"`
	Mobile     string     `gorm:"column:mobile;type:varchar(50);comment:手机号;NOT NULL" json:"mobile"`
	Address    string     `gorm:"column:address;type:varchar(255);comment:联系地址;NOT NULL" json:"address"`
	SysCreated *time.Time `gorm:"autoCreateTime;column:sys_created;type:datetime;default null;comment:创建时间;NOT NULL" json:"sys_created"`
	SysUpdated *time.Time `gorm:"autoUpdateTime;column:sys_updated;type:datetime;default null;comment:修改时间;NOT NULL" json:"sys_updated"`
	SysIp      string     `gorm:"column:sys_ip;type:varchar(50);comment:IP地址;NOT NULL" json:"sys_ip"`
}

func (m *BlackUser) TableName() string {
	return "t_black_user"
}

type BlackUserRepo interface {
	GetByUserID(uid uint) (*BlackUser, error)
	GetByUserIDWithCache(uid uint) (*BlackUser, error)
	GetAll() ([]*BlackUser, error)
	CountAll() (int64, error)
	Create(blackUser *BlackUser) error
	Delete(id uint) error
	DeleteWithCache(uid uint) error
	Update(userID uint, blackUser *BlackUser, cols ...string) error
	UpdateWithCache(userID uint, blackUser *BlackUser, cols ...string) error
	GetFromCache(id uint) (*BlackUser, error)
	GetByCache(uid uint) (*BlackUser, error)
	SetByCache(blackUser *BlackUser) error
	UpdateByCache(blackUser *BlackUser) error
}
