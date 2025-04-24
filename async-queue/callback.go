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
 @Time    : 2025/4/24 -- 15:48
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: async_queue async-queue/callback.go
*/

package async_queue

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Callback 定义回调接口
type Callback interface {
	// OnTaskComplete 在任务完成时调用
	OnTaskComplete(task *Task) error
}

// HTTPCallback 是一个基于HTTP的回调实现
type HTTPCallback struct {
	client *http.Client
}

// NewHTTPCallback 创建一个新的HTTP回调
func NewHTTPCallback(timeout time.Duration) *HTTPCallback {
	return &HTTPCallback{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// OnTaskComplete 在任务完成时发送HTTP回调
func (c *HTTPCallback) OnTaskComplete(t *Task) error {
	// 如果没有设置回调URL，则跳过
	if t.CallbackURL == "" {
		return nil
	}

	// 准备回调数据
	callbackData := map[string]interface{}{
		"task_id":      t.ID,
		"task_type":    t.Type,
		"status":       t.Status,
		"result":       t.Result,
		"completed_at": t.Result.CompletedAt,
	}

	// 转换为JSON
	jsonData, err := json.Marshal(callbackData)
	if err != nil {
		return fmt.Errorf("failed to marshal callback data: %w", err)
	}

	// 发送HTTP请求
	resp, err := c.client.Post(t.CallbackURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send callback: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("callback failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// FunctionCallback 是一个基于函数的回调实现
type FunctionCallback struct {
	callback func(*Task) error
}

// NewFunctionCallback 创建一个新的函数回调
func NewFunctionCallback(callback func(*Task) error) *FunctionCallback {
	return &FunctionCallback{
		callback: callback,
	}
}

// OnTaskComplete 在任务完成时调用回调函数
func (c *FunctionCallback) OnTaskComplete(t *Task) error {
	if c.callback == nil {
		return errors.New("callback function is nil")
	}
	return c.callback(t)
}

// CompositeCallback 组合多个回调
type CompositeCallback struct {
	callbacks []Callback
}

// NewCompositeCallback 创建一个新的组合回调
func NewCompositeCallback(callbacks ...Callback) *CompositeCallback {
	return &CompositeCallback{
		callbacks: callbacks,
	}
}

// OnTaskComplete 在任务完成时调用所有回调
func (c *CompositeCallback) OnTaskComplete(t *Task) error {
	var lastErr error
	for _, callback := range c.callbacks {
		if err := callback.OnTaskComplete(t); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
