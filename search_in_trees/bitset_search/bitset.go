package bitset_search

type bitset []uint64

func (b *bitset) Set(n int) {
	i, j := b.split(n)
	for i >= len(*b) {
		*b = append(*b, 0)
	}

	(*b)[i] = (*b)[i] | (1 << j)
}

func (b bitset) IsSet(n int) bool {
	i, j := b.split(n)
	if i >= len(b) {
		return false
	}

	return b[i]&(1<<j) != 0
}

func (b bitset) Or(with bitset) bitset {
	bits1 := b
	bits2 := with

	if len(bits1) > len(bits2) {
		bits1, bits2 = bits2, bits1
	}

	or := make(bitset, len(bits2))
	for i := 0; i < len(bits1); i++ {
		or[i] = bits1[i] | bits2[i]
	}
	for i := len(bits1); i < len(bits2); i++ {
		or[i] = bits2[i]
	}

	return or
}

func (b *bitset) split(n int) (int, int) {
	return n >> 6, n & 0x3F
}
