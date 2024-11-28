// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/conf"
	"github.com/BitofferHub/lotterysvr/internal/data"
	"github.com/BitofferHub/lotterysvr/internal/interfaces"
	"github.com/BitofferHub/lotterysvr/internal/server"
	"github.com/BitofferHub/lotterysvr/internal/service"
	"github.com/BitofferHub/lotterysvr/internal/task"
	"github.com/go-kratos/kratos/v2"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description: wireApp init kratos application.
//	@param *conf.Server
//	@param *conf.Data
//	@return *kratos.App
//	@return func()
//	@return error
func wireApp(confServer *conf.Server, confData *conf.Data) (*kratos.App, func(), error) {
	db := data.NewDatabase(confData)
	client := data.NewCache(confData)
	dataData := data.NewData(db, client)
	prizeRepo := data.NewPrizeRepo(dataData)
	couponRepo := data.NewCouponRepo(dataData)
	blackUserRepo := data.NewBlackUserRepo(dataData)
	blackIpRepo := data.NewBlackIpRepo(dataData)
	resultRepo := data.NewResultRepo(dataData)
	transaction := data.NewTransaction(dataData)
	lotteryCase := biz.NewLotteryCase(prizeRepo, couponRepo, blackUserRepo, blackIpRepo, resultRepo, transaction)
	lotteryTimesRepo := data.NewLotteryTimesRepo(dataData)
	limitCase := biz.NewLimitCase(blackUserRepo, blackIpRepo, lotteryTimesRepo, transaction)
	adminCase := biz.NewAdminCase(prizeRepo, couponRepo, lotteryTimesRepo, resultRepo)
	lotteryService := service.NewLotteryService(lotteryCase, limitCase, adminCase)
	grpcServer := server.NewGRPCServer(confServer, lotteryService)
	adminService := service.NewAdminService(adminCase)
	handler := interfaces.NewHandler(lotteryService, adminService)
	httpServer := server.NewHTTPServer(confServer, handler)
	taskServer := task.NewTaskServer(lotteryService, confServer)
	app := newApp(grpcServer, httpServer, taskServer)
	return app, func() {
	}, nil
}