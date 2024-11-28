package biz

import "time"

// Prize 奖品表
type Prize struct {
	Id           uint       `gorm:"column:id;type:int(10) unsigned;primary_key;AUTO_INCREMENT" json:"id"`
	Title        string     `gorm:"column:title;type:varchar(255);comment:奖品名称;NOT NULL" json:"title"`
	PrizeNum     int        `gorm:"column:prize_num;type:int(11);default:-1;comment:奖品数量，0 无限量，>0限量，<0无奖品;NOT NULL" json:"prize_num"`
	LeftNum      int        `gorm:"column:left_num;type:int(11);default:0;comment:剩余数量;NOT NULL" json:"left_num"`
	PrizeCode    string     `gorm:"column:prize_code;type:varchar(50);comment:0-9999表示100%，0-0表示万分之一的中奖概率;NOT NULL" json:"prize_code"`
	PrizeTime    uint       `gorm:"column:prize_time;type:int(10) unsigned;default:0;comment:发奖周期，多少天，以天为单位;NOT NULL" json:"prize_time"`
	Img          string     `gorm:"column:img;type:varchar(255);comment:奖品图片;NOT NULL" json:"img"`
	DisplayOrder uint       `gorm:"column:display_order;type:int(10) unsigned;default:0;comment:位置序号，小的排在前面;NOT NULL" json:"display_order"`
	PrizeType    uint       `gorm:"column:prize_type;type:int(10) unsigned;default:0;comment:奖品类型，0 虚拟币，1 虚拟券，2 实物-小奖，3 实物-大奖;NOT NULL" json:"prize_type"`
	PrizeProfile string     `gorm:"column:prize_profile;type:varchar(255);comment:奖品扩展数据，如：虚拟币数量;NOT NULL" json:"prize_profile"`
	BeginTime    time.Time  `gorm:"column:begin_time;type:datetime;default:1000-01-01 00:00:00;comment:奖品有效周期：开始时间;NOT NULL" json:"begin_time"`
	EndTime      time.Time  `gorm:"column:end_time;type:datetime;default:1000-01-01 00:00:00;comment:奖品有效周期：结束时间;NOT NULL" json:"end_time"`
	PrizePlan    string     `gorm:"column:prize_plan;type:mediumtext;comment:发奖计划，[[时间1,数量1],[时间2,数量2]]" json:"prize_plan"`
	PrizeBegin   time.Time  `gorm:"column:prize_begin;type:int(11);default:1000-01-01 00:00:00;comment:发奖计划周期的开始;NOT NULL" json:"prize_begin"`
	PrizeEnd     time.Time  `gorm:"column:prize_end;type:int(11);default:1000-01-01 00:00:00;comment:发奖计划周期的结束;NOT NULL" json:"prize_end"`
	SysStatus    uint       `gorm:"column:sys_status;type:smallint(5) unsigned;default:0;comment:状态，0 正常，1 删除;NOT NULL" json:"sys_status"`
	SysCreated   *time.Time `gorm:"autoCreateTime:datetime;column:sys_created;type:datetime;default null;comment:创建时间;NOT NULL" json:"sys_created"`
	SysUpdated   *time.Time `gorm:"autoUpdateTime:datetime;column:sys_updated;type:datetime;default null;comment:修改时间;NOT NULL" json:"sys_updated"`
	SysIp        string     `gorm:"column:sys_ip;type:varchar(50);comment:操作人IP;NOT NULL" json:"sys_ip"`
}

func (p *Prize) TableName() string {
	return "t_prize"
}

type PrizeRepo interface {
	Get(id uint) (*Prize, error)
	GetWithCache(id uint) (*Prize, error)
	GetAll() ([]*Prize, error)
	GetAllWithCache() ([]*Prize, error)
	CountAll() (int64, error)
	CountAllWithCache() (int64, error)
	Create(prize *Prize) error
	CreateInBatches(prizeList []Prize) error
	CreateWithCache(prize *Prize) error
	Delete(id uint) error
	DeleteAll() error
	DeleteWithCache(id uint) error
	Update(prize *Prize, cols ...string) error
	UpdateWithCache(prize *Prize, cols ...string) error
	GetFromCache(id uint) (*Prize, error)
	GetAllUsefulPrizeList() ([]*Prize, error)
	GetAllUsefulPrizeListWithCache() ([]*Prize, error)
	DecrLeftNum(id int, num int) (bool, error)
	DecrLeftNumByPool(prizeID int) (int64, error)
	IncrLeftNum(id int, column string, num int) error
	SetAllByCache(prizeList []*Prize) error
	GetAllByCache() ([]*Prize, error)
	UpdateByCache(prize *Prize) error
	GetPrizePoolNum(prizeID uint) (int, error)
	SetPrizePoolNum(key string, prizeID uint, num int) error
	IncrPrizePoolNum(key string, prizeID uint, num int) (int, error)
}
