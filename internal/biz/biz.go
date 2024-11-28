package biz

import (
	"context"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewLotteryCase, NewLimitCase, NewAdminCase)

// Transaction 解耦biz与data层，biz层只调用接口的方法
type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}
