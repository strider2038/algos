package byte_suffix_trie_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/prefix_trees/byte_suffix_trie"
	"github.com/strider2038/algos/testdata/fixtures"
)

func TestArray_Basic(t *testing.T) {
	items := byte_suffix_trie.Array[int]{}

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

func TestArray_Put_Suffixes(t *testing.T) {
	tests := []struct {
		name      string
		keys      [][]byte
		wantCount int
	}{
		{
			name: "split by last byte of suffix",
			keys: [][]byte{
				{1, 2, 3, 4, 5},
				{1, 2, 3, 4, 6},
			},
			wantCount: 2,
		},
		{
			name: "split with long suffix",
			keys: [][]byte{
				{1, 2, 3, 4, 5},
				{1, 2, 3, 4, 5, 6, 7, 8},
			},
			wantCount: 2,
		},
		{
			name: "split by prefix",
			keys: [][]byte{
				{1, 2, 3, 4},
				{1, 2},
			},
			wantCount: 2,
		},
		{
			name: "small key",
			keys: [][]byte{
				{1, 2, 3},
				{1},
			},
			wantCount: 2,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			items := byte_suffix_trie.Array[int]{}

			for i, key := range test.keys {
				items.Put(key, i+1)
			}

			assert.Equal(t, test.wantCount, items.Count())
			for i, key := range test.keys {
				value, exists := items.Find(key)
				if !exists {
					t.Errorf("key not found: %v", key)
				}
				assert.Equal(t, i+1, value)
			}
		})
	}
}

func TestArray_Put_Countries(t *testing.T) {
	countries := byte_suffix_trie.Array[int]{}
	m := map[string]int{}

	for i, country := range fixtures.Countries {
		countries.Put([]byte(country), i+1)
		m[country] = i + 1
	}

	_ = countries.Walk(func(key []byte, value int) error {
		assert.Equal(t, m[string(key)], value)

		return nil
	})
}

func TestArray_Put_RandomStrings(t *testing.T) {
	const count = 100_000
	tree := byte_suffix_trie.Array[int]{}
	m := map[string]int{}

	ss := randomStrings(15, count)
	for i, s := range ss {
		tree.Put([]byte(s), i)
		m[s] = i
	}

	isFailed := false
	for key, value := range m {
		v, ok := tree.Find([]byte(key))
		if !ok {
			t.Error("key not found:", key)
			isFailed = true
		}
		if !assert.Equal(t, value, v, "at key: %s", key) {
			isFailed = true
		}
	}

	if isFailed {
		for _, s := range ss {
			fmt.Println(s)
		}
	}
}

func TestArray_MarshalJSON(t *testing.T) {
	items := byte_suffix_trie.Array[int]{}
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

func TestArray_WalkPrefix(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		prefix string
		want   []string
	}{
		{
			name:   "single value",
			values: []string{"foo"},
			prefix: "fo",
			want:   []string{"foo"},
		},
		{
			name:   "not existing prefix",
			values: []string{"foo"},
			prefix: "fos",
			want:   []string{},
		},
		{
			name:   "no suffixes",
			values: []string{"cat", "cap", "car", "foo", "bar"},
			prefix: "ca",
			want:   []string{"cap", "car", "cat"},
		},
		{
			name:   "too long prefix",
			values: []string{"cap", "car", "cat"},
			prefix: "catapult",
			want:   []string{},
		},
		{
			name:   "empty prefix",
			values: []string{"cat", "cap", "car", "foo", "bar"},
			prefix: "",
			want:   []string{"bar", "cap", "car", "cat", "foo"},
		},
		{
			name: "suffix case",
			values: []string{
				"capacity",
			},
			prefix: "cap",
			want:   []string{"capacity"},
		},
		{
			name: "case 1",
			values: []string{
				"QtGgCdh8S",
				"QjelrTqoqGZV",
				"QNaEhkK9E",
				"Q1iq coOLuBe5c",
				"GQE9WruzR1p8",
			},
			prefix: "Qj",
			want:   []string{"QjelrTqoqGZV"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tree := byte_suffix_trie.Array[int]{}
			for i, value := range test.values {
				tree.Put([]byte(value), i)
			}

			got := make([]string, 0)
			_ = tree.WalkPrefix([]byte(test.prefix), func(key []byte, value int) error {
				got = append(got, string(key))
				return nil
			})
			assert.Equal(t, test.want, got)
		})
	}
}

func TestArray_WalkPrefix_Random(t *testing.T) {
	const (
		itemsCount = 10_000
		testCount  = 100
	)
	tree := byte_suffix_trie.Array[int]{}
	m := map[string]int{}
	ss := randomStrings(15, itemsCount)
	for i, s := range ss {
		tree.Put([]byte(s), i)
		m[s] = i
	}

	for i := 0; i < testCount; i++ {
		s := ""
		for key := range m {
			s = key
			break
		}
		prefix := s[:rand.Intn(len(s))]
		t.Run(prefix, func(t *testing.T) {
			got := make([]string, 0)
			_ = tree.WalkPrefix([]byte(prefix), func(key []byte, value int) error {
				got = append(got, string(key))
				return nil
			})

			want := make([]string, 0)
			for key := range m {
				if strings.HasPrefix(key, prefix) {
					want = append(want, key)
				}
			}
			sort.Slice(want, func(i, j int) bool {
				return want[i] < want[j]
			})

			assert.Equal(t, want, got)
		})
	}
}

func BenchmarkArray64_Fill(b *testing.B) {
	cities := fixtures.CitiesT(b)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		t := byte_suffix_trie.Array[int]{}
		for n, city := range cities {
			t.Put([]byte(city), n+1)
		}
	}
}

func BenchmarkArray64_Get(b *testing.B) {
	cities := fixtures.CitiesT(b)
	t := byte_suffix_trie.Array[int]{}
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
