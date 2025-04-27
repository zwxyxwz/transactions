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
 @Time    : 2025/4/24 -- 18:26
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: example tcc/example/main.go
*/

package main

import (
	"context"
	"fmt"
	"github.com/zwxyxwz/transactions/tcc"
	"github.com/zwxyxwz/transactions/tcc/impl/order"
	"github.com/zwxyxwz/transactions/tcc/impl/purse"
)

func main() {
	// 创建 TCC 协调器
	coordinator := tcc.NewCoordinator()

	// 创建账户和订单
	account := purse.NewPurse(1, 1001, 1000.0)
	od := order.NewOrder(1, 1001, 200.0)

	fmt.Println("=== 成功场景演示 ===")
	demoSuccessCase(coordinator, account, od)

	// 重置账户和订单
	account = purse.NewPurse(2, 1002, 100.0) // 余额不足
	od = order.NewOrder(2, 1002, 200.0)

	fmt.Println("\n=== 失败场景演示 ===")
	demoFailureCase(coordinator, account, od)
}

// 演示成功场景
func demoSuccessCase(coordinator *tcc.Coordinator, account *purse.Purse, od *order.Order) {
	ctx := context.Background()

	// 开始事务
	tx := coordinator.Begin(ctx)

	// 注册参与者
	accountService := purse.NewPurseService(account, od.Amount)
	orderService := order.NewService(od)

	tx.RegisterParticipant(accountService)
	tx.RegisterParticipant(orderService)

	// 执行事务
	err := coordinator.Execute(tx)
	if err != nil {
		fmt.Printf("事务执行失败: %v\n", err)
		return
	}

	fmt.Printf("事务 %s 执行成功, 状态: %s\n", tx.ID, tx.GetStatus())
	fmt.Printf("账户余额: %.2f, 冻结金额: %.2f\n", account.Balance, account.FrozenAmount)
	fmt.Printf("订单状态: %d\n", od.Status)
}

// 演示失败场景
func demoFailureCase(coordinator *tcc.Coordinator, account *purse.Purse, od *order.Order) {
	ctx := context.Background()

	// 开始事务
	tx := coordinator.Begin(ctx)

	// 注册参与者
	accountService := purse.NewPurseService(account, od.Amount)
	orderService := order.NewService(od)

	tx.RegisterParticipant(accountService)
	tx.RegisterParticipant(orderService)

	// 执行事务
	err := coordinator.Execute(tx)
	if err != nil {
		fmt.Printf("事务执行失败: %v\n", err)
		fmt.Printf("事务 %s 状态: %s\n", tx.ID, tx.GetStatus())
		fmt.Printf("账户余额: %.2f, 冻结金额: %.2f\n", account.Balance, account.FrozenAmount)
		fmt.Printf("订单状态: %d\n", od.Status)
		return
	}
}
