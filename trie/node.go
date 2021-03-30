package trie

type hashNode []byte

type node struct {
	Children [17]hashNode
	Val      []byte
}

// turns a node into a hash
func (n *node) Encode() hashNode {
	hasher := newHasher(false)
	hashed := hasher.hash(*n)
	return hashed
}
