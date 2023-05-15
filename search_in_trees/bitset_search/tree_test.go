package bitset_search_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/search_in_trees/bitset_search"
	"github.com/strider2038/algos/search_in_trees/testcases"
)

func TestTree_Filter(t *testing.T) {
	tree := loadTree(t)

	for _, test := range testcases.PrefixCases {
		t.Run(test.Name, func(t *testing.T) {
			nodes := tree.Filter(test.Search)

			codes := make([]string, 0, len(nodes))
			for _, node := range nodes {
				codes = append(codes, node.Value.Code)
			}
			assert.Equal(t, test.Want, codes)
		})
	}
}

func BenchmarkTree_Filter(b *testing.B) {
	tree := loadTree(b)
	b.Log("size:", tree.Bytes()/1024, "Kb")

	for _, bm := range testcases.PrefixCases {
		b.Run(bm.Name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				tree.Filter(bm.Search)
			}
		})
	}
}

func loadTree(tb testing.TB) *bitset_search.Tree {
	_, filename, _, _ := runtime.Caller(0)
	tree, err := bitset_search.LoadFromFile(filepath.Dir(filename) + "/../../testdata/classifiers/okved.csv")
	if err != nil {
		tb.Fatal("load okved:", err)
	}

	return tree
}
