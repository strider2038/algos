package prefix_trees_test

import (
	"testing"

	"github.com/strider2038/algos/testdata/fixtures"
)

func BenchmarkMap_Fill(b *testing.B) {
	cities := fixtures.CitiesT(b)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		t := make(map[string]int, 0)
		for n, city := range cities {
			t[city] = n + 1
		}
	}
}

func BenchmarkMap_Get(b *testing.B) {
	cities := fixtures.CitiesT(b)
	t := make(map[string]int, 0)
	for n, city := range cities {
		t[city] = n + 1
	}

	b.ResetTimer()

	benchmarks := []struct {
		name     string
		cityName string
	}{
		{
			name:     "short name",
			cityName: "Adville",
		},
		{
			name:     "long name",
			cityName: "Advocate Lutheran General Childrens Hospital",
		},
		{
			name:     "long unique suffix",
			cityName: "Advocate Services Medical Transportation",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, found := t[bm.cityName]
				if !found {
					b.Fatal("element not found")
				}
			}
		})
	}
}
