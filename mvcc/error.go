package mvcc

import "errors"

var (
	ErrKeyNotFound      = errors.New("key not found")
	ErrTxnNotFound      = errors.New("transaction not found")
	ErrTxnNotActive     = errors.New("transaction not active")
	ErrWriteConflict    = errors.New("write conflict")
	ErrIsolationInvalid = errors.New("invalid isolation level")
)
