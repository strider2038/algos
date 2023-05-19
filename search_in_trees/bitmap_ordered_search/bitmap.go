package bitmap_ordered_search

type bitmap map[int]uint64

func (b *bitmap) Set(n int) {
	i, j := b.split(n)

	(*b)[i] = (*b)[i] | (1 << j)
}

func (b bitmap) IsSet(n int) bool {
	i, j := b.split(n)

	return b[i]&(1<<j) != 0
}

func (b bitmap) Or(with bitmap) bitmap {
	bits1 := b
	bits2 := with

	if len(bits1) > len(bits2) {
		bits1, bits2 = bits2, bits1
	}

	or := make(bitmap, len(bits2))
	for i := range bits1 {
		bits := bits1[i] | bits2[i]
		if bits != 0 {
			or[i] = bits
		}
	}
	for i := range bits2 {
		bits := bits1[i] | bits2[i]
		if bits != 0 {
			or[i] = bits
		}
	}

	return or
}

func (b *bitmap) split(n int) (int, int) {
	return n >> 6, n & 0x3F
}
