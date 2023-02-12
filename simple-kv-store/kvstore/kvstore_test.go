package kvstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jacktantram/backend-go/simple-kv-store/kvstore"
)

func TestStore_Set(t *testing.T) {
	t.Parallel()
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func(s *kvstore.KVStore)
	}{
		{
			name: "given a create request with an empty key it should be rejected",
			args: args{
				key:   "",
				value: "some-val",
			},
			wantErr: true,
		},
		{
			name: "given a create request with a key that is purely whitespace it should be rejected",
			args: args{
				key:   "    ",
				value: "some-val",
			},
			wantErr: true,
		},
		{
			name: "given a create request with a key,value and the value is blank it should be rejected",
			args: args{
				key:   "valid-key",
				value: "",
			},
			wantErr: true,
		},
		{
			name: "given a create request with a key,value and the value is only whitespace",
			args: args{
				key:   "valid-key",
				value: "  ",
			},
			wantErr: true,
		},
		{
			name: "given a set request with a valid key value then it should be accepted and stored",
			args: args{
				key:   "valid-key",
				value: "valid-value",
			},
		},
		{
			name: "given a set request with multiple layers of transactions then value should change at each transaction but not affect the others if aborted",
			args: args{
				key:   "valid-key",
				value: "valid-value",
			},
			setup: func(s *kvstore.KVStore) {
				s.Begin()
				require.NoError(t, s.Set("a-key", "another-value"))
				s.Begin()
				require.NoError(t, s.Set("a-key", "another-value-2"))
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			s := kvstore.NewStore()
			if tc.setup != nil {
				tc.setup(&s)
			}
			err := s.Set(tc.args.key, tc.args.value)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStore_Get(t *testing.T) {
	t.Parallel()
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		setup  func(s *kvstore.KVStore)
		args   args
		expOk  bool
		expVal string
	}{
		{
			name: "given a get request when the key does not exist should return false",
			args: args{
				key: "a-random-key",
			},
			expOk:  false,
			expVal: "",
		},
		{
			name: "given a get request when the key does exists it should return the correct value",
			setup: func(s *kvstore.KVStore) {
				require.NoError(t, s.Set("valid-key", "valid-value"))
			},
			args: args{
				key: "valid-key",
			},
			expOk:  true,
			expVal: "valid-value",
		},
		{
			name: "given a get request when the value is set in a transaction then it should fetch from transaction store.",
			setup: func(s *kvstore.KVStore) {
				_, ok := s.Get("a-key")
				require.False(t, ok)
				s.Begin()
				require.NoError(t, s.Set("a-key", "a-value"))
			},
			args: args{
				key: "a-key",
			},
			expOk:  true,
			expVal: "a-value",
		},
		{
			name: "given a get request when the value is set in a 2 nested transaction then it should fetch from most recent transaction store.",
			setup: func(s *kvstore.KVStore) {
				_, ok := s.Get("a-key")
				require.False(t, ok)
				s.Begin()
				require.NoError(t, s.Set("a-key", "a-value"))
				s.Begin()
				require.NoError(t, s.Set("a-key", "the-value-to-be"))
			},
			args: args{
				key: "a-key",
			},
			expOk:  true,
			expVal: "the-value-to-be",
		},
		{
			name: "given a get request when the value is set in a 3 nested transaction then it should fetch from most recent transaction store.",
			setup: func(s *kvstore.KVStore) {
				_, ok := s.Get("a-key")
				require.False(t, ok)
				s.Begin()
				require.NoError(t, s.Set("a-key", "a-value"))
				s.Begin()
				require.NoError(t, s.Set("a-key", "the-value-to-be"))
				s.Begin()
				require.NoError(t, s.Set("a-key", "the-value-to-be-3"))

			},
			args: args{
				key: "a-key",
			},
			expOk:  true,
			expVal: "the-value-to-be-3",
		},
		{
			name: "given a get request when the value is set in a 2 nested transaction and 1 is aborted then it should fetch the value from the first transaction store.",
			setup: func(s *kvstore.KVStore) {
				_, ok := s.Get("a-key")
				require.False(t, ok)
				s.Begin()
				require.NoError(t, s.Set("a-key", "a-value"))
				s.Begin()
				require.NoError(t, s.Set("a-key", "the-value-to-be"))
				require.NoError(t, s.Abort())
			},
			args: args{
				key: "a-key",
			},
			expOk:  true,
			expVal: "a-value",
		},
		{
			name: "given a get request when the value is nested transaction and both are aborted then it should fetch the value from the root store.",
			setup: func(s *kvstore.KVStore) {
				_, ok := s.Get("a-key")
				require.False(t, ok)
				s.Begin()
				require.NoError(t, s.Set("a-key", "a-value"))
				s.Begin()
				require.NoError(t, s.Set("a-key", "the-value-to-be"))
				require.NoError(t, s.Abort())
				require.NoError(t, s.Abort())
			},
			args: args{
				key: "a-key",
			},
			expOk:  false,
			expVal: "",
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
			val, ok := s.Get(tc.args.key)
			assert.Equal(t, tc.expOk, ok)
			assert.Equal(t, tc.expVal, val)
		})
	}
}

func TestStore_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		setup func(s *kvstore.KVStore)
		args  args
		expOk bool
	}{
		{
			name: "given a delete request where the key does not exist should return false",
			args: args{
				key: "a-random-key",
			},
			expOk: false,
		},
		{
			name: "given a delete request when the key does exists then it should be deleted and True is returned",
			setup: func(s *kvstore.KVStore) {
				require.NoError(t, s.Set("valid-key", "valid-value"))
			},
			args: args{
				key: "valid-key",
			},
			expOk: true,
		},
		{
			name: "given a delete request when the value is set in a transaction then it should be deleted from transaction store.",
			setup: func(s *kvstore.KVStore) {
				require.NoError(t, s.Set("a-key", "a-value"))
				s.Begin()
			},
			args: args{
				key: "a-key",
			},
			expOk: true,
		},
		{
			name: "given a delete request when the value is set in a transaction then it should not affect the root store",
			setup: func(s *kvstore.KVStore) {
				require.NoError(t, s.Set("a-key", "a-value"))
				s.Begin()
				assert.True(t, s.Delete("a-key"))
				require.NoError(t, s.Abort())
			},
			args: args{
				key: "a-key",
			},
			expOk: true,
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
			ok := s.Delete(tc.args.key)
			assert.Equal(t, tc.expOk, ok)
		})
	}
}
