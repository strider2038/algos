package simple_search_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/strider2038/algos/search_in_trees/simple_search"
)

func TestTree_Filter(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	tree, err := simple_search.LoadFromFile(filepath.Dir(filename) + "/../../testdata/classifiers/okved.csv")
	if err != nil {
		t.Fatal("load okved:", err)
	}

	tests := []struct {
		name   string
		search string
		want   []string
	}{
		{
			name:   "basic",
			search: "выращиван",
			want:   []string{"A"},
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
