package interfaces

import (
	"context"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"github.com/BitofferHub/lotterysvr/internal/utils"
	"github.com/BitofferHub/pkg/middlewares/log"
	pb "github.com/BitofferHub/proto_center/api/lotterysvr/v1"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) LotteryV1(c *gin.Context) {
	rsp := HttpResponse{}
	// HTTP响应
	req := LotteryReq{}
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("LotteryV1|Error binding:%v", err)
		rsp.Code = constant.ErrShouldBind
		rsp.Msg = constant.GetErrMsg(rsp.Code)
		c.JSON(http.StatusOK, rsp)
		return
	}
	userIDStr := c.Request.Header.Get(constant.UserID)
	userID, err := strconv.Atoi(userIDStr)
	if err == nil {
		req.UserID = uint(userID)
	}
	//log.Infof("LotteryV1|Handler|req=====%+v", req)
	h.lotteryV1(&req, &rsp)
	c.JSON(http.StatusOK, rsp)
}

func (h *Handler) lotteryV1(lotteryReq *LotteryReq, lotteryRsp *HttpResponse) {
	ctx := context.WithValue(context.Background(), constant.ReqID, utils.NewUuid())
	req := &pb.LotteryReq{
		UserId:   uint32(lotteryReq.UserID),
		UserName: lotteryReq.UserName,
		Ip:       lotteryReq.IP,
	}
	// 2. 验证用户今日抽奖次数
	rsp, err := h.lotteryService.LotteryV1(ctx, req)
	if err != nil {
		log.ErrorContextf(ctx, "http lotteryv1|err:%v", err)
		return
	}
	lotteryRsp.Code = constant.ErrCode(rsp.CommonRsp.Code)
	lotteryRsp.Msg = rsp.CommonRsp.Msg
	lotteryRsp.Data = rsp.PrizeInfo
	lotteryRsp.UserID = rsp.CommonRsp.UserId
}
