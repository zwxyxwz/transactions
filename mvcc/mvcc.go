package mvcc

// MVCC is the main coordinator for MVCC operations.
type MVCC struct {
	storage    *MemoryStorage
	txnManager *TransactionManager
}

func NewMVCC() *MVCC {
	storage := NewMemoryStorage()
	txnMgr := NewTransactionManager()
	txnMgr.AttachStorage(storage)
	return &MVCC{
		storage:    storage,
		txnManager: txnMgr,
	}
}

func (m *MVCC) BeginTransaction(isolation IsolationLevel) *Transaction {
	return m.txnManager.Begin(isolation)
}

func (m *MVCC) Get(txn *Transaction, key string) ([]byte, error) {
	if txn == nil || !txn.Active {
		return nil, ErrTxnNotActive
	}
	val, _, err := m.storage.Get(key, uint64(txn.ID))
	return val, err
}

func (m *MVCC) Set(txn *Transaction, key string, value []byte) error {
	if txn == nil || !txn.Active {
		return ErrTxnNotActive
	}
	return m.storage.Set(key, value, uint64(txn.ID))
}

func (m *MVCC) Delete(txn *Transaction, key string) error {
	if txn == nil || !txn.Active {
		return ErrTxnNotActive
	}
	return m.storage.Delete(key, uint64(txn.ID))
}

func (m *MVCC) Commit(txn *Transaction) error {
	return m.txnManager.Commit(txn)
}

func (m *MVCC) Rollback(txn *Transaction) error {
	return m.txnManager.Rollback(txn)
}
