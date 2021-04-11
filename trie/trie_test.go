package trie

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

func TestNull(t *testing.T) {
	tree := NewTree(true)
	key := make([]byte, 32)
	value := []byte("test")
	tree.Update(key, value)
	got, err := tree.TryGet(key)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(got, value) {
		t.Fatal("wrong value")
	}
}
func TestEmptyValues(t *testing.T) {
	tree := NewTree(true)

	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		updateString(tree, val.k, val.v)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vals), func(i, j int) {
		vals[i], vals[j] = vals[j], vals[i]
	})

	for _, val := range vals {
		got, err := getString(tree, val.k)
		if err != nil {
			t.Error(err)
		}

		if got != val.v {
			t.Errorf("Inserted %s -> %s, got %s", val.k, val.v, got)
		}
	}
}

func TestLargeValue(t *testing.T) {
	tree := NewTree(true)
	tree.Update([]byte("key1"), []byte{99, 99, 99, 99})
	tree.Update([]byte("key2"), bytes.Repeat([]byte{1}, 32))
}

func TestUpdateAndGet(t *testing.T) {
	T := NewTree(true)

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
	T := NewTree(true)

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
		_, err2 := T.TryGet(test.key)
		if err2 != nil {
			t.Logf("Insertion of %#X was not succesfull: %#v", test.key, err2)
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

func updateString(t *Tree, key string, value string) {
	t.Update([]byte(key), []byte(value))
}

func getString(t *Tree, key string) (string, error) {
	val, err := t.TryGet([]byte(key))
	return string(val), err
}
