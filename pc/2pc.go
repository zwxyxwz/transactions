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
 @Time    : 2025/4/9 -- 17:32
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2025 亓官竹
 @Description: pc23 pc23/2pc.go
*/

package pc

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// Participant ...
type Participant interface {
	Prepare(ctx context.Context) (bool, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// Coordinator ...
type Coordinator interface {
	RegisterParticipant(id string, p Participant)
	ExecuteTransaction(ctx context.Context, operation func() error) error
}

type DefaultCoordinator struct {
	participants map[string]Participant
	store        Storage
	timeout      time.Duration
}

func NewCoordinator(timeout time.Duration, store Storage) *DefaultCoordinator {
	return &DefaultCoordinator{
		participants: make(map[string]Participant),
		store:        store,
		timeout:      timeout,
	}
}

func (c *DefaultCoordinator) ExecuteTransaction(ctx context.Context, operation func() error) error {
	// 阶段1: 准备阶段
	prepareResults := make(chan struct{}, len(c.participants))
	go func() {
		for _, p := range c.participants {
			go func(p Participant) {
				ok, err := p.Prepare(ctx)
				if err != nil || !ok {
					prepareResults <- struct{}{}
					return
				}
				prepareResults <- nil
			}(p)
		}
	}()

	// 超时控制
	select {
	case <-time.After(c.timeout):
		c.store.Rollback()
		return errors.New("transaction timeout")
	case <-notifyAll(prepareResults):
		// 阶段2: 提交阶段
		commitResults := make(chan error, len(c.participants))
		for _, p := range c.participants {
			go func(p Participant) {
				commitResults <- p.Commit(ctx)
			}(p)
		}

		for range c.participants {
			if err := <-commitResults; err != nil {
				c.store.Rollback()
				return err
			}
		}
		return nil
	}
}

type DatabaseParticipant struct {
	db *sql.DB
}

func (p *DatabaseParticipant) Prepare(ctx context.Context) (bool, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	// 执行预提交操作
	_, err = tx.ExecContext(ctx, "PREPARE TRANSACTION")
	return err == nil, nil
}

func (p *DatabaseParticipant) Commit(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return errors.New("invalid transaction context")
	}
	return tx.Commit()
}

func (p *DatabaseParticipant) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return errors.New("invalid transaction context")
	}
	return tx.Rollback()
}
