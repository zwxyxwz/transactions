package mvcc

// TxnID is the type for transaction IDs.
type TxnID uint64

// Version is the type for version numbers.
type Version uint64

// IsolationLevel defines supported isolation levels.
type IsolationLevel int

const (
	ReadUncommitted IsolationLevel = iota
	ReadCommitted
	RepeatableRead
	Serializable
)
