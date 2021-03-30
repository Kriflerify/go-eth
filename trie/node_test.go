package trie

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type vnode []byte

func TestNodeHashing(t *testing.T) {
	h := hasherPool.Get().(*hasher)
	h.parallel = true
	h.tmp.Reset()

	// b := new(bytes.Buffer)

	tests := []node{
		{},
		{Val: []byte{6}},
		{Children: [17]hashNode{hashNode([]byte{1})}, Val: []byte{6}},
	}

	for _, test := range tests {
		reflection := reflect.ValueOf(test)
		reflectionValueType := reflection.Type()
		rlpEncoding, _ := rlp.EncodeToBytes(test)
		hash := common.BytesToHash(test.Encode())
		_ = hash
		_ = rlpEncoding
		// t.Log(hash)
		// t.Log(rlpEncoding)
		t.Log(reflection)
		t.Log(reflectionValueType)
	}
}

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

func Test_node_Encode(t *testing.T) {
	tests := []struct {
		name string
		n    *node
		want hashNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("node.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
