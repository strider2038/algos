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

func (b SparseBitmap64p3) Or(with SparseBitmap64p3) SparseBitmap64p3 {
	or := SparseBitmap64p3{}
	indices := b.indices | with.indices
	for i := 0; i < 64 && indices != 0; i++ {
		if indices&1 != 0 {
			node := b.nodeAt(i).Or(with.nodeAt(i))
			if !node.IsZero() {
				or.setNodeAt(i, node)
			}
		}
		indices >>= 1
	}

	return or
}

func (b *SparseBitmap64p3) nodeAt(blockPosition int) SparseBitmap {
	if !b.indices.IsSet(blockPosition) {
		return SparseBitmap{}
	}

	index := b.indices.Index(blockPosition)

	return b.nodes[index]
}

func (b *SparseBitmap64p3) setNodeAt(blockPosition int, node SparseBitmap) {
	index := b.indices.Index(blockPosition)
	if !b.indices.IsSet(blockPosition) {
		if len(b.nodes) == index {
			b.nodes = append(b.nodes, SparseBitmap{})
		} else {
			b.nodes = append(b.nodes[:index+1], b.nodes[index:]...)
		}
		b.indices.Set(blockPosition)
	}
	b.nodes[index] = node
}

func (b *SparseBitmap64p3) Bytes() int {
	bytes := 8
	for _, node := range b.nodes {
		bytes += node.Bytes()
	}

	return bytes
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

func (b SparseBitmap) Or(with SparseBitmap) SparseBitmap {
	or := SparseBitmap{}
	indices := b.indices | with.indices
	for i := 0; i < 64 && indices != 0; i++ {
		if indices&1 != 0 {
			bits := b.bitsAt(i) | with.bitsAt(i)
			if bits != 0 {
				or.setBitsAt(i, bits)
			}
		}
		indices >>= 1
	}

	return or
}

func (b *SparseBitmap) IsZero() bool {
	return b.indices == 0
}

func (b *SparseBitmap) bitsAt(blockPosition int) Bitmap64 {
	if !b.indices.IsSet(blockPosition) {
		return 0
	}

	index := b.indices.Index(blockPosition)

	return b.bits[index]
}

func (b *SparseBitmap) setBitsAt(blockPosition int, bits Bitmap64) {
	index := b.indices.Index(blockPosition)
	if !b.indices.IsSet(blockPosition) {
		if len(b.bits) == index {
			b.bits = append(b.bits, 0)
		} else {
			b.bits = append(b.bits[:index+1], b.bits[index:]...)
		}
		b.indices.Set(blockPosition)
	}
	b.bits[index] = bits
}

func (b *SparseBitmap) Bytes() int {
	return 8 + len(b.bits)*8
}

func split6bits(n int) (int, int) {
	return n >> 6, n & 0b111111
}

func split12bits(n int) (int, int) {
	return n >> 12, n & 0b111111111111
}
