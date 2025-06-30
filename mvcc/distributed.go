package mvcc

import (
	"errors"
	"sync"
)

// DistributedManager 预留分布式/高可用相关接口

type DistributedManager interface {
	// 节点注册
	RegisterNode(nodeID string, addr string) error
	// 节点注销
	UnregisterNode(nodeID string) error
	// 获取所有节点
	ListNodes() ([]string, error)
	// 主节点切换
	PromoteToLeader(nodeID string) error
	// 分布式事务提交
	CommitDistributed(txnID TxnID) error
	// 分布式事务回滚
	RollbackDistributed(txnID TxnID) error
	// 分布式锁
	AcquireLock(key string, txnID TxnID) error
	ReleaseLock(key string, txnID TxnID) error
}

// DummyDistributedManager 是一个简单的内存实现，便于单元测试和后续扩展

type DummyDistributedManager struct {
	mu     sync.Mutex
	nodes  map[string]string // nodeID -> addr
	leader string
	locks  map[string]TxnID // key -> txnID
}

func NewDummyDistributedManager() *DummyDistributedManager {
	return &DummyDistributedManager{
		nodes: make(map[string]string),
		locks: make(map[string]TxnID),
	}
}

func (d *DummyDistributedManager) RegisterNode(nodeID, addr string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.nodes[nodeID] = addr
	if d.leader == "" {
		d.leader = nodeID
	}
	return nil
}

func (d *DummyDistributedManager) UnregisterNode(nodeID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.nodes, nodeID)
	if d.leader == nodeID {
		d.leader = ""
		for id := range d.nodes {
			d.leader = id
			break
		}
	}
	return nil
}

func (d *DummyDistributedManager) ListNodes() ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	ids := make([]string, 0, len(d.nodes))
	for id := range d.nodes {
		ids = append(ids, id)
	}
	return ids, nil
}

func (d *DummyDistributedManager) PromoteToLeader(nodeID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.nodes[nodeID]; !ok {
		return errors.New("node not found")
	}
	d.leader = nodeID
	return nil
}

func (d *DummyDistributedManager) CommitDistributed(txnID TxnID) error {
	// 假实现：直接返回成功
	return nil
}

func (d *DummyDistributedManager) RollbackDistributed(txnID TxnID) error {
	// 假实现：直接返回成功
	return nil
}

func (d *DummyDistributedManager) AcquireLock(key string, txnID TxnID) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if owner, ok := d.locks[key]; ok && owner != txnID {
		return errors.New("lock already held")
	}
	d.locks[key] = txnID
	return nil
}

func (d *DummyDistributedManager) ReleaseLock(key string, txnID TxnID) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if owner, ok := d.locks[key]; ok && owner == txnID {
		delete(d.locks, key)
		return nil
	}
	return errors.New("lock not held by txn")
}
