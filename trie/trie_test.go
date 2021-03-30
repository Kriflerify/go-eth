package trie

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func newEmpty() *Tree {
	root := node{}
	rootHash := common.BytesToHash(root.Encode())
	db := make(map[common.Hash]node)

	db[rootHash] = root
	T := Tree{db, rootHash}

	return &T
}

func TestUpdate(t *testing.T) {
	T := *newEmpty()

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

func TestEncoding(t *testing.T) {
	var node1, node2 node

	node1.Val = []byte{3}

	hash1 := node1.Encode()
	hash2 := node2.Encode()

	// hash1 = common.BytesToHash(hash1)
	// hash2 = common.BytesToHash(hash2)
	t.Logf("Hash1: %#X, Hash2: %#X", hash1, hash2)
}

func TestDb(t *testing.T) {
	T := *newEmpty()

	var node1, node2 node

	node1.Val = []byte{3}

	hash1 := node1.Encode()
	hash2 := node2.Encode()

	T.db[common.BytesToHash(hash1)] = node1
	T.db[common.BytesToHash(hash2)] = node2

	t.Log(T)
}
