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

func TestSparseBitmap64p3_Or(t *testing.T) {
	tests := []struct {
		bits1 []int
		bits2 []int
	}{
		{},
		{
			bits1: []int{0, 64, 127},
			bits2: []int{63, 64, 127, 128},
		},
		{
			bits1: []int{4095, 2048, 1024, 1023},
			bits2: []int{4095, 2047, 1024},
		},
		{
			bits1: []int{100_000, 8192, 4095, 2048, 1024, 1023},
			bits2: []int{200_000, 8192, 8191, 2048, 2047, 1024},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v | %v", test.bits1, test.bits2), func(t *testing.T) {
			bitmap1 := bitmaps.SparseBitmap64p3{}
			bitmap2 := bitmaps.SparseBitmap64p3{}
			m := make(map[int]struct{})

			for _, bit := range test.bits1 {
				bitmap1.Set(bit)
				m[bit] = struct{}{}
			}
			for _, bit := range test.bits2 {
				bitmap2.Set(bit)
				m[bit] = struct{}{}
			}
			bitmap := bitmap1.Or(bitmap2)

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

func TestSparseBitmap_Or(t *testing.T) {
	tests := []struct {
		bits1 []int
		bits2 []int
	}{
		{},
		{
			bits1: []int{0, 64, 127},
			bits2: []int{63, 64, 127, 128},
		},
		{
			bits1: []int{4095, 2048, 1024, 1023},
			bits2: []int{4095, 2047, 1024},
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v | %v", test.bits1, test.bits2), func(t *testing.T) {
			bitmap1 := bitmaps.SparseBitmap{}
			bitmap2 := bitmaps.SparseBitmap{}
			m := make(map[int]struct{})

			for _, bit := range test.bits1 {
				bitmap1.Set(bit)
				m[bit] = struct{}{}
			}
			for _, bit := range test.bits2 {
				bitmap2.Set(bit)
				m[bit] = struct{}{}
			}
			bitmap := bitmap1.Or(bitmap2)

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
