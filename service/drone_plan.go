package service

import (
	"fmt"
	"math"

	"github.com/SawitProRecruitment/UserService/models"
)

type Point struct {
	X int
	Y int
}

type DronePlan struct {
	Distance int
	RestX    int
	RestY    int
}

func GenerateFlightPath(
	width int,
	length int,
) []Point {

	path := make([]Point, 0, width*length)

	for y := 1; y <= length; y++ {

		if y%2 == 1 {

			for x := 1; x <= width; x++ {
				path = append(path, Point{
					X: x,
					Y: y,
				})
			}

		} else {

			for x := width; x >= 1; x-- {
				path = append(path, Point{
					X: x,
					Y: y,
				})
			}
		}
	}

	return path
}

func CalculateDronePlan(
	estate *models.Estate,
	trees []models.Tree,
	maxDistance *int,
) DronePlan {

	treeMap := make(map[string]int)

	for _, tree := range trees {

		key := fmt.Sprintf(
			"%d:%d",
			tree.X,
			tree.Y,
		)

		treeMap[key] = tree.Height
	}

	path := GenerateFlightPath(
		estate.Width,
		estate.Length,
	)

	horizontalDistance := 0
	verticalDistance := 0
	currentHeight := 0
	travelled := 0
	lastPlotX := 1
	lastPlotY := 1
	plotDistance := 10 // Each plot is 10 meters apart
	prevX := 1
	prevY := 1

	for i, point := range path {

		lastPlotX = point.X
		lastPlotY = point.Y

		// Default target height is 1 if there is no tree at the current point
		targetHeight := 1

		key := fmt.Sprintf(
			"%d:%d",
			point.X,
			point.Y,
		)

		// Check if the current point has a tree and get its height
		if treeHeight, ok := treeMap[key]; ok {
			targetHeight = treeHeight + 1
		}

		//if next tree is higher than current height, we need to move up,
		// else we need to move down
		verticalMove := int(
			math.Abs(
				float64(
					targetHeight - currentHeight,
				),
			),
		)

		// Add the vertical move to the total vertical distance and travelled distance
		verticalDistance += verticalMove
		travelled += verticalMove
		// Update the current height to the target height for the next iteration
		currentHeight = targetHeight

		// Add horizontal distance for each move except the first one
		if i > 0 {
			// Each plot is 10 meters apart,
			// so we add 10 meters to the horizontal distance for each move
			horizontalDistance += plotDistance
			travelled += plotDistance

			prevX = path[i-1].X
			prevY = path[i-1].Y
		}

		//check maxDistance
		if maxDistance != nil {
			if *maxDistance <= travelled {
				restX, restY := lastPlotX, lastPlotY
				if *maxDistance <= travelled-(plotDistance/2) {
					restX, restY = prevX, prevY
				}
				return DronePlan{
					Distance: *maxDistance,
					RestX:    restX,
					RestY:    restY,
				}
			}
		}
	} //end of for loop

	// Add the final vertical distance to return to ground level (height 0)
	verticalDistance += currentHeight
	// Add the final vertical distance to the travelled distance
	totalDistance := horizontalDistance + verticalDistance

	result := DronePlan{
		Distance: totalDistance,
	}

	if maxDistance != nil {
		result.RestX = lastPlotX
		result.RestY = lastPlotY
	}

	return result
}
