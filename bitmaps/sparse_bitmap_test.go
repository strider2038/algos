package bitmaps_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/bitmaps"
)

func TestSparseBitmap64p3_IsSet(t *testing.T) {
	tests := []struct {
		bits []int
	}{
		{bits: []int{}},
		{bits: []int{0, 63, 64}},
		{bits: []int{4095, 2048, 2047, 1024, 1023}},
		{bits: []int{200_000, 100_000, 8192, 8191, 4095, 2048, 2047, 1024, 1023}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.bits), func(t *testing.T) {
			bitmap := bitmaps.SparseBitmap64p3{}
			m := make(map[int]struct{})

			for _, bit := range test.bits {
				bitmap.Set(bit)
				m[bit] = struct{}{}
			}

			for i := 0; i <= 64*64*64; i++ {
				if _, isSet := m[i]; isSet {
					assert.True(t, bitmap.IsSet(i), "at index: %d", i)
				} else {
					assert.False(t, bitmap.IsSet(i), "at index: %d", i)
				}
			}
		})
	}
}

func TestSparseBitmap_IsSet(t *testing.T) {
	tests := []struct {
		bits []int
	}{
		{bits: []int{}},
		{bits: []int{0, 63, 64}},
		{bits: []int{4095, 2048, 2047, 1024, 1023}},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.bits), func(t *testing.T) {
			bitmap := bitmaps.SparseBitmap{}
			m := make(map[int]struct{})

			for _, bit := range test.bits {
				bitmap.Set(bit)
				m[bit] = struct{}{}
			}

			for i := 0; i <= 64*64; i++ {
				if _, isSet := m[i]; isSet {
					assert.True(t, bitmap.IsSet(i), "at index: %d", i)
				} else {
					assert.False(t, bitmap.IsSet(i), "at index: %d", i)
				}
			}
		})
	}
}

func BenchmarkSparseBitmap64p3_IsSet(b *testing.B) {
	bitmap := bitmaps.SparseBitmap64p3{}
	for i := 0; i < 64*64*64; i += 4 {
		bitmap.Set(i)
	}
	bitmap.Set(123456)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bitmap.IsSet(123456)
	}
}
