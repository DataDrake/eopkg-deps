package storage

import (
	"math"
)

type BitMap struct {
	RowCol []uint64
	ColRow []uint64
	Rows   uint64
	Cols   uint64
}

func NewBitMap(rows, cols uint64) *BitMap {
	size := uint64(math.Ceil(float64(rows*cols) / 64))
	return &BitMap{
		RowCol: make([]uint64, size),
		ColRow: make([]uint64, size),
		Rows:   rows,
		Cols:   cols,
	}
}

func (b *BitMap) getIndices(row, col uint64) (uint64, uint64, uint64, uint64) {
	bitsRC := (row*b.Cols + col)
	bitsCR := (col*b.Rows + row)
	return bitsRC / 64, bitsRC % 64, bitsCR / 64, bitsCR % 64
}

func (b *BitMap) IsSet(row, col uint64) bool {
	word, bit, _, _ := b.getIndices(row, col)
	return (b.RowCol[word] & (1 << bit)) > 0
}

func (b *BitMap) Set(row, col uint64) {
	wordRC, bitRC, wordCR, bitCR := b.getIndices(row, col)
	b.RowCol[wordRC] |= (1 << bitRC)
	b.ColRow[wordCR] |= (1 << bitCR)
}

func (b *BitMap) Clear(row, col uint64) {
	wordRC, bitRC, wordCR, bitCR := b.getIndices(row, col)
	b.RowCol[wordRC] &= ^(1 << bitRC)
	b.ColRow[wordCR] &= ^(1 << bitCR)
}

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
