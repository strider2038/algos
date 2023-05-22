package bitmap_search_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/search_in_trees/bitmap_search"
	"github.com/strider2038/algos/search_in_trees/testcases"
)

func TestTree_Filter(t *testing.T) {
	tree, err := loadTree("testdata/classifiers/okved.csv")
	if err != nil {
		t.Fatal(err)
	}

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
	tree, err := loadTree("testdata/classifiers/okved.csv")
	if err != nil {
		b.Fatal(err)
	}
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

func BenchmarkTree_Filter_BigDataSet(b *testing.B) {
	tree, err := loadTree("var/ksr.csv")
	if err != nil {
		b.Skip(err)
	}
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

func loadTree(dataFilename string) (*bitmap_search.Tree, error) {
	_, goFilename, _, _ := runtime.Caller(0)
	tree, err := bitmap_search.LoadFromFile(filepath.Dir(goFilename) + "/../../" + dataFilename)
	if err != nil {
		return nil, fmt.Errorf(`load classifiers from "%s": %w`, dataFilename, err)
	}

	return tree, nil
}
