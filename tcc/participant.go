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
 @Time    : 2025/4/24 -- 18:16
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: tcc tcc/participant.go
*/

package tcc

import (
	"context"
)

// Participant 表示 TCC 事务的参与者
type Participant interface {
	// GetName 获取参与者名称
	GetName() string

	// Try 尝试执行业务操作，预留资源
	Try(ctx context.Context) error

	// Confirm 确认执行业务操作
	Confirm(ctx context.Context) error

	// Cancel 取消执行业务操作，释放预留的资源
	Cancel(ctx context.Context) error
}
