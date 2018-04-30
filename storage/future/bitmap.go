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
	"math"
)

// BitMap efficiently stores relationships by only requiring 2-bits per relationship
type BitMap struct {
	RowCol []uint64
	ColRow []uint64
	Rows   uint64
	Cols   uint64
}

// NewBitMap returns an empty BitMap of the specified size
func NewBitMap(rows, cols uint64) *BitMap {
	size := uint64(math.Ceil(float64(rows*cols) / 64))
	return &BitMap{
		RowCol: make([]uint64, size),
		ColRow: make([]uint64, size),
		Rows:   rows,
		Cols:   cols,
	}
}

// getIndices performes index calculations
func (b *BitMap) getIndices(row, col uint64) (uint64, uint64, uint64, uint64) {
	bitsRC := (row*b.Cols + col)
	bitsCR := (col*b.Rows + row)
	return bitsRC / 64, bitsRC % 64, bitsCR / 64, bitsCR % 64
}

// IsSet checks if a relationship exists
func (b *BitMap) IsSet(row, col uint64) bool {
	word, bit, _, _ := b.getIndices(row, col)
	return (b.RowCol[word] & (1 << bit)) > 0
}

// Set create a relationship
func (b *BitMap) Set(row, col uint64) {
	wordRC, bitRC, wordCR, bitCR := b.getIndices(row, col)
	b.RowCol[wordRC] |= (1 << bitRC)
	b.ColRow[wordCR] |= (1 << bitCR)
}

// Clear breaks a relationship
func (b *BitMap) Clear(row, col uint64) {
	wordRC, bitRC, wordCR, bitCR := b.getIndices(row, col)
	b.RowCol[wordRC] &= ^(1 << bitRC)
	b.ColRow[wordCR] &= ^(1 << bitCR)
}

// GetRow reads all of the relationships for a given Row
func (b *BitMap) GetRow(index uint64) []uint64 {
	bRow := make([]uint64, uint64(math.Ceil(float64(b.Cols)/64)))
	col := uint64(0)
	srcWord, bit, _, _ := b.getIndices(index, 0)
	dstWord := 0
	for col+64 < b.Cols {
		bRow[dstWord] = (b.RowCol[srcWord] << (64 - bit)) | (b.RowCol[srcWord+1] >> bit)
		col += 64
		dstWord++
		srcWord++
	}
	if col < b.Cols {
		bRow[dstWord] = b.RowCol[srcWord] << (64 - bit)
	}
	return bRow
}

// GetCol reads all of the relationships for a given Column
func (b *BitMap) GetCol(index uint64) []uint64 {
	bRow := make([]uint64, uint64(math.Ceil(float64(b.Rows)/64)))
	row := uint64(0)
	_, _, srcWord, bit := b.getIndices(index, 0)
	dstWord := 0
	for row+64 < b.Rows {
		bRow[dstWord] = (b.ColRow[srcWord] << (64 - bit)) | (b.ColRow[srcWord+1] >> bit)
		row += 64
		dstWord++
		srcWord++
	}
	if row < b.Rows {
		bRow[dstWord] = b.ColRow[srcWord] << (64 - bit)
	}
	return bRow
}

// Print prints out a BitMap for manual verification
func (b *BitMap) Print() {
	for row := uint64(0); row < b.Rows; row++ {
		for col := uint64(0); col < b.Cols; col++ {
			if b.IsSet(row, col) {
				print("1,")
			} else {
				print("0,")
			}
		}
		print("\n")
	}
}
