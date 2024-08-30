package iterdun

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelSimple(t *testing.T) {
	items1 := slices.Values([]int{1, 2, 3})
	items2 := slices.Values([]int{4, 5, 6, 7})

	ctx := context.Background()
	result := slices.Collect(Parallel(ctx, items1, items2))
	slices.Sort(result)
	expected := []int{1, 2, 3, 4, 5, 6, 7}
	assert.Equal(t, expected, result)
}

func TestParallelCancel(t *testing.T) {
	items1 := slices.Values([]int{1, 2, 3, 4, 5})

	ctx, cancel := context.WithCancel(context.Background())

	var result []int
	for item := range Parallel(ctx, items1) {
		if item == 3 {
			cancel()
		}
		result = append(result, item)
	}

	expected := []int{1, 2, 3}
	assert.Equal(t, expected, result)

}
