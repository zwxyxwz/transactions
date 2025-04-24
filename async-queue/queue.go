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
 @Time    : 2025/4/24 -- 15:50
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: async_queue async-queue/queue.go
*/

package async_queue

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Queue 定义队列接口
type Queue interface {
	// Push 将任务推入队列
	Push(ctx context.Context, task *Task) error
	// Pop 从队列中取出任务
	Pop(ctx context.Context) (*Task, error)
	// Size 返回队列中的任务数量
	Size() int
	// Close 关闭队列
	Close() error
}

// MemoryQueue 是一个基于内存的队列实现
type MemoryQueue struct {
	tasks     []*Task
	mutex     sync.Mutex
	cond      *sync.Cond
	closed    bool
	maxSize   int
	waitTime  time.Duration
	scheduler Scheduler
}

// Scheduler 定义任务调度器接口
type Scheduler interface {
	// Schedule 对任务进行排序
	Schedule(tasks []*Task) []*Task
}

// FIFOScheduler 是一个先进先出的调度器
type FIFOScheduler struct{}

// Schedule 按照先进先出的顺序排序任务
func (s *FIFOScheduler) Schedule(tasks []*Task) []*Task {
	return tasks // 先进先出，不需要额外排序
}

// PriorityScheduler 是一个基于优先级的调度器
type PriorityScheduler struct{}

// Schedule 按照优先级排序任务
func (s *PriorityScheduler) Schedule(tasks []*Task) []*Task {
	// 复制任务列表，避免修改原始列表
	result := make([]*Task, len(tasks))
	copy(result, tasks)

	// 按优先级排序（高优先级在前）
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Priority < result[j].Priority {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// NewMemoryQueue 创建一个新的内存队列
func NewMemoryQueue(maxSize int, waitTime time.Duration, scheduler Scheduler) *MemoryQueue {
	q := &MemoryQueue{
		tasks:     make([]*Task, 0),
		maxSize:   maxSize,
		waitTime:  waitTime,
		scheduler: scheduler,
	}
	q.cond = sync.NewCond(&q.mutex)
	return q
}

// Push 将任务推入队列
func (q *MemoryQueue) Push(ctx context.Context, t *Task) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.closed {
		return errors.New("queue is closed")
	}

	if q.maxSize > 0 && len(q.tasks) >= q.maxSize {
		return errors.New("queue is full")
	}

	q.tasks = append(q.tasks, t)
	q.cond.Signal() // 通知等待的消费者

	return nil
}

// Pop 从队列中取出任务
func (q *MemoryQueue) Pop(ctx context.Context) (*Task, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// 如果队列为空，等待任务到达或超时
	if len(q.tasks) == 0 {
		if q.closed {
			return nil, errors.New("queue is closed")
		}

		// 设置等待超时
		waitCh := make(chan struct{})
		go func() {
			q.cond.Wait()
			close(waitCh)
		}()

		q.mutex.Unlock()
		select {
		case <-waitCh:
			// 重新获取锁
			q.mutex.Lock()
		case <-time.After(q.waitTime):
			// 超时
			q.mutex.Lock()
			return nil, errors.New("timeout waiting for task")
		case <-ctx.Done():
			// 上下文取消
			q.mutex.Lock()
			return nil, ctx.Err()
		}

		// 再次检查队列状态
		if len(q.tasks) == 0 {
			if q.closed {
				return nil, errors.New("queue is closed")
			}
			return nil, errors.New("no task available")
		}
	}

	// 使用调度器对任务进行排序
	sortedTasks := q.scheduler.Schedule(q.tasks)

	// 取出第一个任务
	t := sortedTasks[0]

	// 从原始队列中移除该任务
	for i, task := range q.tasks {
		if task.ID == t.ID {
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
			break
		}
	}

	return t, nil
}

// Size 返回队列中的任务数量
func (q *MemoryQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.tasks)
}

// Close 关闭队列
func (q *MemoryQueue) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.closed {
		return errors.New("queue already closed")
	}

	q.closed = true
	q.cond.Broadcast() // 通知所有等待的消费者
	return nil
}
