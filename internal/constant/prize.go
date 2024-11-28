package constant

const (
	UserPrizeMax = 20    // 用户每天最多抽奖次数
	IpPrizeMax   = 30000 // 同一个IP每天最多抽奖次数
	IpLimitMax   = 3000  // 同一个IP每天最多抽奖次数
)

const (
	IpFrameSize   = 2
	UserFrameSize = 2
)

const (
	PrizeCodeMax = 10000
)

const (
	PrizeTypeVirtualCoin  = 0 // 虚拟币
	PrizeTypeCouponSame   = 1 // 虚拟券，相同的码
	PrizeTypeCouponDiff   = 2 // 虚拟券，不同的码
	PrizeTypeEntitySmall  = 3 // 实物小奖
	PrizeTypeEntityMiddle = 4 // 实物中等将
	PrizeTypeEntityLarge  = 5 // 实物大奖
)

const (
	DefaultBlackTime    = 7 * 86400  // 默认1周
	AllPrizeCacheTime   = 30 * 86400 // 默认1周
	CouponDiffLockLimit = 10000000
)

const (
	AllPrizeCacheKey        = "all_prize"
	UserCacheKeyPrefix      = "black_user_info_"
	IpCacheKeyPrefix        = "black_ip_info_"
	UserLotteryDayNumPrefix = "user_lottery_day_num_"
	PrizePoolCacheKey       = "prize_pool"
	PrizeCouponCacheKey     = "prize_coupon_"
)
