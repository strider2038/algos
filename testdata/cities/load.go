// cities - data source: https://www.geonames.org/datasources/

package cities

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func LoadT(tb testing.TB) []string {
	cities, err := Load()
	if err != nil {
		tb.Fatal(err)
	}

	return cities
}

func Load() ([]string, error) {
	_, filename, _, _ := runtime.Caller(0)

	input, err := os.Open(filepath.Dir(filename) + "/cities.txt")
	if err != nil {
		return nil, err
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	cities := make([]string, 0)
	for scanner.Scan() {
		cities = append(cities, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}
