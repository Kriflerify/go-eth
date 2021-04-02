package trie

import "github.com/ethereum/go-ethereum/common"

type node interface {
	// TODO Dummy method to make node interface more specific
}

//TODO capital cases/ small cases of struct fields ???
type (
	branchNode struct {
		Children [16]node
		Val      []byte
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
