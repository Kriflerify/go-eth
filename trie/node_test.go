package trie

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func unhex(str string) []byte {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}

func tohex(data []byte) string {
	b := hex.EncodeToString(data)
	return b
}

// Playing around with RLP Encoding:
// The RLP Encoding of an extension Node and a leaf Node may be equal
// This is why the path-key includes information about the type of node
func TestRLP(t *testing.T) {
	en := extensionNode{Key: []byte{1}, Extension: common.BytesToHash(bytes.Repeat([]byte{1}, 32))}
	ln := leafNode{Key: []byte{2}, Value: bytes.Repeat([]byte{1}, 32)}

	w1 := make(sliceBuffer, 0, 550)
	w2 := make(sliceBuffer, 0, 550)

	rlp.Encode(&w1, en)
	rlp.Encode(&w2, ln)

	t.Log(w1)
	t.Log(w2)

	r1 := bytes.NewReader(w1)
	r2 := bytes.NewReader(w2)

	enr := extensionNode{}
	lnr := leafNode{}
	rlp.Decode(r1, &enr)
	rlp.Decode(r2, &lnr)
	t.Log(enr)
	t.Log(lnr)

	r1.Reset(w1)
	r2.Reset(w2)
	enr = extensionNode{}
	lnr = leafNode{}
	rlp.Decode(r2, &enr)
	rlp.Decode(r1, &lnr)
	t.Log(enr)
	t.Log(lnr)

	t.Fail()
}

// TODO: Copy all tests from go-ethereum
