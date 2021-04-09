package trie

import "github.com/ethereum/go-ethereum/common"

type node interface {
	// optional: Dummy method to make node interface more specific
}

type (
	branchNode struct {
		Children [16]common.Hash
		Value    []byte
	}
	extensionNode struct {
		Key       []byte
		Extension common.Hash
	}
	leafNode struct {
		Key   []byte
		Value []byte
	}
)

var nilValueNode = []byte(nil)
