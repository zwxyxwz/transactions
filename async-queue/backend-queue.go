/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2025/4/14 -- 18:35
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: async_queue async-queue/backend-queue.go
*/

package async_queue

import (
	"context"
	"time"
)

// AsyncQueue 是异步队列系统的主要入口
type AsyncQueue struct {
	queue    Queue
	worker   *Worker
	callback Callback
}

// Config 配置异步队列系统
type Config struct {
	QueueMaxSize    int
	QueueWaitTime   time.Duration
	WorkerCount     int
	MaxRetries      int
	CallbackTimeout time.Duration
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		QueueMaxSize:    1000,
		QueueWaitTime:   5 * time.Second,
		WorkerCount:     5,
		MaxRetries:      3,
		CallbackTimeout: 10 * time.Second,
	}
}

// New 创建一个新的异步队列系统
func New(config *Config) *AsyncQueue {
	if config == nil {
		config = DefaultConfig()
	}

	// 创建调度器
	scheduler := &PriorityScheduler{}

	// 创建队列
	q := NewMemoryQueue(config.QueueMaxSize, config.QueueWaitTime, scheduler)

	// 创建回调
	cb := NewHTTPCallback(config.CallbackTimeout)

	// 创建工作者
	w := NewWorker("default", q, cb)
	w.SetMaxRetries(config.MaxRetries)

	return &AsyncQueue{
		queue:    q,
		worker:   w,
		callback: cb,
	}
}

// Start 启动异步队列系统
func (aq *AsyncQueue) Start(workerCount int) {
	aq.worker.Start(workerCount)
}

// Stop 停止异步队列系统
func (aq *AsyncQueue) Stop() {
	aq.worker.Stop()
	aq.queue.Close()
}

// RegisterHandler 注册任务处理器
func (aq *AsyncQueue) RegisterHandler(taskType string, handler Handler) {
	aq.worker.RegisterHandler(taskType, handler)
}

// RegisterHandlerFunc 注册任务处理函数
func (aq *AsyncQueue) RegisterHandlerFunc(taskType string, handlerFunc HandlerFunc) {
	aq.worker.RegisterHandlerFunc(taskType, handlerFunc)
}

// EnqueueTask 将任务加入队列
func (aq *AsyncQueue) EnqueueTask(ctx context.Context, t *Task) error {
	return aq.queue.Push(ctx, t)
}

// NewTask 创建一个新任务并加入队列
func (aq *AsyncQueue) NewTask(ctx context.Context, taskType string, payload map[string]interface{}) (*Task, error) {
	t := NewTask(taskType, payload)
	err := aq.EnqueueTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// QueueSize 返回队列中的任务数量
func (aq *AsyncQueue) QueueSize() int {
	return aq.queue.Size()
}
