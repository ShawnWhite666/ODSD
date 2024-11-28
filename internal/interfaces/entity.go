package interfaces

import (
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/constant"
)

// HttpResponse http独立请求返回结构体,这个通用的，不需要修改
type HttpResponse struct {
	Code   constant.ErrCode `json:"code"`
	Msg    string           `json:"msg"`
	Data   interface{}      `json:"data"`
	UserID uint32           `json:"user_id"`
}

type LotteryReq struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	IP       string `json:"ip"`
}

type AddPrizeReq struct {
	UserID uint           `json:"user_id"`
	Prize  *biz.ViewPrize `json:"prize"`
}

type AddPrizeListReq struct {
	UserID    uint             `json:"user_id"`
	PrizeList []*biz.ViewPrize `json:"prize_list"`
}

type ClearPrizeReq struct {
	UserID uint `json:"user_id"`
}

type ImportCouponReq struct {
	UserID     uint                `json:"user_id"`
	CouponInfo *biz.ViewCouponInfo `json:"coupon"`
}

type ClearCouponReq struct {
	UserID uint `json:"user_id"`
}

type ClearLotteryTimesReq struct {
	UserID uint `json:"user_id"`
}

type ClearResultReq struct {
	UserID uint `json:"user_id"`
}
