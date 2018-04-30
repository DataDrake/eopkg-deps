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
	"testing"
)

var bMap *BitMap

func init() {
	bMap = NewBitMap(MapTestSize, MapTestSize)
	for i := uint64(0); i < MapTestSize; i++ {
		bMap.Set(i, i)
	}
}

func TestBitSet(t *testing.T) {
	b := NewBitMap(4, 4)
	println("Before")
	b.Print()
	for i := uint64(0); i < 4; i++ {
		for j := uint64(0); j < 4; j++ {
			if (i * j) > 4 {
				b.Set(i, j)
				if !b.IsSet(i, j) {
					t.Errorf("Should have set <%d,%d>", i, j)
				}
			}
		}
	}
	for i := uint64(0); i < 4; i++ {
		for j := uint64(0); j < 4; j++ {
			if (i*j) > 4 && !b.IsSet(i, j) {
				t.Errorf("Should have set <%d,%d>", i, j)
			}
		}
	}
	println("After")
	b.Print()
}

func BenchmarkBitSet(t *testing.B) {
	n := uint64(t.N)
	b := NewBitMap(10000, 10000)
	for i := uint64(0); i < n; i++ {
		b.Set(i%10000, i%10000)
	}
}

func BenchmarkBitClear(t *testing.B) {
	n := uint64(t.N)
	b := NewBitMap(10000, 10000)
	for i := uint64(0); i < 10000; i++ {
		b.Set(i, i)
	}
	t.ResetTimer()
	for i := uint64(0); i < n; i++ {
		b.Clear(i%10000, i%10000)
	}
}

func BenchmarkBitGetRow(t *testing.B) {
	n := uint64(t.N)
	for i := uint64(0); i < n; i++ {
		bMap.GetRow(i % MapTestSize)
	}
}

func BenchmarkBitGetCol(t *testing.B) {
	n := uint64(t.N)
	for i := uint64(0); i < n; i++ {
		bMap.GetCol(i % MapTestSize)
	}
}
