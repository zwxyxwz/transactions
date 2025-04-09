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
 @Time    : 2025/4/9 -- 10:35
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: retry /retry.go
*/

package retry

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Config struct {
	MaxAttempts   int                                               // 最大重试次数
	InitialDelay  time.Duration                                     // 初始延迟
	MaxDelay      time.Duration                                     // 最大延迟
	BackoffFactor float64                                           // 退避因子（默认2倍）
	Jitter        float64                                           // 抖动系数（0-1之间）
	RetryOn       func(error) bool                                  // 错误过滤条件
	OnRetry       func(attempt int, delay time.Duration, err error) // 回调函数
}

type Retry struct {
	config Config
	logger *log.Logger
}

func NewRetry(config Config) *Retry {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = 3
	}
	if config.InitialDelay == 0 {
		config.InitialDelay = 100 * time.Millisecond
	}
	return &Retry{
		config: config,
		logger: log.New(os.Stderr, "[retry] ", log.LstdFlags),
	}
}

func (r *Retry) Do(ctx context.Context, operation func() error) error {
	var attempt int
	var delay time.Duration

	for {
		err := operation()
		if err == nil {
			return nil
		}

		if !r.shouldRetry(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		attempt++
		if attempt >= r.config.MaxAttempts {
			return fmt.Errorf("max attempts (%d) exceeded: last error: %w",
				r.config.MaxAttempts, err)
		}

		delay = r.calculateDelay(attempt)
		r.logger.Printf("retry %d/%d after %v (error: %v)",
			attempt, r.config.MaxAttempts, delay, err)

		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return fmt.Errorf("retry canceled: %w", ctx.Err())
		}
	}
}

func (r *Retry) shouldRetry(err error) bool {
	if r.config.RetryOn == nil {
		return true
	}
	return r.config.RetryOn(err)
}

func (r *Retry) calculateDelay(attempt int) time.Duration {
	baseDelay := r.config.InitialDelay * time.Duration(math.Pow(r.config.BackoffFactor, float64(attempt-1)))
	jitter := time.Duration(rand.Float64() * float64(baseDelay) * r.config.Jitter)
	delay := baseDelay + jitter

	if r.config.MaxDelay > 0 {
		if delay > r.config.MaxDelay {
			return r.config.MaxDelay
		}
	}
	return delay
}
