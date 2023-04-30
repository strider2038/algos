package byte_suffix_trie

import "math/bits"

// bitIndex - битовая маска для хранения 64 индексов.
type bitIndex [4]uint64

func (b *bitIndex) set(n byte) {
	hi, lo := b.splitN(n)
	b[hi] = b[hi] | (1 << lo)
}

func (b *bitIndex) unset(n byte) {
	hi, lo := b.splitN(n)
	b[hi] = b[hi] & ^(1 << lo)
}

func (b *bitIndex) isSet(n byte) bool {
	hi, lo := b.splitN(n)

	return b[hi]&(1<<lo) != 0
}

// getOneNumber возвращает порядковый номер установленного бита. Перед вызовом функции
// необходимо обязательно проверить установлен ли бит с помощью функции isSet.
//
// Пример маски и номеров
//
//	маска             0 0 1 0 0 1 1 0
//	номер бита        7 6 5 4 3 2 1 0
//	порядковый номер  - - 2 - - 1 0 -
//
// Примеры:
//
//	маска bitIndex = 0010 0110, номер бита n = 1, вернется число 0
//	маска bitIndex = 0010 0110, номер бита n = 2, вернется число 1
//	маска bitIndex = 0010 0110, номер бита n = 6, вернется число 2
func (b *bitIndex) getOneNumber(n byte) int {
	hi, lo := b.splitN(n)

	index := bits.OnesCount64(b[hi] & ^(uint64(0xFFFFFFFFFFFFFFFF) << lo))
	for i := byte(0); i < hi; i++ {
		index += bits.OnesCount64(b[i])
	}

	return index
}

func (b *bitIndex) splitN(n byte) (byte, byte) {
	return n >> 6, n & 0x3F
}
