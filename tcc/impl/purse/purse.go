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
 @Time    : 2025/4/24 -- 18:27
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: purse tcc/impl/purse/purse.go
*/

package purse

import (
	"errors"
	"fmt"
	"sync"
)

// Purse 表示用户钱包
type Purse struct {
	ID           int64
	UserID       int64
	Balance      float64
	FrozenAmount float64
	mu           sync.Mutex
}

// NewPurse 创建一个新钱包
func NewPurse(id, userID int64, initialBalance float64) *Purse {
	return &Purse{
		ID:           id,
		UserID:       userID,
		Balance:      initialBalance,
		FrozenAmount: 0,
	}
}

// TryDeduct 尝试扣款（冻结金额）
func (a *Purse) TryDeduct(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Balance < amount {
		return errors.New("余额不足")
	}

	// 原子操作 cas 或其他方式
	a.Balance -= amount
	a.FrozenAmount += amount
	fmt.Printf("账户 %d: 冻结金额 %.2f, 当前余额 %.2f, 冻结金额 %.2f\n",
		a.ID, amount, a.Balance, a.FrozenAmount)

	return nil
}

// ConfirmDeduct 确认扣款
func (a *Purse) ConfirmDeduct(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.FrozenAmount < amount {
		return errors.New("冻结金额不足")
	}

	// 原子操作 cas 或其他方式
	a.FrozenAmount -= amount
	fmt.Printf("账户 %d: 确认扣款 %.2f, 当前余额 %.2f, 冻结金额 %.2f\n",
		a.ID, amount, a.Balance, a.FrozenAmount)

	return nil
}

// CancelDeduct 取消扣款
func (a *Purse) CancelDeduct(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.FrozenAmount < amount {
		return errors.New("冻结金额不足")
	}

	// 原子操作 cas 或其他方式
	a.Balance += amount
	a.FrozenAmount -= amount
	fmt.Printf("账户 %d: 取消扣款 %.2f, 当前余额 %.2f, 冻结金额 %.2f\n",
		a.ID, amount, a.Balance, a.FrozenAmount)

	return nil
}
