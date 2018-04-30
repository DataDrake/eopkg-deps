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

// Graph maps relationship efficiently
type Graph struct {
	Keys []Pair
	Map  *BitMap
}

// FindIndex gets the index in the BitMap for a given Key
func (g *Graph) FindIndex(key []byte) (index uint64) {
	for _, p := range g.Keys {
		if p.MatchKey(key) {
			index = p.Index
			return
		}
	}
	return
}

// FindKey gets the key for a given index in the BitMap
func (g *Graph) FindKey(index uint64) (key []byte) {
	for _, p := range g.Keys {
		if p.MatchIndex(index) {
			key = p.Key
			return
		}
	}
	return
}

// AddKey assigns an index to a key
func (g *Graph) AddKey(key []byte) (index uint64) {
	if index = g.FindIndex(key); index == 0 {
		index = uint64(len(g.Keys)) + 1
		g.Keys = append(g.Keys, Pair{key, index})
	}
	return
}

// FinalizeMap sets up the internal BitMap based on the current Keys
func (g *Graph) FinalizeMap() {
	g.Map = NewBitMap(uint64(len(g.Keys)), uint64(len(g.Keys)))
}
