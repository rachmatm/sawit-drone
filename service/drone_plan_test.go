package service

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func makeEstate(width, length int) *models.Estate {
	return &models.Estate{Width: width, Length: length}
}

func intPtr(v int) *int {
	return &v
}

// ---------------------------------------------------------------------------
// GenerateFlightPath
// ---------------------------------------------------------------------------

func TestGenerateFlightPath(t *testing.T) {
	t.Run("1x1 estate returns single point", func(t *testing.T) {
		path := GenerateFlightPath(1, 1)
		require.Len(t, path, 1)
		require.Equal(t, Point{X: 1, Y: 1}, path[0])
	})

	t.Run("6x6 estate returns single point", func(t *testing.T) {
		path := GenerateFlightPath(6, 6)
		t.Logf("path: %+v", path)
		require.Len(t, path, 36)
		require.Equal(t, Point{X: 1, Y: 6}, path[35])
	})

	t.Run("total points equals width x length", func(t *testing.T) {
		path := GenerateFlightPath(5, 3)
		require.Len(t, path, 15)
	})

	t.Run("odd row goes left to right", func(t *testing.T) {
		path := GenerateFlightPath(3, 1) // row 1 is odd
		require.Equal(t, Point{X: 1, Y: 1}, path[0])
		require.Equal(t, Point{X: 2, Y: 1}, path[1])
		require.Equal(t, Point{X: 3, Y: 1}, path[2])
	})

	t.Run("even row goes right to left (boustrophedon)", func(t *testing.T) {
		path := GenerateFlightPath(3, 2) // row 2 is even
		require.Equal(t, Point{X: 3, Y: 2}, path[3])
		require.Equal(t, Point{X: 2, Y: 2}, path[4])
		require.Equal(t, Point{X: 1, Y: 2}, path[5])
	})

	t.Run("3x2 estate full path order", func(t *testing.T) {
		path := GenerateFlightPath(3, 2)
		expected := []Point{
			{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, // row 1: left → right
			{X: 3, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 2}, // row 2: right → left
		}
		require.Equal(t, expected, path)
	})
}

// ---------------------------------------------------------------------------
// CalculateDronePlan
// ---------------------------------------------------------------------------

func TestCalculateDronePlan(t *testing.T) {

	t.Run("1x1 estate no trees: rise 1 + land 1 = 2", func(t *testing.T) {
		// no trees → target height at (1,1) = 1
		// vertical: abs(1-0) = 1 up, then land 1 down = 2
		// horizontal: 0 (only one plot, no movement)
		// total = 2
		plan := CalculateDronePlan(makeEstate(1, 1), nil, nil)
		require.Equal(t, 2, plan.Distance)
	})

	t.Run("no max distance: rest coordinates are zero", func(t *testing.T) {
		plan := CalculateDronePlan(makeEstate(1, 1), nil, nil)
		require.Equal(t, 0, plan.RestX)
		require.Equal(t, 0, plan.RestY)
	})

	t.Run("single tree raises drone to treeHeight+1", func(t *testing.T) {
		// tree height 10 → target height = 11
		// vertical: 11 up + 11 down = 22
		// horizontal: 0
		// total = 22
		trees := []models.Tree{{X: 1, Y: 1, Height: 10}}
		plan := CalculateDronePlan(makeEstate(1, 1), trees, nil)
		require.Equal(t, 22, plan.Distance)
	})

	t.Run("2x1 estate no trees: horizontal 10 + vertical 2 = 12", func(t *testing.T) {
		// path: (1,1) → (2,1)
		// at (1,1): rise 1, travelled=1
		// at (2,1): horizontal +10, travelled=11, target=1, vertical abs(1-1)=0
		// land: +1
		// total = 10 + 1 + 1 = 12
		plan := CalculateDronePlan(makeEstate(2, 1), nil, nil)
		require.Equal(t, 12, plan.Distance)
	})

	t.Run("drone adjusts height between trees of different heights", func(t *testing.T) {
		// estate 3x1, trees at (1,1)h=5 and (3,1)h=10
		// path: (1,1) → (2,1) → (3,1)
		// (1,1): target=6,  vertical abs(6-0)=6,  travelled=6
		// (2,1): target=1,  vertical abs(1-6)=5,  horizontal=10, travelled=21
		// (3,1): target=11, vertical abs(11-1)=10, horizontal=10, travelled=41
		// land: +11
		// total = (6+5+10) vertical + (10+10) horizontal + 11 land = 52
		trees := []models.Tree{
			{X: 1, Y: 1, Height: 5},
			{X: 3, Y: 1, Height: 10},
		}
		plan := CalculateDronePlan(makeEstate(3, 1), trees, nil)
		require.Equal(t, 52, plan.Distance)
	})

	t.Run("normal case: 5x1 estate with 3 trees returns distance 82", func(t *testing.T) {
		// mirrors the api_test "Normal 2" scenario
		trees := []models.Tree{
			{X: 2, Y: 1, Height: 10},
			{X: 3, Y: 1, Height: 20},
			{X: 4, Y: 1, Height: 10},
		}
		plan := CalculateDronePlan(makeEstate(5, 1), trees, nil)
		require.Equal(t, 82, plan.Distance)
	})

	t.Run("max distance hit mid-path: returns rest point and capped distance", func(t *testing.T) {
		// 6x6 estate, no trees
		// (1,1): rise 1, travelled=1
		// (2,1): horizontal +10, travelled=11, vertical 0
		// maxDistance=27 → last flight above plot 4
		// Claude responded: (1,1) 1-2-6, (2,1) 7-11-16, (3,1) 17-21-26,
		// (1,1) 1-2-6, (2,1) 7-11-16, (3,1) 17-21-26, (4,1) 27-31-36, (5,1) 37-41-46, (6,1) 47-51-56,
		// (6,2) 57-61-66, (5,2) 67-71-76, (4,2) 77-81-86, (3,2) 87-91-96, (2,2) 97-101-106,
		// (1,2) 107-111-116, (1,3) 117-121-126, (2,3) 127-131-136, (3,3) 137-141-146, (3,4) 147-148-150
		maxDistance := 141
		plan := CalculateDronePlan(makeEstate(6, 6), nil, intPtr(maxDistance))
		require.Equal(t, maxDistance, plan.Distance)
		require.Equal(t, 3, plan.RestX)
		require.Equal(t, 3, plan.RestY)
	})

	t.Run("max distance larger than total: travels full path, rest is last plot", func(t *testing.T) {
		// maxDistance larger than total distance → completes full path
		// rest coordinates = last plot visited
		fullPlan := CalculateDronePlan(makeEstate(2, 1), nil, nil)
		limitedPlan := CalculateDronePlan(makeEstate(2, 1), nil, intPtr(9999))
		require.Equal(t, fullPlan.Distance, limitedPlan.Distance)
		// when maxDistance is set and path completes, rest = last plot
		require.Equal(t, 2, limitedPlan.RestX)
		require.Equal(t, 1, limitedPlan.RestY)
	})

	t.Run("max distance exactly equals total: rest is last plot", func(t *testing.T) {
		// 1x1, no trees: total = 2
		// maxDistance = 2 → travelled hits exactly at (1,1) after landing
		plan := CalculateDronePlan(makeEstate(1, 1), nil, intPtr(2))
		require.Equal(t, 2, plan.Distance)
		require.Equal(t, 1, plan.RestX)
		require.Equal(t, 1, plan.RestY)
	})

	t.Run("tree not on flight path does not affect distance", func(t *testing.T) {
		// 1x1 estate, tree placed outside the single plot (2,1) — not visited
		// distance should be same as no trees: 2
		trees := []models.Tree{{X: 2, Y: 1, Height: 50}}
		plan := CalculateDronePlan(makeEstate(1, 1), trees, nil)
		require.Equal(t, 2, plan.Distance)
	})

	t.Run("multiple trees on same row: drone adjusts height per plot", func(t *testing.T) {
		// estate 2x1, tree at (1,1)h=5 and (2,1)h=5
		// (1,1): target=6, vertical abs(6-0)=6, travelled=6
		// (2,1): target=6, vertical abs(6-6)=0, horizontal=10, travelled=16
		// land: +6
		// total = 6 + 0 + 10 + 6 = 22
		trees := []models.Tree{
			{X: 1, Y: 1, Height: 5},
			{X: 2, Y: 1, Height: 5},
		}
		plan := CalculateDronePlan(makeEstate(2, 1), trees, nil)
		require.Equal(t, 22, plan.Distance)
	})
}
