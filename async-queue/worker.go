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
 @Time    : 2025/4/24 -- 15:46
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: async_queue async-queue/worker.go
*/

package async_queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Worker 表示一个工作者
type Worker struct {
	id         string
	queue      Queue
	handlers   map[string]Handler
	callback   Callback
	maxRetries int
	wg         sync.WaitGroup
	stopCh     chan struct{}
	stopped    bool
	mutex      sync.Mutex
}

// NewWorker 创建一个新的工作者
func NewWorker(id string, q Queue, cb Callback) *Worker {
	return &Worker{
		id:         id,
		queue:      q,
		handlers:   make(map[string]Handler),
		callback:   cb,
		maxRetries: 3,
		stopCh:     make(chan struct{}),
	}
}

// RegisterHandler 注册任务处理器
func (w *Worker) RegisterHandler(taskType string, handler Handler) {
	w.handlers[taskType] = handler
}

// RegisterHandlerFunc 注册任务处理函数
func (w *Worker) RegisterHandlerFunc(taskType string, handlerFunc HandlerFunc) {
	w.handlers[taskType] = handlerFunc
}

// SetMaxRetries 设置最大重试次数
func (w *Worker) SetMaxRetries(maxRetries int) {
	w.maxRetries = maxRetries
}

// Start 启动工作者
func (w *Worker) Start(numWorkers int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.stopped {
		return
	}

	for i := 0; i < numWorkers; i++ {
		w.wg.Add(1)
		go w.processLoop(fmt.Sprintf("%s-%d", w.id, i))
	}
}

// Stop 停止工作者
func (w *Worker) Stop() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.stopped {
		return
	}

	w.stopped = true
	close(w.stopCh)
	w.wg.Wait()
}

// processLoop 处理任务循环
func (w *Worker) processLoop(workerId string) {
	defer w.wg.Done()

	log.Printf("Worker %s started", workerId)

	for {
		select {
		case <-w.stopCh:
			log.Printf("Worker %s stopped", workerId)
			return
		default:
			// 创建一个带超时的上下文
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			t, err := w.queue.Pop(ctx)
			cancel()

			if err != nil {
				// 如果是因为队列为空或超时，则继续尝试
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// 处理任务
			w.processTask(t)
		}
	}
}

// processTask 处理单个任务
func (w *Worker) processTask(t *Task) {
	// 更新任务状态
	t.Status = StatusRunning
	now := time.Now()
	t.StartedAt = &now

	// 查找处理器
	handler, ok := w.handlers[t.Type]
	if !ok {
		log.Printf("No handler registered for task type: %s", t.Type)
		t.Status = StatusFailed
		t.Result = &Result{
			Error:       fmt.Sprintf("no handler for task type: %s", t.Type),
			CompletedAt: time.Now(),
		}
		w.handleResult(t)
		return
	}

	// 执行任务
	log.Printf("Processing task %s of type %s", t.ID, t.Type)
	result, err := handler.Handle(context.Background(), t)

	// 处理结果
	if err != nil {
		log.Printf("Task %s failed: %v", t.ID, err)

		// 检查是否需要重试
		if t.RetryCount < t.MaxRetries {
			t.RetryCount++
			t.Status = StatusPending
			log.Printf("Retrying task %s (%d/%d)", t.ID, t.RetryCount, t.MaxRetries)

			// 重新入队
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := w.queue.Push(ctx, t); err != nil {
				log.Printf("Failed to requeue task %s: %v", t.ID, err)
				t.Status = StatusFailed
				t.Result = &Result{
					Error:       fmt.Sprintf("failed to requeue: %v", err),
					CompletedAt: time.Now(),
				}
				w.handleResult(t)
			}
			return
		}

		// 达到最大重试次数，标记为失败
		t.Status = StatusFailed
		t.Result = &Result{
			Error:       err.Error(),
			CompletedAt: time.Now(),
		}
	} else {
		// 任务成功完成
		log.Printf("Task %s completed successfully", t.ID)
		t.Status = StatusCompleted
		t.Result = &Result{
			Data:        result,
			CompletedAt: time.Now(),
		}
	}

	// 处理结果回调
	w.handleResult(t)
}

// handleResult 处理任务结果
func (w *Worker) handleResult(t *Task) {
	if w.callback != nil {
		if err := w.callback.OnTaskComplete(t); err != nil {
			log.Printf("Failed to process callback for task %s: %v", t.ID, err)
		}
	}
}
