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
 @Time    : 2025/4/24 -- 18:32
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: purse tcc/impl/purse/service.go
*/

package purse

import (
	"context"
	"fmt"
)

// Service 账户服务
type Service struct {
	purse        *Purse
	deductAmount float64
}

// NewPurseService 创建账户服务
func NewPurseService(Purse *Purse, amount float64) *Service {
	return &Service{
		purse:        Purse,
		deductAmount: amount,
	}
}

// GetName 获取参与者名称
func (s *Service) GetName() string {
	return fmt.Sprintf("PurseService-%d", s.purse.ID)
}

// Try 尝试扣款
func (s *Service) Try(ctx context.Context) error {
	fmt.Printf("账户服务 Try: 尝试从账户 %d 扣款 %.2f\n", s.purse.ID, s.deductAmount)
	return s.purse.TryDeduct(s.deductAmount)
}

// Confirm 确认扣款
func (s *Service) Confirm(ctx context.Context) error {
	fmt.Printf("账户服务 Confirm: 确认从账户 %d 扣款 %.2f\n", s.purse.ID, s.deductAmount)
	return s.purse.ConfirmDeduct(s.deductAmount)
}

// Cancel 取消扣款
func (s *Service) Cancel(ctx context.Context) error {
	fmt.Printf("账户服务 Cancel: 取消从账户 %d 扣款 %.2f\n", s.purse.ID, s.deductAmount)
	return s.purse.CancelDeduct(s.deductAmount)
}
