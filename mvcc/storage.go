package mvcc

import (
	"sync"
)

// Storage defines the interface for MVCC storage backend.
type Storage interface {
	Get(key string, txnID uint64) (value []byte, version uint64, err error)
	Set(key string, value []byte, txnID uint64) error
	Delete(key string, txnID uint64) error
	Scan(start, end string, txnID uint64) (map[string][]byte, error)
}

// MemoryStorage is an in-memory implementation of Storage.
type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string][]*VersionedValue
}

// VersionedValue represents a value with version info.
type VersionedValue struct {
	Value     []byte
	TxnID     uint64
	Committed bool
	Timestamp int64
	Deleted   bool
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]*VersionedValue),
	}
}

func (ms *MemoryStorage) Get(key string, txnID uint64) (value []byte, version uint64, err error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	versions, ok := ms.data[key]
	if !ok || len(versions) == 0 {
		return nil, 0, ErrKeyNotFound
	}
	// 优先返回本事务未提交版本，否则返回最新已提交版本
	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i]
		if v.TxnID == txnID {
			if v.Deleted {
				return nil, 0, ErrKeyNotFound
			}
			return v.Value, v.TxnID, nil
		}
	}
	for i := len(versions) - 1; i >= 0; i-- {
		v := versions[i]
		if v.Committed {
			if v.Deleted {
				return nil, 0, ErrKeyNotFound
			}
			return v.Value, v.TxnID, nil
		}
	}
	return nil, 0, ErrKeyNotFound
}

func (ms *MemoryStorage) Set(key string, value []byte, txnID uint64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	v := &VersionedValue{
		Value:     value,
		TxnID:     txnID,
		Committed: false,
		Timestamp: timeNowUnixNano(),
		Deleted:   false,
	}
	ms.data[key] = append(ms.data[key], v)
	return nil
}

func (ms *MemoryStorage) Delete(key string, txnID uint64) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	v := &VersionedValue{
		Value:     nil,
		TxnID:     txnID,
		Committed: false,
		Timestamp: timeNowUnixNano(),
		Deleted:   true,
	}
	ms.data[key] = append(ms.data[key], v)
	return nil
}

func (ms *MemoryStorage) Scan(start, end string, txnID uint64) (map[string][]byte, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string][]byte)
	for k, versions := range ms.data {
		if (start == "" || k >= start) && (end == "" || k < end) {
			for i := len(versions) - 1; i >= 0; i-- {
				v := versions[i]
				if v.Committed && !v.Deleted {
					result[k] = v.Value
					break
				}
			}
		}
	}
	return result, nil
}

// ...后续实现方法将在后续步骤补充...
