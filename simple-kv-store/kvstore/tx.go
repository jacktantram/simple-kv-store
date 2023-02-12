package kvstore

import (
	"errors"
)

var (
	ErrOperationTransactionDoesNotExist = errors.New("operation not allowed as a transaction does not exist")
)

// Transaction is an in-progress database transaction.
// A transaction must end with a call to Commit or Abort.
type Transaction struct {
	dataStore dataStore
	next      *Transaction
}

// Begin initiates a transaction
// The head of the KV store will always hold the most recent transaction.
func (s *KVStore) Begin() {
	if s.head == nil {
		s.head = &Transaction{dataStore: s.dataStore.Copy()}
		return
	}

	s.head = &Transaction{dataStore: s.head.dataStore.Copy(), next: s.head}
}

// Commit will "commit" the latest state of the head transaction to either
// the stores global state or to the next transaction along the chain.
// If there are no current transactions ErrOperationTransactionDoesNotExist will be returned.
func (s *KVStore) Commit() error {
	if s.head == nil {
		return ErrOperationTransactionDoesNotExist
	}
	if s.head.next == nil {
		s.dataStore = s.head.dataStore
		s.head = nil
		return nil
	}

	s.head.next.dataStore = s.head.dataStore.Copy()
	s.head = s.head.next
	return nil
}

// Abort will "abort" the latest state of the head transaction.
// If there is a transaction along the chain from past head it will now be the new pointer.
// If there are no current transactions ErrOperationTransactionDoesNotExist will be returned.
func (s *KVStore) Abort() error {
	if s.head == nil {
		return ErrOperationTransactionDoesNotExist
	}
	if s.head.next == nil {
		s.head = nil
		return nil
	}
	s.head = s.head.next
	return nil
}
