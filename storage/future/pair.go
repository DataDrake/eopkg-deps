//
// Copyright 2018 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package storage

import (
	"bytes"
)

// Pair associates a Key with its index in the BitMap
type Pair struct {
	Key   []byte
	Index uint64
}

// MatchKey checks if a key matches the current Pair's key
func (p Pair) MatchKey(key []byte) bool {
	return bytes.Equal(key, p.Key)
}

// MatchIndex checks if an index matches the current Pair's index
func (p Pair) MatchIndex(index uint64) bool {
	return p.Index == index
}
