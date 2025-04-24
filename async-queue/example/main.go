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
 @Time    : 2025/4/24 -- 15:57
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: example async-queue/example/main.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	asyncqueue "github.com/zwxyxwz/transactions/async-queue"
)

func main() {
	// 创建异步队列系统
	config := asyncqueue.DefaultConfig()
	config.WorkerCount = 2
	aq := asyncqueue.New(config)

	// 注册任务处理器
	aq.RegisterHandlerFunc("email", func(ctx context.Context, t *asyncqueue.Task) (interface{}, error) {
		to, _ := t.Payload["to"].(string)
		subject, _ := t.Payload["subject"].(string)
		body, _ := t.Payload["body"].(string)

		log.Printf("发送邮件: 收件人=%s, 主题=%s", to, subject)
		// 模拟发送邮件的耗时操作
		time.Sleep(1 * time.Second)

		return map[string]interface{}{
			"sent_to":    to,
			"sent_at":    time.Now().Format(time.RFC3339),
			"message_id": fmt.Sprintf("msg-%d", time.Now().UnixNano()),
			"body":       body,
		}, nil
	})

	aq.RegisterHandlerFunc("notification", func(ctx context.Context, t *asyncqueue.Task) (interface{}, error) {
		userID, _ := t.Payload["user_id"].(string)
		message, _ := t.Payload["message"].(string)

		log.Printf("发送通知: 用户ID=%s, 消息=%s", userID, message)
		// 模拟发送通知的耗时操作
		time.Sleep(500 * time.Millisecond)

		return map[string]interface{}{
			"delivered": true,
			"timestamp": time.Now().Unix(),
		}, nil
	})

	// 启动工作者
	aq.Start(config.WorkerCount)
	defer aq.Stop()

	// 创建并入队一些任务
	ctx := context.Background()

	// 创建邮件任务
	emailTask, err := aq.NewTask(ctx, "email", map[string]interface{}{
		"to":      "user@example.com",
		"subject": "测试邮件",
		"body":    "这是一封测试邮件",
	})
	if err != nil {
		log.Fatalf("Failed to enqueue email task: %v", err)
	}
	log.Printf("邮件任务已入队: %s", emailTask.ID)

	// 创建通知任务
	notificationTask, err := aq.NewTask(ctx, "notification", map[string]interface{}{
		"user_id": "user123",
		"message": "您有一条新消息",
	})
	if err != nil {
		log.Fatalf("Failed to enqueue notification task: %v", err)
	}
	log.Printf("通知任务已入队: %s", notificationTask.ID)

	// 等待任务完成
	log.Println("等待任务完成...")
	time.Sleep(3 * time.Second)

	log.Println("所有任务已处理完成")
}
