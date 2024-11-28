package biz

import "time"

// DayPrizeWeights 定义一天中24小时内，每个小时的发奖比例权重，100的数组，0-23出现的次数为权重大小
var DayPrizeWeights = [100]int{
	0, 0, 0,
	1, 1, 1,
	2, 2, 2,
	3, 3, 3,
	4, 4, 4,
	5, 5, 5,
	6, 6, 6,
	7, 7, 7,
	8, 8, 8, 8, 8, 8, 8,
	9, 9, 9,
	10, 10, 10,
	11, 11, 11,
	12, 12, 12,
	13, 13, 13,
	14, 14, 14, 14, 14, 14, 14,
	15, 15, 15, 15, 15, 15, 15,
	16, 16, 16, 16, 16, 16, 16,
	17, 17, 17, 17, 17, 17, 17,
	18, 18, 18,
	19, 19, 19,
	20, 20, 20, 20, 20, 20, 20,
	21, 21, 21, 21, 21, 21, 21,
	22, 22, 22,
	23, 23, 23,
}

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

type LoginRsp struct {
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
}

// LotteryPrize 中奖奖品信息
type LotteryPrize struct {
	Id            uint   `json:"id"`
	Title         string `json:"title"`
	PrizeNum      int    `json:"-"`
	LeftNum       int    `json:"-"`
	PrizeCodeLow  int    `json:"-"`
	PrizeCodeHigh int    `json:"-"`
	Img           string `json:"img"`
	DisplayOrder  uint   `json:"display_order"`
	PrizeType     uint   `json:"prize_type"`
	PrizeProfile  string `json:"prize_profile"`
	CouponCode    string `json:"coupon_code"` // 如果中奖奖品是优惠券，这个字段位优惠券编码，否则为空
}

type LotteryUserInfo struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	IP       string `json:"ip"`
}

type ViewCouponInfo struct {
	Id         uint      `json:"id"`
	PrizeId    uint      `json:"prize_id"`
	Code       string    `json:"code"`
	SysCreated time.Time `json:"sys_created"`
	SysUpdated time.Time `json:"sys_updated"`
	SysStatus  uint      `json:"sys_status"`
}

type TimePrizeInfo struct {
	Time string `json:"time"`
	Num  int    `json:"num"`
}
