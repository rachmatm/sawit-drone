package service

import (
	"sort"

	"github.com/SawitProRecruitment/UserService/models"
)

type EstateStats struct {
	Count  int
	Min    int
	Max    int
	Median int
}

// Median adalah nilai tengah dari sekumpulan data yang telah diurutkan dari yang terkecil hingga yang terbesar.
// Median membagi data menjadi dua bagian yang sama besar,
// di mana setengah data nilainya lebih kecil dan setengahnya lagi lebih besar.
func CalculateStats(trees []models.Tree) EstateStats {
	if len(trees) == 0 {
		return EstateStats{}
	}

	heights := make([]int, len(trees))
	for i, t := range trees {
		heights[i] = t.Height
	}
	sort.Ints(heights)

	n := len(heights)
	median := 0
	if n%2 == 1 {
		median = heights[n/2]
	} else {
		median = (heights[n/2-1] + heights[n/2]) / 2
	}

	return EstateStats{
		Count:  n,
		Min:    heights[0],
		Max:    heights[n-1],
		Median: median,
	}
}
