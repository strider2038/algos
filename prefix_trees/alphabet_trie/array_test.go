package alphabet_trie_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/prefix_trees/alphabet_trie"
	"github.com/strider2038/algos/testdata/fixtures"
)

var testChars = []rune("abcdefghijklmnopqrstuvwxyz")

func TestArray64_Basic(t *testing.T) {
	items := alphabet_trie.NewArray64[int]("abcdefghijklmnopqrstuvwxyz")

	items.Put("alpha", 1)
	items.Put("beta", 2)
	items.Put("gamma", 3)
	items.Put("delta", 4)
	items.Delete("beta")
	items.Put("beta", 5)
	items.Put("cap", 6)
	items.Put("cat", 7)
	items.Put("car", 8)
	items.Delete("delta")
	items.Delete("delta")
	items.Delete("unknown")

	assert.Equal(t, 6, items.Count())
	assert.Equal(t, 1, items.Get("alpha"))
	assert.Equal(t, 5, items.Get("beta"))
	assert.Equal(t, 3, items.Get("gamma"))
	assert.Equal(t, 6, items.Get("cap"))
	assert.Equal(t, 7, items.Get("cat"))
	assert.Equal(t, 8, items.Get("car"))
	assert.Equal(t, 0, items.Get("delta"))
	if _, exist := items.Find("delta"); exist {
		t.Error("delta value is found in map")
	}
}

func TestArray64_RealData(t *testing.T) {
	countries := alphabet_trie.NewArray64[int](`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "-&`)
	m := map[string]int{}

	for i, country := range fixtures.Countries {
		countries.Put(country, i+1)
		m[country] = i + 1
	}

	countries.Walk(func(key string, value int) error {
		assert.Equal(t, m[key], value)

		return nil
	})
}

func TestArray64_MarshalJSON(t *testing.T) {
	items := alphabet_trie.NewArray64[int]("abcdefghijklmnopqrstuvwxyz")
	items.Put("alpha", 1)
	items.Put("beta", 2)
	items.Put("gamma", 3)
	items.Put("delta", 4)

	data, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, `{"alpha":1,"beta":2,"delta":4,"gamma":3}`, string(data))
}

func BenchmarkArray64_Fill(b *testing.B) {
	cities := fixtures.CitiesT(b)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		t := alphabet_trie.NewArray64[int](`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "-&`)
		for n, city := range cities {
			t.Put(city, n+1)
		}
	}
}

func BenchmarkArray64_Get(b *testing.B) {
	cities := fixtures.CitiesT(b)
	t := alphabet_trie.NewArray64[int](`abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "-&`)
	for n, city := range cities {
		t.Put(city, n+1)
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
				_, found := t.Find(bm.cityName)
				if !found {
					b.Fatal("element not found")
				}
			}
		})
	}
}
