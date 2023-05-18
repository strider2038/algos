package bitmaps_test

import (
	"testing"

	"github.com/strider2038/algos/bitmaps"
)

func BenchmarkSparseBitmap_IsSet(b *testing.B) {
	bitmap := make(bitmaps.Bitmap, 0)
	for i := 0; i < 64*64*64; i += 4 {
		bitmap.Set(i)
	}
	bitmap.Set(123456)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bitmap.IsSet(123456)
	}
}
