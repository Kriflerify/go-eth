package trie

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type sliceBuffer []byte

func (b *sliceBuffer) Write(data []byte) (n int, err error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *sliceBuffer) Reset() {
	*b = (*b)[:0]
}

// hasher is a type used for the trie Hash operation. A hasher has some
// internal preallocated temp space
type hasher struct {
	sha      crypto.KeccakState
	tmp      sliceBuffer
	parallel bool // Whether to use paralallel threads when hashing
}

// hasherPool holds pureHashers
var hasherPool = sync.Pool{
	New: func() interface{} {
		return &hasher{
			tmp: make(sliceBuffer, 0, 550), // cap is as large as a full fullNode.
			sha: sha3.NewLegacyKeccak256().(crypto.KeccakState),
		}
	},
}

func newHasher(parallel bool) *hasher {
	h := hasherPool.Get().(*hasher)
	h.parallel = parallel
	return h
}

func returnHasherToPool(h *hasher) {
	hasherPool.Put(h)
}

// hash collapses a node down into a hash node, also returning a copy of the
// original node initialized with the computed hash to replace the original one.
func (h *hasher) hash(n node) common.Hash {
	rlp.Encode(&h.tmp, n)

	hashed := common.BytesToHash(h.hashData(h.tmp))
	return hashed
}

// hashData hashes the provided data
func (h *hasher) hashData(data []byte) []byte {
	n := make([]byte, 32)
	h.sha.Reset()
	h.sha.Write(data)
	h.sha.Read(n)
	return n
}

// // proofHash is used to construct trie proofs, and returns the 'collapsed'
// // node (for later RLP encoding) aswell as the hashed node -- unless the
// // node is smaller than 32 bytes, in which case it will be returned as is.
// // This method does not do anything on value- or hash-nodes.
// func (h *hasher) proofHash(original node) (collapsed, hashed node) {
// 	switch n := original.(type) {
// 	case *shortNode:
// 		sn, _ := h.hashShortNodeChildren(n)
// 		return sn, h.shortnodeToHash(sn, false)
// 	case *fullNode:
// 		fn, _ := h.hashFullNodeChildren(n)
// 		return fn, h.fullnodeToHash(fn, false)
// 	default:
// 		// Value and hash nodes don't have children so they're left as were
// 		return n, n
// 	}
// }
