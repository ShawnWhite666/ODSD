// routes.go

package interfaces

import (
	"github.com/BitofferHub/lotterysvr/internal/service"
	engine "github.com/BitofferHub/pkg/middlewares/gin"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	lotteryService *service.LotteryService
	adminService   *service.AdminService
}

func NewHandler(s *service.LotteryService, a *service.AdminService) *Handler {
	return &Handler{
		lotteryService: s,
		adminService:   a,
	}
}

func NewRouter(h *Handler) *gin.Engine {
	r := engine.NewEngine(engine.WithLogger(false))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	adminGroup := r.Group("admin")
	// 获取奖品列表
	//adminGroup.GET("/get_prize_list", handlers.GetPrizeList)
	// 添加奖品
	adminGroup.POST("/add_prize", h.AddPrize)
	// 添加奖品列表
	adminGroup.POST("/add_prize_list", h.AddPrizeList)
	// 清空奖品
	adminGroup.POST("/clear_prize", h.ClearPrize)
	// 导入优惠券
	adminGroup.POST("/import_coupon", h.ImportCoupon)
	// 导入优惠券，同时导入缓存
	adminGroup.POST("/import_coupon_cache", h.ImportCouponWithCache)
	// 清空优惠券
	adminGroup.POST("/clear_coupon", h.ClearCoupon)
	// 清空用户抽奖次数
	adminGroup.POST("/clear_lottery_times", h.ClearLotteryTimes)
	// 清空获奖结果
	adminGroup.POST("/clear_result", h.ClearResult)

	lotteryGroup := r.Group("lottery")
	// V1基础版获取中奖
	lotteryGroup.POST("/v1/get_lucky", h.LotteryV1)
	// 优化V2版中奖逻辑
	lotteryGroup.POST("/v2/get_lucky", h.LotteryV2)
	// 优化V3版中奖逻辑
	lotteryGroup.POST("/v3/get_lucky", h.LotteryV3)
	return r
}
