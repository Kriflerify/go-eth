package trie

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/rlp"
)

type Tree struct {
	db map[common.Hash]node

	root   common.Hash
	hasher *hasher

	// The hash value of the nil node; does not depend on tree,
	// but only on RLP encoding and hashing algortim of nodes
	nilValueNodeHash common.Hash
}

func newTree(parallel bool) *Tree {

	t := new(Tree)
	t.hasher = newHasher(parallel)
	t.db = make(map[common.Hash]node)

	t.nilValueNodeHash = t.encodeAndStore(nilValueNode)

	return t
}

// enocdeAndStore calculates the hash of a given node and adds it to the database of tree
func (t *Tree) encodeAndStore(n node) common.Hash {
	hashed := t.hasher.hash(n)
	t.db[hashed] = n
	return hashed
}

// Update associates key with value.
func (t *Tree) Update(key, value []byte) error {
	k := keybytesToHex(key)

	newRoot, err := t.insert(t.root, []byte{}, k, value)
	if err != nil {
		return err
	}
	t.root = newRoot
	delete(t.db, t.root)
	return nil
}

func (t *Tree) insert(h common.Hash, prefix []byte, key []byte, value []byte) (common.Hash, error) {

	if _, ok := t.db[h]; (h == common.Hash{}) || h == t.nilValueNodeHash || !ok {
		n := leafNode{hexToCompact(key), value}
		newHash := t.encodeAndStore(n)
		return newHash, nil
	}

	a := key[0]

	switch n := t.db[h].(type) {
	case branchNode:
		if bytes.Equal(key, []byte{16}) {
			n.Value = value
			newHash := t.encodeAndStore(n)
			delete(t.db, h)
			return newHash, nil
		}

		newChild, err := t.insert(n.Children[a], append(prefix, a), key[1:], value)
		n.Children[a] = newChild
		newHash := t.encodeAndStore(n)
		delete(t.db, h)
		return newHash, err
	case extensionNode:
		nKey := compactToHex(n.Key)
		prefixLen := prefixLen(key, nKey)

		if prefixLen == len(nKey) {
			// key must be nKey + suffix; n.Extension must point to a branchNode
			newChild, err := t.insert(n.Extension, append(prefix, a), key[prefixLen:], value)
			n.Extension = newChild
			newHash := t.encodeAndStore(n)
			delete(t.db, h)
			return newHash, err
		}

		// Replace the extensionNode n with:
		// (1) extension Node -> (2) branchNode -> {(3) extensionNode, (4) leafNode}
		n.Key = hexToCompact(nKey[prefixLen:])
		n3 := extensionNode{hexToCompact(nKey[prefixLen:]), n.Extension}
		n3hash := t.encodeAndStore(n3)

		n2Children := [16]common.Hash{}
		n2Children[nKey[prefixLen]] = common.BytesToHash((n3hash[:]))
		n2 := branchNode{n2Children, []byte{}}
		n2hash := t.encodeAndStore(n2) // TODO: avoid unnecessary save and delete from db
		n2hash, err := t.insert(n2hash, append(prefix, key[:prefixLen]...), key[prefixLen:], value)
		if err != nil {
			return n2hash, err
		}

		// we konw that prefixLen != 0
		n1 := extensionNode{hexToCompact(key[:prefixLen]), n2hash}
		n1hash := t.encodeAndStore(n1)
		delete(t.db, h)
		return n1hash, nil
	case leafNode:
		nKey := compactToHex(n.Key)
		prefixLen := prefixLen(key, nKey)

		if prefixLen == len(key) && prefixLen == len(nKey) {
			n.Value = value
			newHash := t.encodeAndStore(n)
			delete(t.db, h)
			return newHash, nil
		}

		// two leafnodes with different are replaced by:
		// (1) extension Node -> (2) branchNode -> {(3) Leaf node, (4) Leaf node}
		// (4) is evtl. omitted by set (2).Val
		n2 := branchNode{}
		if len(nKey) == prefixLen+1 {
			n2.Value = n.Value
			n3 := leafNode{hexToCompact(key[prefixLen:]), value}
			n3hash := t.encodeAndStore(n3)
			n2.Children[key[prefixLen]] = n3hash
		} else if len(key) == prefixLen+1 {
			n2.Value = value
			n4 := leafNode{hexToCompact(nKey[prefixLen:]), n.Value}
			n4hash := t.encodeAndStore(n4)
			n2.Children[nKey[prefixLen]] = n4hash
		} else {
			n3 := leafNode{hexToCompact(key[prefixLen:]), value}
			n3hash := t.encodeAndStore(n3)
			n2.Children[key[prefixLen]] = n3hash
			n4 := leafNode{hexToCompact(nKey[prefixLen:]), n.Value}
			n4hash := t.encodeAndStore(n4)
			n2.Children[nKey[prefixLen]] = n4hash
		}

		n2hash := t.encodeAndStore(n2)
		n1 := extensionNode{Key: hexToCompact(key[:prefixLen]),
			Extension: n2hash,
		}
		n1hash := t.encodeAndStore(n1)
		delete(t.db, h)
		return n1hash, nil
	}
	return h, errors.New("unexpected node type")
}

// TryGet retrieves the value associated with key
func (t *Tree) TryGet(key []byte) ([]byte, error) {
	k := keybytesToHex(key)
	return t.tryGet(t.root, k, 0)
}

func (t *Tree) tryGet(h common.Hash, key []byte, pos int) (value []byte, err error) {
	if (h == common.Hash{}) || (h == t.nilValueNodeHash) {
		return nil, errors.New("key not found")
	}
	n, inDb := t.db[h]
	if !inDb {
		return nil, errors.New("key not found")
	}

	// key shold have at least a Termination symbol at the end
	a := key[0]
	switch n := n.(type) {
	case branchNode:
		if bytes.Equal(key, []byte{16}) {
			return n.Value, nil
		}
		aNode := n.Children[a]
		if (aNode == common.Hash{}) || (aNode == t.nilValueNodeHash) {
			return nil, errors.New("key not found")
		}
		return t.tryGet(aNode, key[1:], pos+1)
	case leafNode:
		if bytes.Equal(key, []byte{16}) {
			return n.Value, nil
		}
		return nil, errors.New("key not found")
	case extensionNode:
		if prefixLen := prefixLen(key, n.Key); prefixLen == len(n.Key) {
			return t.tryGet(n.Extension, key[prefixLen:], pos+prefixLen)
		}
	}
	return nil, errors.New("unexpected node")
}
