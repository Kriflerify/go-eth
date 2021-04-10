package trie

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
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

	for i, test := range tests {
		T.Update(test.key, test.val)

		got, err := T.TryGet(test.key)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(got, test.val) {
			t.Errorf("%d) inserted %#X  %#X, got %#X", i, test.key, test.val, got)
		}
	}
}

func TestDelete(t *testing.T) {
	T := newTree(true)

	tests := []struct {
		key []byte
		val []byte
	}{
		{key: []byte{'a'}, val: []byte{3}},
		{key: []byte{'b'}, val: []byte{7}},
		{key: []byte("cat"), val: []byte{4, 7, 1, 1}},
		{key: []byte("catat"), val: []byte{30}},
		{key: []byte("catattack"), val: []byte{30}},
	}

	for _, test := range tests {
		err := T.Update(test.key, test.val)

		if err != nil {
			t.Error(err)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(tests), func(i, j int) {
		tests[i], tests[j] = tests[j], tests[i]
	})

	for i, test := range tests {
		err := T.Delete(test.key)

		if err != nil {
			t.Error(err)
		}

		got, err2 := T.TryGet(test.key)
		if err2 != nil {
			t.Log(err2)
		}

		if !bytes.Equal(got, []byte{}) {
			t.Errorf("%d) inserted and then deleted %#X -> %#X, but received %#X.",
				i, test.key, test.val, got)
		}
	}

}
