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

	t.nilValueNodeHash = t.hasher.hash(nilValueNode)

	t.root = t.nilValueNodeHash

	return t
}

// enocdeAndStore calculates the hash of a given node and adds it to the database of tree
func (t *Tree) encodeAndStore(n node) common.Hash {
	if n != nil {
		hashed := t.hasher.hash(n)
		t.db[hashed] = n
		return hashed
	}
	return t.nilValueNodeHash
}

// Update associates key with value.
func (t *Tree) Update(key []byte, value []byte) error {
	if bytes.Equal(value, []byte{}) {
		return t.Delete(key)
	}

	k := keybytesToHex(key)
	newRoot, err := t.insert(t.root, []byte{}, k, value)
	if err != nil {
		return err
	}
	delete(t.db, t.root)
	t.root = newRoot
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
			n3 := leafNode{hexToCompact(key[prefixLen+1:]), value}
			n3hash := t.encodeAndStore(n3)
			n2.Children[key[prefixLen]] = n3hash
		} else if len(key) == prefixLen+1 {
			n2.Value = value
			n4 := leafNode{hexToCompact(nKey[prefixLen+1:]), n.Value}
			n4hash := t.encodeAndStore(n4)
			n2.Children[nKey[prefixLen]] = n4hash
		} else {
			n3 := leafNode{hexToCompact(key[prefixLen+1:]), value}
			n3hash := t.encodeAndStore(n3)
			n2.Children[key[prefixLen]] = n3hash
			n4 := leafNode{hexToCompact(nKey[prefixLen+1:]), n.Value}
			n4hash := t.encodeAndStore(n4)
			n2.Children[nKey[prefixLen]] = n4hash
		}

		n2hash := t.encodeAndStore(n2)

		if prefixLen == 0 {
			delete(t.db, h)
			return n2hash, nil
		} else {
			n1 := extensionNode{Key: hexToCompact(key[:prefixLen]),
				Extension: n2hash,
			}
			n1hash := t.encodeAndStore(n1)
			delete(t.db, h)
			return n1hash, nil
		}
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
		if bytes.Equal(key, compactToHex(n.Key)) {
			return n.Value, nil
		}
		return nil, errors.New("key not found")
	case extensionNode:
		nKey := compactToHex(n.Key)
		if prefixLen := prefixLen(key, nKey); prefixLen == len(nKey) {
			return t.tryGet(n.Extension, key[prefixLen:], pos+prefixLen)
		}
		return nil, errors.New("key not found")
	}
	return nil, errors.New("unexpected node")
}

//Delete removes any mapping from key
func (t *Tree) Delete(key []byte) error {
	k := keybytesToHex(key)

	newRoot, err := t.delete(t.root, k)
	delete(t.db, t.root)
	t.root = t.encodeAndStore(newRoot)
	return err
}

func (t *Tree) delete(h common.Hash, key []byte) (node, error) {
	n, ok := t.db[h]
	if (h == common.Hash{}) || (h == t.nilValueNodeHash) || !ok {
		return nil, errors.New("trying to delete a nonexisting key")
	}

	switch n := n.(type) {
	case branchNode:
		if bytes.Equal(key, []byte{16}) {
			n.Value = nil
			if onlyChild, yes := t.hasOnlyChild(n); yes {
				child := t.db[onlyChild]
				delete(t.db, onlyChild)
				delete(t.db, h)
				return child, nil
			}
			delete(t.db, h)
			return n, nil
		}
		child := n.Children[key[0]]
		newChild, err := t.delete(child, key[1:])
		if err != nil {
			return nil, err
		}

		if newChild == nil {
			n.Children[key[0]] = common.Hash{}
			if t.isEmptyBranchNode(n) {
				delete(t.db, h)
				return nil, nil
			}

			if onlyChild, yes := t.hasOnlyChild(n); yes {
				child := t.db[onlyChild]
				delete(t.db, onlyChild)
				delete(t.db, h)
				return child, nil
			}
		}

		n.Children[key[0]] = t.encodeAndStore(newChild)
		delete(t.db, h)
		return n, nil
	case leafNode:
		if nKey := compactToHex(n.Key); bytes.Equal(nKey, key) {
			return nil, nil
		}
		return n, errors.New("trying to delete a nonexisting key")
	case extensionNode:
		nKey := compactToHex(n.Key)
		prefixLen := prefixLen(nKey, key)

		if prefixLen < len(nKey) {
			return nil, errors.New("trying to delete a nonexisting key")
		}

		//from now on assume that extension used to be a branchNode
		newChild, err := t.delete(n.Extension, key[prefixLen:])
		if err != nil {
			return nil, err
		}
		if newChild == nil {
			// the branch node must not be needed anymore
			delete(t.db, h)
			return nil, nil
		}

		switch newChild := newChild.(type) {
		case extensionNode:
			combinedKey := append(nKey, compactToHex(newChild.Key)...)
			newChild.Key = hexToCompact(combinedKey)
		case leafNode:
			combinedKey := append(nKey, compactToHex(newChild.Key)...)
			newChild.Key = hexToCompact(combinedKey)
		}

		delete(t.db, h)
		return newChild, nil
	}
	return common.Hash{}, errors.New("tree error: unrecognizable node type")
}

func (t *Tree) isEmptyBranchNode(n branchNode) bool {
	if !bytes.Equal(n.Value, []byte{}) {
		return false
	}
	for _, c := range n.Children {
		if (c != common.Hash{}) && (c != t.nilValueNodeHash) {
			return false
		}
	}
	return true
}

func (t *Tree) hasOnlyChild(n branchNode) (common.Hash, bool) {
	i := -1
	if bytes.Equal(n.Value, []byte{}) {
		return common.Hash{}, false
	}
	for j, c := range n.Children {
		if (c != common.Hash{}) && (c != t.nilValueNodeHash) {
			if i == -1 {
				i = j
			} else {
				return common.Hash{}, false
			}
		}
	}
	return n.Children[i], true
}
