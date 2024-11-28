package constant

import "time"

const ReqID = "req_id"

// 时间标准化
const (
	SysTimeFormat      = "2006-01-02 15:04:05"
	SysTimeFormatShort = "2006-01-02"
)

const (
	TraceID = "Trace-ID"
	UserID  = "User-ID"
)

// 奖品状态
const (
	PrizeStatusNormal = 1 // 正常
	PrizeStatusDelete = 2 // 删除
)

const (
	Issuer              = "lottery"
	Expires             = 3600
	SecretKey           = "lottery-single"
	TokenExpireDuration = time.Hour * 2
)

const (
	LotteryLockKeyPrefix = "lucky_lock_"
)
