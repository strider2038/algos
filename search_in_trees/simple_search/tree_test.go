package simple_search_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/search_in_trees/simple_search"
)

func TestTree_Filter(t *testing.T) {
	tree := loadTree(t)

	tests := []struct {
		name   string
		search string
		want   []string
	}{
		{
			name:   "top node filter, single word",
			search: "рыболовство",
			want:   []string{"A"},
		},
		{
			name:   "deep node filter, single word",
			search: "картофел",
			want:   []string{"A", "C", "G"},
		},
		{
			name:   "deep node filter, single word with stemming",
			search: "картофель",
			want:   []string{"A", "C", "G"},
		},
		{
			name:   "deep node filter, many words",
			search: "Деятельность прочих общественных и прочих некоммерческих организаций, кроме религиозных и политических организаций, не включенных в другие группировки",
			want:   []string{"S"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nodes := tree.Filter(test.search)

			codes := make([]string, 0, len(nodes))
			for _, node := range nodes {
				codes = append(codes, node.Value.Code)
			}
			assert.Equal(t, test.want, codes)
		})
	}
}

func BenchmarkTree_Filter(b *testing.B) {
	tree := loadTree(b)

	benchmarks := []struct {
		name   string
		search string
	}{
		{
			name:   "top node filter, single word",
			search: "рыболовство",
		},
		{
			name:   "deep node filter, single word",
			search: "картофел",
		},
		{
			name:   "deep node filter, single word with stemming",
			search: "картофель",
		},
		{
			name:   "deep node filter, many words",
			search: "Деятельность прочих общественных и прочих некоммерческих организаций, кроме религиозных и политических организаций, не включенных в другие группировки",
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				tree.Filter(bm.search)
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
