// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"reflect"
	"testing"
)

func Test_hexToCompact(t *testing.T) {
	type args struct {
		hex []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexToCompact(tt.args.hex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hexToCompact() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hexToCompactInPlace(t *testing.T) {
	type args struct {
		hex []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexToCompactInPlace(tt.args.hex); got != tt.want {
				t.Errorf("hexToCompactInPlace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compactToHex(t *testing.T) {
	type args struct {
		compact []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compactToHex(tt.args.compact); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("compactToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keybytesToHex(t *testing.T) {
	type args struct {
		str []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keybytesToHex(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keybytesToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hexToKeybytes(t *testing.T) {
	type args struct {
		hex []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexToKeybytes(tt.args.hex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("hexToKeybytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeNibbles(t *testing.T) {
	type args struct {
		nibbles []byte
		bytes   []byte
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decodeNibbles(tt.args.nibbles, tt.args.bytes)
		})
	}
}

func Test_prefixLen(t *testing.T) {
	type args struct {
		a []byte
		b []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prefixLen(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("prefixLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasTerm(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasTerm(tt.args.s); got != tt.want {
				t.Errorf("hasTerm() = %v, want %v", got, tt.want)
			}
		})
	}
}
