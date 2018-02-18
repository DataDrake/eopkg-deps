package storage

import (
	"bytes"
)

type Pair struct {
	Key   []byte
	Index uint64
}

func (p Pair) MatchKey(key []byte) bool {
	return bytes.Equal(key, p.Key)
}

func (p Pair) MatchIndex(index uint64) bool {
	return p.Index == index
}
