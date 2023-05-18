package bitmaps

type SparseBitmap64p3 struct {
	indices Bitmap64
	nodes   []SparseBitmap
}

func (b *SparseBitmap64p3) Set(n int) {
	blockPosition, nodePosition := split12bits(n)
	if blockPosition >= 64 {
		panic("out of range")
	}

	index := b.indices.Index(blockPosition)
	if !b.indices.IsSet(blockPosition) {
		if len(b.nodes) == index {
			b.nodes = append(b.nodes, SparseBitmap{})
		} else {
			b.nodes = append(b.nodes[:index+1], b.nodes[index:]...)
			b.nodes[index] = SparseBitmap{}
		}
		b.indices.Set(blockPosition)
	}
	b.nodes[index].Set(nodePosition)
}

func (b *SparseBitmap64p3) IsSet(n int) bool {
	blockPosition, nodePosition := split12bits(n)
	if blockPosition >= 64 {
		return false
	}
	if !b.indices.IsSet(blockPosition) {
		return false
	}

	index := b.indices.Index(blockPosition)

	return b.nodes[index].IsSet(nodePosition)
}

type SparseBitmap struct {
	indices Bitmap64
	bits    []Bitmap64
}

func (b *SparseBitmap) Set(n int) {
	blockPosition, bitPosition := split6bits(n)
	if blockPosition >= 64 {
		panic("out of range")
	}

	index := b.indices.Index(blockPosition)
	if !b.indices.IsSet(blockPosition) {
		if len(b.bits) == index {
			b.bits = append(b.bits, 0)
		} else {
			b.bits = append(b.bits[:index+1], b.bits[index:]...)
			b.bits[index] = 0
		}
		b.indices.Set(blockPosition)
	}
	b.bits[index] = b.bits[index] | (1 << bitPosition)
}

func (b *SparseBitmap) IsSet(n int) bool {
	blockPosition, bitPosition := split6bits(n)
	if blockPosition >= 64 {
		return false
	}
	if !b.indices.IsSet(blockPosition) {
		return false
	}

	index := b.indices.Index(blockPosition)

	return b.bits[index].IsSet(bitPosition)
}

func split6bits(n int) (int, int) {
	return n >> 6, n & 0b111111
}

func split12bits(n int) (int, int) {
	return n >> 12, n & 0b111111111111
}
