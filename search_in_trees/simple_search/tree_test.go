package simple_search_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/search_in_trees/simple_search"
	"github.com/strider2038/algos/search_in_trees/testcases"
)

func TestTree_Filter(t *testing.T) {
	tree := loadTree(t)

	for _, test := range testcases.SimpleCases {
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

	for _, bm := range testcases.SimpleCases {
		b.Run(bm.Name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				tree.Filter(bm.Search)
			}
		})
	}
}

func loadTree(tb testing.TB) *simple_search.Tree {
	_, filename, _, _ := runtime.Caller(0)
	tree, err := simple_search.LoadFromFile(filepath.Dir(filename) + "/../../testdata/classifiers/okved.csv")
	if err != nil {
		tb.Fatal("load okved:", err)
	}

	return tree
}