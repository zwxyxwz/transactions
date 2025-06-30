package mvcc

// VersionedValue 已在 storage.go 定义，这里补充版本链管理相关结构和方法。

// VersionChain represents a chain of versions for a key.
type VersionChain struct {
	Versions []*VersionedValue
}

func (vc *VersionChain) LatestVisible(txn *Transaction) *VersionedValue {
	// 简化实现：返回事务可见的最新已提交版本
	for i := len(vc.Versions) - 1; i >= 0; i-- {
		v := vc.Versions[i]
		if v.Committed && !v.Deleted && v.Timestamp <= txn.BeginTS {
			return v
		}
	}
	return nil
}
