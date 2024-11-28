//go:build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/conf"
	"github.com/BitofferHub/lotterysvr/internal/data"
	"github.com/BitofferHub/lotterysvr/internal/interfaces"
	"github.com/BitofferHub/lotterysvr/internal/server"
	"github.com/BitofferHub/lotterysvr/internal/service"
	"github.com/BitofferHub/lotterysvr/internal/task"
	"github.com/google/wire"
)

// wireApp
//
//	@Author <a href="https://bitoffer.cn">狂飙训练营</a>
//	@Description: wireApp init kratos application.
//	@param *conf.Server
//	@param *conf.Data
//	@return *kratos.App
//	@return func()
//	@return error
func wireApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		interfaces.ProviderSet,
		task.ProviderSet,
		newApp))
}
