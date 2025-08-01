package spider

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/panjf2000/ants/v2"
	"harmel.cn/financial/internal/public"
)

// 爬虫管理器
type SpiderManager struct {
	// 进度管理器
	progressManager *ProgressManager
	// 待处理任务通道
	taskChan chan PendingTask
	// 线程池
	pool *ants.Pool
	// 是否有在执行的任务
	hasTaskRunning atomic.Bool
	// 完成通知
	finishNotify chan bool
}

func NewSpiderManager(rootDir string) *SpiderManager {
	pool, err := ants.NewPool(public.SpiderExecutorPoolSize, ants.WithPreAlloc(true))
	if err != nil {
		panic(err)
	}

	return &SpiderManager{
		progressManager: NewProgressManager(rootDir),
		taskChan:        make(chan PendingTask, PENDING_TASKS_INIT_CAPACITY),
		pool:            pool,
		finishNotify:    make(chan bool),
	}
}

// 开启爬虫
func (s *SpiderManager) Start(ctx context.Context) (err error) {
	g.Log("spider").Debugf(ctx, "spider is running")

	// 加载历史进度
	err = s.progressManager.Load(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "ProgressManager.Load failed, err is %v", err)
		return err
	}

	// 如果到了五月一日，清空所有任务全部重跑（年报全部出了）
	if time.Now().Format("01-02") == "05-01" {
		s.progressManager.ClearTasks()
	}

	// 如果上次成功了，判断时间是否大于等于配置天数
	if s.progressManager.Done() {
		if time.Now().Unix()-s.progressManager.LastTS() >= public.SpiderTaskIntervalDays*24*3600 {
			s.progressManager.ClearTasks()
		} else {
			g.Log("spider").Debugf(ctx, "task finish, the time since the last successful task completion is less than %d days", public.SpiderTaskIntervalDays)
			return
		}
	}

	// 基础数据
	err = s.fetchIndexSample()
	if err != nil {
		g.Log("spider").Errorf(ctx, "fetch index sample data failed, err is %v", err)
		return
	}
	err = s.fetchCategory()
	if err != nil {
		g.Log("spider").Errorf(ctx, "fetch category data failed, err is %v", err)
		return
	}

	// 启动处理任务线程
	go s.doProcTaskWorker(ctx)

	// 等待完成
OUT:
	for {
		select {
		case <-s.finishNotify:
			g.Log("spider").Debug(ctx, "all task execute finish")
			break OUT
		default:
			if !s.hasTaskRunning.Load() {
				tasks := s.progressManager.UnexecutedTasks()
				if len(tasks) == 0 {
					s.progressManager.SetDone()
					err = s.progressManager.Save(ctx)
					if err != nil {
						g.Log("spider").Errorf(ctx, "save process failed, err is %v", err)
					}
					return
				}
				for _, task := range tasks {
					s.taskChan <- task
				}
			}
		}
	}

	err = s.calcFinancialRatios()
	if err != nil {
		g.Log("spider").Errorf(ctx, "calc financial ratios failed, err is %v", err)
		return
	}

	return
}

// 最新指数样本信息
func (s *SpiderManager) fetchIndexSample() error {
	// TODO
	return nil
}

// 最新行业分类信息（含分类下的股票）
func (s *SpiderManager) fetchCategory() error {
	// TODO
	task := PendingTask{Id: "xxx"}
	s.progressManager.PutTask(task)
	s.taskChan <- task
	return nil
}

// 计算财务比率
func (s *SpiderManager) calcFinancialRatios() error {
	// TODO
	return nil
}

// 处理任务
func (s *SpiderManager) doProcTaskWorker(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			g.Log("spiser").Critical(ctx, "doProcTaskWorker panic: %v", err)
		}
	}()

	for {
		task := <-s.taskChan
		err := s.executeTask(ctx, task)
		if err != nil {
			g.Log("spider").Errorf(ctx, "execute task failed, err is %v", err)
		}
	}
}

// 执行实际任务
func (s *SpiderManager) executeTask(ctx context.Context, task PendingTask) (err error) {
	err = s.pool.Submit(func() {
		s.hasTaskRunning.Store(true)
		defer func() {
			s.hasTaskRunning.Store(false)
		}()

		// 是否处理完
		isFinished, err := s.progressManager.TaskStatus(ctx, task.Id)
		if err != nil {
			g.Log("spider").Errorf(ctx, "query task status failed, err is %v", err)
			return
		}
		if isFinished {
			return
		}

		g.Log("spider").Debugf(ctx, "start execute task %s", task.Id)

		// 标记完成
		s.progressManager.MarkTask(ctx, task.Id, true)

		// 写入磁盘
		err = s.progressManager.Save(ctx)
		if err != nil {
			g.Log("spider").Errorf(ctx, "save process failed, err is %v", err)
			return
		}

		// 通知完成
		if s.progressManager.Done() {
			s.finishNotify <- true
		}
	})
	return
}
