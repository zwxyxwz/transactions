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
 @Description: order tcc/impl/order/service.go
*/

package order

import (
	"context"
	"fmt"
)

// Service 订单服务
type Service struct {
	order *Order
}

// NewService 创建订单服务
func NewService(order *Order) *Service {
	return &Service{
		order: order,
	}
}

// GetName 获取参与者名称
func (s *Service) GetName() string {
	return fmt.Sprintf("Service-%d", s.order.ID)
}

// Try 尝试创建订单
func (s *Service) Try(ctx context.Context) error {
	fmt.Printf("订单服务 Try: 尝试创建订单 %d\n", s.order.ID)
	return s.order.TryCreate()
}

// Confirm 确认创建订单
func (s *Service) Confirm(ctx context.Context) error {
	fmt.Printf("订单服务 Confirm: 确认创建订单 %d\n", s.order.ID)
	return s.order.ConfirmCreate()
}

// Cancel 取消创建订单
func (s *Service) Cancel(ctx context.Context) error {
	fmt.Printf("订单服务 Cancel: 取消创建订单 %d\n", s.order.ID)
	return s.order.CancelCreate()
}
