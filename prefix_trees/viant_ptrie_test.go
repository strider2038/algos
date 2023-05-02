package prefix_trees

import (
	"testing"

	"github.com/strider2038/algos/testdata/fixtures"
	"github.com/viant/ptrie"
)

func BenchmarkViantPtrie_Fill(b *testing.B) {
	cities := fixtures.CitiesT(b)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		t := ptrie.New()
		for n, city := range cities {
			if err := t.Put([]byte(city), n+1); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkViantPtrie_Get(b *testing.B) {
	cities := fixtures.CitiesT(b)
	t := ptrie.New()
	for n, city := range cities {
		if err := t.Put([]byte(city), n+1); err != nil {
			b.Fatal(err)
		}
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
				_, found := t.Get([]byte(bm.cityName))
				if !found {
					b.Fatal("element not found")
				}
			}
		})
	}
}
