package task

import (
	"context"
	"github.com/BitofferHub/lotterysvr/internal/conf"
	"github.com/BitofferHub/lotterysvr/internal/service"
	"github.com/BitofferHub/lotterysvr/internal/utils"
	"github.com/google/wire"
	"time"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewTaskServer)

type TaskServer struct {
	// 需要什么service, 就修改成自己的service
	service   *service.LotteryService
	scheduler *TaskScheduler
}

func (t *TaskServer) Stop(ctx context.Context) error {
	t.scheduler.Stop()
	return nil
}

// NewJobs 添加Job方法
func (t *TaskServer) NewJobs() []Job {
	return []Job{t.job1, t.job2, t.job3, t.job4}
}

// NewTaskServer 注入对应service
func NewTaskServer(s *service.LotteryService, c *conf.Server) *TaskServer {
	t := &TaskServer{
		service: s,
	}
	conf := c.GetTask()
	t.scheduler = NewScheduler(conf.GetAddr(), NewTasks(conf, t.NewJobs()))

	return t
}
func NewTasks(c *conf.Server_TASK, jobs []Job) []*Task {
	var tasks []*Task
	for i, job := range jobs {
		tasks = append(tasks, &Task{
			Name:     c.Tasks[i].Name,
			Type:     c.Tasks[i].Type,
			Schedule: c.Tasks[i].Schedule,
			Handler:  job,
		})
	}

	return tasks
}

func (t *TaskServer) job1() {
	t.service.CronJobResetIPLotteryNumsTask()
	next := utils.NextDayTime()
	t.scheduler.AddTask(Task{
		Name:     "job1",
		Type:     "once",
		NextTime: next,
		Handler:  t.job1,
	})
}

func (t *TaskServer) job2() {
	t.service.CronJobResetUserLotteryNumsTask()
	next := utils.NextDayTime()
	t.scheduler.AddTask(Task{
		Name:     "job2",
		Type:     "once",
		NextTime: next,
		Handler:  t.job2,
	})
}

func (t *TaskServer) job3() {
	t.service.CronJobResetAllPrizePlanTask()
	next := time.Now().Add(5 * time.Minute)
	t.scheduler.AddTask(Task{
		Name:     "job3",
		Type:     "once",
		NextTime: next,
		Handler:  t.job3,
	})
}

func (t *TaskServer) job4() {
	t.service.CronJobFillAllPrizePoolTask()
	next := time.Now().Add(1 * time.Minute)
	t.scheduler.AddTask(Task{
		Name:     "job4",
		Type:     "once",
		NextTime: next,
		Handler:  t.job4,
	})
}
