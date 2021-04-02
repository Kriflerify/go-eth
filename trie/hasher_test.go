package trie

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test_sliceBuffer_Write(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		b       *sliceBuffer
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.b.Write(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("sliceBuffer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("sliceBuffer.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_sliceBuffer_Reset(t *testing.T) {
	tests := []struct {
		name string
		b    *sliceBuffer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Reset()
		})
	}
}

func Test_newHasher(t *testing.T) {
	type args struct {
		parallel bool
	}
	tests := []struct {
		name string
		args args
		want *hasher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHasher(tt.args.parallel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHasher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_returnHasherToPool(t *testing.T) {
	type args struct {
		h *hasher
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnHasherToPool(tt.args.h)
		})
	}
}

func Test_hasher_hash(t *testing.T) {
	type args struct {
		n node
	}
	tests := []struct {
		name string
		h    *hasher
		args args
		want common.Hash
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.hash(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hasher.hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasher_hashData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		h    *hasher
		args args
		want common.Hash
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.hashData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hasher.hashData() = %v, want %v", got, tt.want)
			}
		})
	}
}
