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
 @Time    : 2025/4/9 -- 10:36
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: example example/main.go
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/zwxyxwz/transactions/retry"
)

func main() {
	// 模拟 stat
	counter := 0

	rty := retry.NewRetry(retry.Config{
		MaxAttempts:   5,
		InitialDelay:  200 * time.Millisecond,
		BackoffFactor: 1.5,
		Jitter:        0.3,
		RetryOn: func(err error) bool {
			return strings.Contains(err.Error(), "transient error")
		},
		OnRetry: func(attempt int, delay time.Duration, err error) {
			counter += 1
			fmt.Println(counter, delay, err)
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := rty.Do(ctx, func() error {
		_, err := http.Get("https://your.example.com")
		return err
	})

	if err != nil {
		log.Fatalf("Operation failed: %v", err)
	}
}
