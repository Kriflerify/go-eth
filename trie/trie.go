package trie

// import (
//  	"bytes"
//  	"errors"

// 	"github.com/ethereum/go-ethereum/common"
// 	// "github.com/ethereum/go-ethereum/rlp"
// )

// type Tree struct {
// 	db map[common.Hash]node

// 	root   common.Hash
// 	hasher *hasher
// }

// func newTree(parallel bool) *Tree {

// 	t := new(Tree)
// 	t.hasher = newHasher(parallel)
// 	t.db = make(map[common.Hash]node)

// 	t.encodeAndStore(nilValueNode)
// 	return t
// }

// // enocdeAndStore calculates the hash of a given node and adds it to the database of tree
// func (t *Tree) encodeAndStore(n node) common.Hash {
// 	hashed := t.hasher.hash(n)
// 	t.db[hashed] = n
// 	return hashed
// }

// // Update associates key with value.
// func (t *Tree) Update(key, value []byte) error {
// 	k := keybytesToHex(key)

// 	rootNode := t.db[t.root]

// 	newRoot, err := t.insert(rootNode, []byte{}, k, value)

// 	delete(t.db, t.root)

// 	newRootHash := newRoot.Encode()
// 	t.db[common.BytesToHash(newRootHash)] = newRoot
// 	t.root = common.BytesToHash(newRootHash)
// 	return err
// }

// func (t *Tree) insert(n node, prefix []byte, key []byte, value []byte) (node, error) {
// 	if len(key) == 0 {
// 		return node{}, errors.New("non-terminated String")
// 	} else if bytes.Equal(key, []byte{16}) {
// 		// Key termination symbol
// 		n.Val = value
// 		return n, nil
// 	}

// 	a := key[0]
// 	key = key[1:]
// 	prefix = append(prefix, a)

// 	var child node

// 	childHash := n.Children[a]

// 	if childHash != nil {
// 		child = t.db[common.BytesToHash(childHash)]
// 	} else {
// 		child = node{Children: [17]hashNode{}}
// 	}

// 	newChild, err := t.insert(child, prefix, key, value)
// 	if err != nil {
// 		return n, err
// 	}
// 	newChildHash := newChild.Encode()
// 	t.db[common.BytesToHash(newChildHash)] = newChild
// 	n.Children[a] = newChildHash

// 	delete(t.db, common.BytesToHash(childHash))

// 	return n, nil

// }

// // TryGet retrieves the value associated with key
// func (t *Tree) TryGet(key []byte) ([]byte, error) {
// 	k := keybytesToHex(key)
// 	root := t.db[t.root]
// 	return t.tryGet(root, k, 0)
// }

// func (t *Tree) tryGet(n node, key []byte, pos int) (value []byte, err error) {
// 	if len(key) == 0 {
// 		return []byte{}, nil
// 	} else if bytes.Equal(key, []byte{16}) {
// 		return n.Val, nil
// 	}

// 	a := key[0]
// 	childHash := n.Children[a]
// 	if childHash == nil {
// 		return nil, errors.New("key not found")
// 	}

// 	key = key[1:]

// 	child := t.db[common.BytesToHash(childHash)]
// 	return t.tryGet(child, key, pos+1)
// }
