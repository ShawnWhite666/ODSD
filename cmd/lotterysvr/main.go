package main

import (
	"flag"
	"github.com/BitofferHub/lotterysvr/internal/conf"
	"github.com/BitofferHub/lotterysvr/internal/task"
	"github.com/BitofferHub/pkg/middlewares/discovery"
	"github.com/BitofferHub/pkg/middlewares/log"
	"github.com/go-kratos/kratos/v2"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "lottery-svr"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")

}

func newApp(gs *grpc.Server, hs *http.Server, ts *task.TaskServer) *kratos.App {
	// new reg with etcd client
	reg := discovery.GetRegistrar()
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Server(
			gs,
			hs,
			ts,
		),
		kratos.Registrar(reg.Reg),
	)
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	InitSource(&bc)
	app, cleanup, err := wireApp(bc.GetServer(), bc.GetData())
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func InitSource(c *conf.Bootstrap) {
	l := c.GetLog()
	// 初始化日志
	log.Init(log.WithLogLevel(l.GetLevel()),
		log.WithFileName(l.GetFilename()),
		log.WithMaxSize(l.GetMaxSize()),
		log.WithMaxBackups(l.GetMaxBackups()),
		log.WithLogPath(l.GetLogPath()),
		log.WithConsole(l.GetConsole()))
	// // 注册服务
	discovery.NewRegistrar(c.GetMicro().GetLb().GetAddr())
}
