package byte_trie_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/prefix_trees/byte_trie"
	"github.com/strider2038/algos/testdata/fixtures"
)

func TestArray_Basic(t *testing.T) {
	items := byte_trie.Array[int]{}

	items.Put([]byte("alpha"), 1)
	items.Put([]byte("beta"), 2)
	items.Put([]byte("gamma"), 3)
	items.Put([]byte("delta"), 4)
	items.Delete([]byte("beta"))
	items.Put([]byte("beta"), 5)
	items.Put([]byte("cap"), 6)
	items.Put([]byte("cat"), 7)
	items.Put([]byte("car"), 8)
	items.Delete([]byte("delta"))
	items.Delete([]byte("delta"))
	items.Delete([]byte("unknown"))

	assert.Equal(t, 6, items.Count())
	assert.Equal(t, 1, items.Get([]byte("alpha")))
	assert.Equal(t, 5, items.Get([]byte("beta")))
	assert.Equal(t, 3, items.Get([]byte("gamma")))
	assert.Equal(t, 6, items.Get([]byte("cap")))
	assert.Equal(t, 7, items.Get([]byte("cat")))
	assert.Equal(t, 8, items.Get([]byte("car")))
	assert.Equal(t, 0, items.Get([]byte("delta")))
	if _, exist := items.Find([]byte("delta")); exist {
		t.Error("delta value is found in map")
	}
}

func TestArray_Put_Countries(t *testing.T) {
	countries := byte_trie.Array[int]{}
	m := map[string]int{}

	for i, country := range fixtures.Countries {
		countries.Put([]byte(country), i+1)
		m[country] = i + 1
	}

	countries.Walk(func(key []byte, value int) error {
		assert.Equal(t, m[string(key)], value)

		return nil
	})
}

func TestArray_Put_RandomStrings(t *testing.T) {
	const count = 100_000
	tree := byte_trie.Array[int]{}
	m := map[string]int{}

	ss := randomStrings(15, count)
	for i, s := range ss {
		tree.Put([]byte(s), i)
		m[s] = i
	}

	for key, value := range m {
		v, ok := tree.Find([]byte(key))
		if !ok {
			t.Error("key not found:", key)
		}
		assert.Equal(t, value, v)
	}
}

func TestArray_MarshalJSON(t *testing.T) {
	items := byte_trie.Array[int]{}
	items.Put([]byte("alpha"), 1)
	items.Put([]byte("beta"), 2)
	items.Put([]byte("gamma"), 3)
	items.Put([]byte("delta"), 4)

	data, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEq(t, `{"alpha":1,"beta":2,"delta":4,"gamma":3}`, string(data))
}

func BenchmarkArray64_Fill(b *testing.B) {
	cities := fixtures.CitiesT(b)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		t := byte_trie.Array[int]{}
		for n, city := range cities {
			t.Put([]byte(city), n+1)
		}
	}
}

func BenchmarkArray64_Get(b *testing.B) {
	cities := fixtures.CitiesT(b)
	t := byte_trie.Array[int]{}
	for n, city := range cities {
		t.Put([]byte(city), n+1)
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
				_, found := t.Find([]byte(bm.cityName))
				if !found {
					b.Fatal("element not found")
				}
			}
		})
	}
}
