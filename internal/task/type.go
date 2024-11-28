package task

import (
	"context"
	"fmt"
	"github.com/BitofferHub/pkg/middlewares/lock"
	cronpkg "github.com/robfig/cron/v3"
	"time"
)

// 此文件是task定时任务设计内容, 无须进行修改
var locolAddr = "localhost:6379"

var ctx, cancel = context.WithCancel(context.Background())

func (t *TaskServer) Start(ctx context.Context) error {
	t.scheduler.Start()
	return nil
}

// TaskType represents the type of task (case1, case2, once)

const (
	Cron = "cron"
	Once = "once"
)

// Task represents a scheduled task
type Task struct {
	Name string
	Type string
	// 有两种写法
	// 第一种: "@every 1h2m30s" 只能写到h, m, s
	// 第二种: cron表达式: "* 1 * * *"
	// 第三种: 只用于单次执行任务, 例如"4s", 代表任务添加后4s执行 只能写h, m, s
	Schedule string
	NextTime time.Time
	Handler  Job
}

type Job func()

func (t *Task) Run() {
	// 判断
	switch t.Type {
	case Once:
		// 单次执行
		t.once()
	case Cron:
		// 定时执行
		t.cron()
	default:
		panic("任务类型有误")
	}

}

func (t *Task) once() {
	// 单次执行
	go func() {
		var waitTime time.Duration
		var err error
		if t.Schedule != "" {
			waitTime, err = time.ParseDuration(t.Schedule)
			if err != nil {
				panic("once 任务类型表达式有误")
			}
		}
		if !t.NextTime.IsZero() {
			// 计算等待时间
			waitTime = t.NextTime.Sub(time.Now())
		}
		t.run(waitTime)
	}()
}

func (t *Task) run(waitTime time.Duration) error {
	// 使用time.After等待指定时间
	select {
	case <-time.After(waitTime):
		// 先抢锁
		locker := lock.NewRedisLock(t.Name, lock.WithExpireSeconds(int64(waitTime.Seconds())))
		err := locker.Lock(context.Background())
		if err != nil {
			// 抢锁失败, 直接跳过执行, 下一轮
			return nil
		}
		if t.Handler == nil {
			panic(fmt.Sprintf("请检查%s任务%s的Handler是否为空", t.Type, t.Name))
		}
		t.Handler()
		if t.Type == Once {
			locker.Unlock(context.Background())
		}
	case <-ctx.Done():
		// 终止
		return ctx.Err()
	}
	return nil
}

func (t *Task) cron() {
	// 定时执行
	go func() {
		// 解析 cron 表达式
		schedule, err := cronpkg.ParseStandard(t.Schedule)
		if err != nil {
			panic(fmt.Sprintf("请检查定时任务%s的cron表达式是否正确", t.Name))
		}
		for {
			// 获取 cron 表达式下一次执行的时间
			nextTime := schedule.Next(time.Now())
			fmt.Println(nextTime)
			// 计算等待时间
			waitTime := nextTime.Sub(time.Now())
			if err := t.run(waitTime); err != nil {
				return
			}
		}
	}()
}

// TaskScheduler represents the task scheduler
type TaskScheduler struct {
	// 定时控制
	tasks []*Task // 要执行的任务
}

// NewScheduler creates a new taskScheduler instance
func NewScheduler(addr string, tasks []*Task) *TaskScheduler {
	if addr != "" {
		locolAddr = addr
	}
	return &TaskScheduler{
		tasks: tasks,
	}
}

// AddTask adds a new task to the scheduler
func (s *TaskScheduler) AddTask(task Task) {
	task.Run()
}

// Start starts the scheduler
func (s *TaskScheduler) Start() {
	// 遍历所有任务
	for _, task := range s.tasks {
		task.Run()
	}
}

// Stop stops the scheduler
func (s *TaskScheduler) Stop() {
	cancel()
}
