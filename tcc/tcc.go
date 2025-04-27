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
 @Time    : 2025/4/24 -- 18:17
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: tcc tcc/tcc.go
*/

package tcc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Status 表示事务的状态
type Status int

const (
	// StatusInit 初始状态
	StatusInit Status = iota
	// StatusTried Try 阶段完成
	StatusTried
	// StatusConfirmed Confirm 阶段完成
	StatusConfirmed
	// StatusCancelled Cancel 阶段完成
	StatusCancelled
)

// Tcc 表示一个 TCC 事务
type Tcc struct {
	ID            string
	Status        Status
	Participants  []Participant
	StartTime     time.Time
	CompletedTime time.Time
	Context       context.Context
	mu            sync.Mutex
}

// NewTcc 创建一个新的 TCC 事务
func NewTcc(ctx context.Context) *Tcc {
	return &Tcc{
		ID:           uuid.New().String(),
		Status:       StatusInit,
		Participants: make([]Participant, 0),
		StartTime:    time.Now(),
		Context:      ctx,
	}
}

// RegisterParticipant 注册一个事务参与者
func (t *Tcc) RegisterParticipant(p Participant) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Participants = append(t.Participants, p)
}

// Try 执行所有参与者的 Try 阶段
func (t *Tcc) Try() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Status != StatusInit {
		return errors.New("tcc is not ready to try")
	}

	// 执行所有参与者的 Try 方法
	for _, p := range t.Participants {
		if err := p.Try(t.Context); err != nil {
			// 如果任何一个 Try 失败，则取消整个事务
			fmt.Printf("Try failed for participant %s: %v\n", p.GetName(), err)
			_ = t.cancel()
			return err
		}
	}

	t.Status = StatusTried
	return nil
}

// Confirm 执行所有参与者的 Confirm 阶段
func (t *Tcc) Confirm() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.Status != StatusTried {
		return errors.New("tcc is not ready to confirm")
	}

	// 执行所有参与者的 Confirm 方法
	for _, p := range t.Participants {
		if err := p.Confirm(t.Context); err != nil {
			fmt.Printf("Confirm failed for participant %s: %v\n", p.GetName(), err)
			// 确认阶段的错误通常需要人工干预
			return err
		}
	}

	t.Status = StatusConfirmed
	t.CompletedTime = time.Now()
	return nil
}

// Cancel 执行所有参与者的 Cancel 阶段
func (t *Tcc) Cancel() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cancel()
}

// cancel 内部取消方法
func (t *Tcc) cancel() error {
	if t.Status == StatusConfirmed || t.Status == StatusCancelled {
		return errors.New("tcc cannot be cancelled in" + t.GetStatus())
	}

	// 执行所有参与者的 Cancel 方法
	var lastErr error
	for i := len(t.Participants) - 1; i >= 0; i-- {
		p := t.Participants[i]
		if err := p.Cancel(t.Context); err != nil {
			fmt.Printf("Cancel failed for participant %s: %v\n", p.GetName(), err)
			lastErr = err
			// 继续尝试取消其他参与者
		}
	}

	t.Status = StatusCancelled
	t.CompletedTime = time.Now()
	return lastErr
}

// GetStatus 获取事务状态的字符串表示
func (t *Tcc) GetStatus() string {
	switch t.Status {
	case StatusInit:
		return "初始化"
	case StatusTried:
		return "Try 阶段完成"
	case StatusConfirmed:
		return "已确认"
	case StatusCancelled:
		return "已取消"
	default:
		return "未知状态"
	}
}
