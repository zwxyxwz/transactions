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
 @Time    : 2025/4/24 -- 15:14
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: async_queue async-queue/task.go
*/

package async_queue

import (
	"context"
	"encoding/json"
	"time"
)

// Status 表示任务的状态
type Status string

const (
	StatusPending   Status = "pending"   // 等待执行
	StatusRunning   Status = "running"   // 正在执行
	StatusCompleted Status = "completed" // 执行完成
	StatusFailed    Status = "failed"    // 执行失败
	StatusCancelled Status = "cancelled" // 已取消
)

// Priority 表示任务的优先级
type Priority int

const (
	PriorityLow      Priority = 1 // 低优先级
	PriorityNormal   Priority = 2 // 普通优先级
	PriorityHigh     Priority = 3 // 高优先级
	PriorityCritical Priority = 4 // 关键优先级
)

// Result 表示任务执行的结果
type Result struct {
	Data        interface{} `json:"data"`         // 结果数据
	Error       string      `json:"error"`        // 错误信息
	CompletedAt time.Time   `json:"completed_at"` // 完成时间
}

// Task 表示一个异步任务
type Task struct {
	ID          string                 `json:"id"`           // 任务ID
	Type        string                 `json:"type"`         // 任务类型
	Payload     map[string]interface{} `json:"payload"`      // 任务负载
	Status      Status                 `json:"status"`       // 任务状态
	Priority    Priority               `json:"priority"`     // 任务优先级
	MaxRetries  int                    `json:"max_retries"`  // 最大重试次数
	RetryCount  int                    `json:"retry_count"`  // 当前重试次数
	Result      *Result                `json:"result"`       // 任务结果
	CreatedAt   time.Time              `json:"created_at"`   // 创建时间
	StartedAt   *time.Time             `json:"started_at"`   // 开始执行时间
	CallbackURL string                 `json:"callback_url"` // 回调URL
}

// Handler 定义任务处理器接口
type Handler interface {
	// Handle 处理任务并返回结果
	Handle(ctx context.Context, task *Task) (interface{}, error)
}

// HandlerFunc 是一个适配器，允许使用普通函数作为任务处理器
type HandlerFunc func(ctx context.Context, task *Task) (interface{}, error)

// Handle 实现Handler接口
func (f HandlerFunc) Handle(ctx context.Context, task *Task) (interface{}, error) {
	return f(ctx, task)
}

// NewTask 创建一个新的任务
func NewTask(taskType string, payload map[string]interface{}) *Task {
	return &Task{
		ID:         generateID(),
		Type:       taskType,
		Payload:    payload,
		Status:     StatusPending,
		Priority:   PriorityNormal,
		MaxRetries: 3,
		RetryCount: 0,
		CreatedAt:  time.Now(),
	}
}

// SetPriority 设置任务优先级
func (t *Task) SetPriority(priority Priority) *Task {
	t.Priority = priority
	return t
}

// SetMaxRetries 设置最大重试次数
func (t *Task) SetMaxRetries(maxRetries int) *Task {
	t.MaxRetries = maxRetries
	return t
}

// SetCallbackURL 设置回调URL
func (t *Task) SetCallbackURL(callbackURL string) *Task {
	t.CallbackURL = callbackURL
	return t
}

// ToJSON 将任务转换为JSON字符串
func (t *Task) ToJSON() (string, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON 从JSON字符串解析任务
func FromJSON(data string) (*Task, error) {
	task := &Task{}
	err := json.Unmarshal([]byte(data), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// 生成唯一ID
func generateID() string {
	// 简单实现，实际应用中可以使用UUID等
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // 确保随机性
	}
	return string(b)
}
