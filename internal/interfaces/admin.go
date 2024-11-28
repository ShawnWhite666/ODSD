package interfaces

import (
	"context"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/lotterysvr/internal/utils"
	"github.com/BitofferHub/pkg/middlewares/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddPrize 添加奖品
func (h *Handler) AddPrize(c *gin.Context) {
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	// HTTP响应
	req := AddPrizeReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("AddPrize|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("AddPrize|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	err := h.adminService.AddPrize(ctx, req.Prize)
	if err != nil {
		log.Errorf("AddPrize|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

// AddPrizeList 添加奖品列表
func (h *Handler) AddPrizeList(c *gin.Context) {
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	// HTTP响应
	req := AddPrizeListReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("AddPrizeList|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("AddPrizeList|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	err := h.adminService.AddPrizeList(ctx, req.PrizeList)
	if err != nil {
		log.Errorf("AddPrizeList|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	c.JSON(http.StatusOK, rsp)
}

// ClearPrize 清空奖品数据
func (h *Handler) ClearPrize(c *gin.Context) {
	req := ClearPrizeReq{}
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("ClearPrize|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("ClearPrize|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	if err := h.adminService.ClearPrize(ctx); err != nil {
		log.Errorf("ClearPrize|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
	}
	c.JSON(http.StatusOK, rsp)
}

// ImportCoupon 导入优惠券
func (h *Handler) ImportCoupon(c *gin.Context) {
	req := ImportCouponReq{}
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("ImportCoupon|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("ImportCoupon|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	if err := h.adminService.ImportCoupon(ctx, req.CouponInfo.PrizeId, req.CouponInfo.Code); err != nil {
		log.Errorf("ImportCoupon|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
	}
	c.JSON(http.StatusOK, rsp)
}

// ImportCouponWithCache 导入优惠券,并导入缓存
func (h *Handler) ImportCouponWithCache(c *gin.Context) {
	req := ImportCouponReq{}
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("ImportCoupon|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("ImportCoupon|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	if err := h.adminService.ImportCouponWithCache(ctx, req.CouponInfo.PrizeId, req.CouponInfo.Code); err != nil {
		log.Errorf("ImportCoupon|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
	}
	c.JSON(http.StatusOK, rsp)
}

// ClearCoupon 清空优惠券
func (h *Handler) ClearCoupon(c *gin.Context) {
	req := ClearCouponReq{}
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("ClearCoupon|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("ClearCoupon|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	if err := h.adminService.ClearCoupon(ctx); err != nil {
		log.Errorf("ClearCoupon|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
	}
	c.JSON(http.StatusOK, rsp)
}

// ClearLotteryTimes 清空用户抽奖次数
func (h *Handler) ClearLotteryTimes(c *gin.Context) {
	req := ClearLotteryTimesReq{}
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("ClearLotteryTimes|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("ClearLotteryTimes|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	if err := h.adminService.ClearLotteryTimes(ctx); err != nil {
		log.Errorf("ClearLotteryTimes|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
	}
	c.JSON(http.StatusOK, rsp)
}

// ClearResult 清空用户抽奖次数
func (h *Handler) ClearResult(c *gin.Context) {
	req := ClearResultReq{}
	rsp := HttpResponse{
		Code: constant.Success,
		Msg:  constant.GetErrMsg(constant.Success),
	}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("ClearResult|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	if req.UserID <= 0 {
		log.Errorf("ClearResult|user_id invalid")
		rsp.Code = constant.ErrInputInvalid
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	if err := h.adminService.ClearResult(ctx); err != nil {
		log.Errorf("ClearResult|err:%v", err)
		rsp.Code = constant.ErrInternalServer
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
	}
	c.JSON(http.StatusOK, rsp)
}
