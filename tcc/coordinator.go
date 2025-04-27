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
 @Time    : 2025/4/24 -- 18:20
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: tcc tcc/coordinator.go
*/

package tcc

import (
	"context"
	"fmt"
	"sync"
)

// Coordinator TCC 事务协调器
type Coordinator struct {
	TccMap map[string]*Tcc
	mu     sync.RWMutex
}

// NewCoordinator 创建一个新的 TCC 协调器
func NewCoordinator() *Coordinator {
	return &Coordinator{
		TccMap: make(map[string]*Tcc),
	}
}

// Begin 开始一个新的 TCC 事务
func (c *Coordinator) Begin(ctx context.Context) *Tcc {
	tx := NewTcc(ctx)

	c.mu.Lock()
	c.TccMap[tx.ID] = tx
	c.mu.Unlock()

	return tx
}

// GetTcc 获取指定 ID 的事务
func (c *Coordinator) GetTcc(id string) (*Tcc, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	tx, ok := c.TccMap[id]
	return tx, ok
}

// Execute 执行 TCC 事务
func (c *Coordinator) Execute(tx *Tcc) error {
	fmt.Printf("开始执行事务 %s\n", tx.ID)

	// 执行 Try 阶段
	fmt.Println("执行 Try 阶段...")
	if err := tx.Try(); err != nil {
		fmt.Printf("Try 阶段失败: %v\n", err)
		return err
	}

	// 执行 Confirm 阶段
	fmt.Println("执行 Confirm 阶段...")
	if err := tx.Confirm(); err != nil {
		fmt.Printf("Confirm 阶段失败: %v\n", err)
		// 在实际应用中，Confirm 失败通常需要人工干预
		return err
	}

	fmt.Printf("事务 %s 执行成功\n", tx.ID)
	return nil
}
