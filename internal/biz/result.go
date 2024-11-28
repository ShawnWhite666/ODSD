package biz

import "time"

// Result 抽奖记录表
type Result struct {
	Id         uint       `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	PrizeId    uint       `gorm:"column:prize_id;type:int(10) unsigned;default:0;comment:奖品ID，关联lt_prize表;NOT NULL" json:"prize_id"`
	PrizeName  string     `gorm:"column:prize_name;type:varchar(255);comment:奖品名称;NOT NULL" json:"prize_name"`
	PrizeType  uint       `gorm:"column:prize_type;type:int(10) unsigned;default:0;comment:奖品类型，同lt_prize. gtype;NOT NULL" json:"prize_type"`
	UserId     uint       `gorm:"column:user_id;type:int(10) unsigned;default:0;comment:用户ID;NOT NULL" json:"user_id"`
	UserName   string     `gorm:"column:user_name;type:varchar(50);comment:用户名;NOT NULL" json:"user_name"`
	PrizeCode  uint       `gorm:"column:prize_code;type:int(10) unsigned;default:0;comment:抽奖编号（4位的随机数）;NOT NULL" json:"prize_code"`
	PrizeData  string     `gorm:"column:prize_data;type:varchar(255);comment:获奖信息;NOT NULL" json:"prize_data"`
	SysCreated *time.Time `gorm:"autoCreateTime;column:sys_created;type:datetime;default null;comment:创建时间;NOT NULL" json:"sys_created"`
	SysIp      string     `gorm:"column:sys_ip;type:varchar(50);comment:用户抽奖的IP;NOT NULL" json:"sys_ip"`
	SysStatus  uint       `gorm:"column:sys_status;type:smallint(5) unsigned;default:0;comment:状态，0 正常，1删除，2作弊;NOT NULL" json:"sys_status"`
}

func (r *Result) TableName() string {
	return "t_result"
}

type ResultRepo interface {
	Get(id uint) (*Result, error)
	GetAll() ([]*Result, error)
	CountAll() (int64, error)
	Create(result *Result) error
	Delete(id uint) error
	DeleteAll() error
	Update(result *Result, cols ...string) error
	GetFromCache(id uint) (*Result, error)
}
