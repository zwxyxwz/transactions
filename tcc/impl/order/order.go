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
 @Time    : 2025/4/24 -- 18:34
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: order tcc/impl/order/order.go
*/

package order

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// OrderStatus 表示订单状态
type OrderStatus int

const (
	// StatusCreated 已创建
	StatusCreated OrderStatus = iota
	// StatusPending 处理中
	StatusPending
	// StatusPaid 已支付
	StatusPaid
	// StatusCancelled 已取消
	StatusCancelled
)

// Order 表示订单
type Order struct {
	ID         int64
	UserID     int64
	Amount     float64
	Status     OrderStatus
	CreateTime time.Time
	UpdateTime time.Time
	mu         sync.Mutex
}

// NewOrder 创建一个新订单
func NewOrder(id, userID int64, amount float64) *Order {
	now := time.Now()
	return &Order{
		ID:         id,
		UserID:     userID,
		Amount:     amount,
		Status:     StatusCreated,
		CreateTime: now,
		UpdateTime: now,
	}
}

// TryCreate 尝试创建订单
func (o *Order) TryCreate() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.Status != StatusCreated {
		return errors.New("订单状态错误")
	}

	o.Status = StatusPending
	o.UpdateTime = time.Now()
	fmt.Printf("订单 %d: 尝试创建, 金额 %.2f, 状态: 处理中\n", o.ID, o.Amount)

	return nil
}

// ConfirmCreate 确认创建订单
func (o *Order) ConfirmCreate() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.Status != StatusPending {
		return errors.New("订单状态错误")
	}

	o.Status = StatusPaid
	o.UpdateTime = time.Now()
	fmt.Printf("订单 %d: 确认创建, 金额 %.2f, 状态: 已支付\n", o.ID, o.Amount)

	return nil
}

// CancelCreate 取消创建订单
func (o *Order) CancelCreate() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.Status != StatusPending && o.Status != StatusCreated {
		return errors.New("订单状态错误")
	}

	o.Status = StatusCancelled
	o.UpdateTime = time.Now()
	fmt.Printf("订单 %d: 取消创建, 金额 %.2f, 状态: 已取消\n", o.ID, o.Amount)

	return nil
}
