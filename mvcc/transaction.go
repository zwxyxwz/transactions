package mvcc

import (
	"sync"
	"time"
)

// Transaction represents a database transaction.
type Transaction struct {
	ID        TxnID
	Isolation IsolationLevel
	Active    bool
	BeginTS   int64
	CommitTS  int64
}

// TransactionManager manages all transactions.
type TransactionManager struct {
	mu           sync.Mutex
	nextTxnID    TxnID
	transactions map[TxnID]*Transaction
	storage      *MemoryStorage // 新增字段用于访问存储
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		nextTxnID:    1,
		transactions: make(map[TxnID]*Transaction),
	}
}

func (tm *TransactionManager) AttachStorage(storage *MemoryStorage) {
	tm.storage = storage
}

func (tm *TransactionManager) Begin(isolation IsolationLevel) *Transaction {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	txn := &Transaction{
		ID:        tm.nextTxnID,
		Isolation: isolation,
		Active:    true,
		BeginTS:   nowTS(),
	}
	tm.transactions[txn.ID] = txn
	tm.nextTxnID++
	return txn
}

func (tm *TransactionManager) Commit(txn *Transaction) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if txn == nil {
		return ErrTxnNotFound
	}
	if !txn.Active {
		return ErrTxnNotActive
	}
	txn.Active = false
	txn.CommitTS = nowTS()
	// 修正：遍历所有数据，将属于该事务的未提交版本标记为已提交
	if tm.storage != nil {
		tm.storage.mu.Lock()
		for _, versions := range tm.storage.data {
			for _, v := range versions {
				if v.TxnID == uint64(txn.ID) && !v.Committed {
					v.Committed = true
				}
			}
		}
		tm.storage.mu.Unlock()
	}
	return nil
}

func (tm *TransactionManager) Rollback(txn *Transaction) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if txn == nil {
		return ErrTxnNotFound
	}
	if !txn.Active {
		return ErrTxnNotActive
	}
	txn.Active = false
	txn.CommitTS = 0
	return nil
}

func nowTS() int64 {
	return timeNowUnixNano()
}

func timeNowUnixNano() int64 {
	return time.Now().UnixNano()
}
