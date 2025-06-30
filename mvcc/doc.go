// Package mvcc provides a memory-based multi-version concurrency control (MVCC) library.
//
// Features:
//   - In-memory MVCC with pluggable storage interface
//   - MySQL-style isolation levels
//   - Transactional API: Begin, Commit, Rollback, Get, Set
//   - Ready for distributed and persistent extension
//
// Usage:
//
//	mvcc := mvcc.NewMVCC()
//	txn := mvcc.BeginTransaction(mvcc.RepeatableRead)
//	err := mvcc.Set(txn, "key", []byte("value"))
//	value, err := mvcc.Get(txn, "key")
//	err = mvcc.Commit(txn)
package mvcc
