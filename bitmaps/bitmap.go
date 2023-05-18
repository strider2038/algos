package bitmaps

import "math/bits"

type Bitmap map[int]uint64

func (b *Bitmap) Set(n int) {
	blockPosition, bitPosition := split6bits(n)

	(*b)[blockPosition] = (*b)[blockPosition] | (1 << bitPosition)
}

func (b Bitmap) IsSet(n int) bool {
	blockPosition, bitPosition := split6bits(n)

	return b[blockPosition]&(1<<bitPosition) != 0
}

type Bitmap64 uint64

func (b *Bitmap64) Set(n int) {
	*b = *b | (1 << n)
}

func (b Bitmap64) IsSet(n int) bool {
	return b&(1<<n) != 0
}

func (b Bitmap64) Index(n int) int {
	return bits.OnesCount64(uint64(b) & ^(uint64(0xFFFFFFFFFFFFFFFF) << n))
}
