package mvcc

import (
	"sync"
	"testing"
)

func TestMVCC_Basic(t *testing.T) {
	mvcc := NewMVCC()
	// 事务1写入
	txn1 := mvcc.BeginTransaction(RepeatableRead)
	err := mvcc.Set(txn1, "foo", []byte("bar"))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	// 事务1提交
	err = mvcc.Commit(txn1)
	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	// 事务2读取
	txn2 := mvcc.BeginTransaction(RepeatableRead)
	val, err := mvcc.Get(txn2, "foo")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(val) != "bar" {
		t.Fatalf("Expected 'bar', got '%s'", string(val))
	}
}

func TestMVCC_Rollback(t *testing.T) {
	mvcc := NewMVCC()
	txn := mvcc.BeginTransaction(RepeatableRead)
	err := mvcc.Set(txn, "k1", []byte("v1"))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	err = mvcc.Rollback(txn)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}
	txn2 := mvcc.BeginTransaction(RepeatableRead)
	_, err = mvcc.Get(txn2, "k1")
	if err == nil {
		t.Fatalf("Expected error for rolled back key, got nil")
	}
}

func TestMVCC_ConcurrentWrite(t *testing.T) {
	mvcc := NewMVCC()
	wg := sync.WaitGroup{}
	write := func(key, val string) {
		txn := mvcc.BeginTransaction(RepeatableRead)
		_ = mvcc.Set(txn, key, []byte(val))
		_ = mvcc.Commit(txn)
		wg.Done()
	}
	wg.Add(2)
	go write("c1", "v1")
	go write("c1", "v2")
	wg.Wait()
	txn := mvcc.BeginTransaction(RepeatableRead)
	v, err := mvcc.Get(txn, "c1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(v) != "v1" && string(v) != "v2" {
		t.Fatalf("Expected 'v1' or 'v2', got '%s'", string(v))
	}
}

func TestMVCC_Delete(t *testing.T) {
	mvcc := NewMVCC()
	txn := mvcc.BeginTransaction(RepeatableRead)
	_ = mvcc.Set(txn, "d1", []byte("toDel"))
	_ = mvcc.Commit(txn)
	txn2 := mvcc.BeginTransaction(RepeatableRead)
	_ = mvcc.Delete(txn2, "d1")
	_ = mvcc.Commit(txn2)
	txn3 := mvcc.BeginTransaction(RepeatableRead)
	_, err := mvcc.Get(txn3, "d1")
	if err == nil {
		t.Fatalf("Expected error for deleted key, got nil")
	}
}

func TestDistributedManagerDummy(t *testing.T) {
	dm := &DummyDistributedManager{}
	if err := dm.RegisterNode("n1", "127.0.0.1"); err != nil {
		t.Fatalf("RegisterNode failed: %v", err)
	}
	if err := dm.UnregisterNode("n1"); err != nil {
		t.Fatalf("UnregisterNode failed: %v", err)
	}
	nodes, err := dm.ListNodes()
	if err != nil || len(nodes) == 0 {
		t.Fatalf("ListNodes failed: %v", err)
	}
	if err := dm.PromoteToLeader("n1"); err != nil {
		t.Fatalf("PromoteToLeader failed: %v", err)
	}
	if err := dm.CommitDistributed(1); err != nil {
		t.Fatalf("CommitDistributed failed: %v", err)
	}
	if err := dm.RollbackDistributed(1); err != nil {
		t.Fatalf("RollbackDistributed failed: %v", err)
	}
	if err := dm.AcquireLock("k", 1); err != nil {
		t.Fatalf("AcquireLock failed: %v", err)
	}
	if err := dm.ReleaseLock("k", 1); err != nil {
		t.Fatalf("ReleaseLock failed: %v", err)
	}
}
