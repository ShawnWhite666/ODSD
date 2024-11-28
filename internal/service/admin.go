package service

import (
	"context"
	"fmt"
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/pkg/middlewares/log"
)

// AddPrize 添加奖品
func (a *AdminService) AddPrize(ctx context.Context, viewPrize *biz.ViewPrize) error {
	if err := a.adminCase.AddPrize(ctx, viewPrize); err != nil {
		log.ErrorContextf(ctx, "adminService|AddPrize err:%v", err)
		return fmt.Errorf("adminService|AddPrize:%v", err)
	}
	return nil
}

// AddPrizeList 添加奖品列表
func (a *AdminService) AddPrizeList(ctx context.Context, viewPrizeList []*biz.ViewPrize) error {
	if err := a.adminCase.AddPrizeList(ctx, viewPrizeList); err != nil {
		log.ErrorContextf(ctx, "adminService|AddPrizeList err:%v", err)
		return fmt.Errorf("adminService|AddPrizeList:%v", err)
	}
	return nil
}

// ClearPrize 清空奖品
func (a *AdminService) ClearPrize(ctx context.Context) error {
	if err := a.adminCase.ClearPrize(ctx); err != nil {
		log.ErrorContextf(ctx, "adminService|ClearPrize err:%v", err)
		return fmt.Errorf("adminService|ClearPrize:%v", err)
	}
	return nil
}

func (a *AdminService) ImportCoupon(ctx context.Context, prizeID uint, codes string) error {
	successNum, failNum, err := a.adminCase.ImportCoupon(ctx, prizeID, codes)
	if err != nil {
		return fmt.Errorf("AdminService|ImportCoupon|%v", err)
	}
	log.Infof("ImportCoupon|successNum=%d|failNum=%d\n", successNum, failNum)
	return nil
}

func (a *AdminService) ClearCoupon(ctx context.Context) error {
	if err := a.adminCase.ClearCoupon(ctx); err != nil {
		log.ErrorContextf(ctx, "adminService|ClearCoupon err:%v", err)
		return fmt.Errorf("adminService|ClearCoupon:%v", err)
	}
	return nil
}

func (a *AdminService) ImportCouponWithCache(ctx context.Context, prizeID uint, codes string) error {
	successNum, failNum, err := a.adminCase.ImportCouponWithCache(ctx, prizeID, codes)
	if err != nil {
		return fmt.Errorf("AdminService|ImportCouponWithCache|%v", err)
	}
	log.Infof("ImportCouponWithCache|successNum=%d|failNum=%d\n", successNum, failNum)
	return nil
}

func (a *AdminService) ClearLotteryTimes(ctx context.Context) error {
	if err := a.adminCase.ClearLotteryTimes(ctx); err != nil {
		log.ErrorContextf(ctx, "adminService|ClearCoupon err:%v", err)
		return fmt.Errorf("adminService|ClearCoupon:%v", err)
	}
	return nil
}

func (a *AdminService) ClearResult(ctx context.Context) error {
	if err := a.adminCase.ClearResult(ctx); err != nil {
		log.ErrorContextf(ctx, "adminService|ClearCoupon err:%v", err)
		return fmt.Errorf("adminService|ClearCoupon:%v", err)
	}
	return nil
}
