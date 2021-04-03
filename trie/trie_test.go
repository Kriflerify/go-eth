package trie

import (
	"bytes"
	"testing"
)

// TODO: Write tests analog to go-ethereum

func TestUpdate(t *testing.T) {
	T := newTree(true)

	tests := []struct {
		key []byte
		val []byte
	}{
		{key: []byte{'a'}, val: []byte{6}},
		{key: []byte{'b'}, val: []byte{7}},
		{key: []byte("b"), val: []byte{9}},
		{key: []byte("ac"), val: []byte{12, 13, 14}},
		{key: []byte("cat"), val: []byte("dog")},
		{key: []byte("cat"), val: []byte("doge")},
	}

	for _, test := range tests {
		T.Update(test.key, test.val)

		got, err := T.TryGet(test.key)
		if err != nil {
			t.Log(err)
		}
		if !bytes.Equal(got, test.val) {
			t.Errorf("inserted %#X  %#X, got %#X", test.key, test.val, got)
		}
	}
}
