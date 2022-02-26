package kvstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jacktantram/backend-go/simple-kv-store/kvstore"
)

func TestStore_Begin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(s *kvstore.KVStore)
		wantErr error
	}{
		{
			name: "should be able to begin a transaction",
		},
		{
			name: "should be able to begin nested transactions",
			setup: func(s *kvstore.KVStore) {
				s.Begin()
				s.Begin()
			},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := kvstore.NewStore()
			if tc.setup != nil {
				tc.setup(&s)
			}
			s.Begin()
		})
	}
}

func TestKVStore_Commit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(s *kvstore.KVStore)
		wantErr error
	}{
		{
			name:    "given there are no active transactions when Commit is called then ErrOperationTransactionDoesNotExist should be returned",
			wantErr: kvstore.ErrOperationTransactionDoesNotExist},
		{
			name: "should be able to commit a transaction",
			setup: func(s *kvstore.KVStore) {
				s.Begin()
			},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := kvstore.NewStore()
			if tc.setup != nil {
				tc.setup(&s)
			}
			assert.Equal(t, tc.wantErr, s.Commit())
		})
	}
}

func TestKVStore_Abort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(s *kvstore.KVStore)
		wantErr error
	}{
		{
			name:    "given there are no active transactions when Abort is called then ErrOperationTransactionDoesNotExist should be returned",
			wantErr: kvstore.ErrOperationTransactionDoesNotExist},
		{
			name: "should be able to abort a 'begun' transaction",
			setup: func(s *kvstore.KVStore) {
				s.Begin()
			},
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := kvstore.NewStore()
			if tc.setup != nil {
				tc.setup(&s)
			}
			assert.Equal(t, tc.wantErr, s.Abort())
		})
	}
}
