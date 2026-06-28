package service

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/stretchr/testify/require"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name       string
		heights    []int
		wantCount  int
		wantMin    int
		wantMax    int
		wantMedian int
	}{
		{"empty", nil, 0, 0, 0, 0},
		{"single", []int{10}, 1, 10, 10, 10},
		{"odd count", []int{5, 10, 20}, 3, 5, 20, 10},
		{"even count", []int{10, 20, 30, 40}, 4, 10, 40, 25},
		{"even rounding", []int{10, 15}, 2, 10, 15, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var trees []models.Tree
			for _, h := range tt.heights {
				trees = append(trees, models.Tree{Height: h})
			}
			got := CalculateStats(trees)
			require.Equal(t, tt.wantCount, got.Count)
			require.Equal(t, tt.wantMin, got.Min)
			require.Equal(t, tt.wantMax, got.Max)
			require.Equal(t, tt.wantMedian, got.Median)
		})
	}
}
