package service

import "time"

// ViewPrize 对外返回的数据（区别于存储层的数据）
type ViewPrize struct {
	Id           uint      `json:"id"`
	Title        string    `json:"title"`
	Img          string    `json:"img"`
	PrizeNum     int       `json:"prize_num"`
	PrizeCode    string    `json:"prize_code"`
	PrizeTime    uint      `json:"prize_time"`
	LeftNum      int       `json:"left_num"`
	PrizeType    uint      `json:"prize_type"`
	PrizePlan    string    `json:"prize_plan"`
	BeginTime    time.Time `json:"begin_time"`
	EndTime      time.Time `json:"end_time"`
	DisplayOrder uint      `json:"display_order"`
	SysStatus    uint      `json:"sys_status"`
}
