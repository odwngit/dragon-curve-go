package main

import (
	"fmt"
	"slices"
	"github.com/fogleman/gg"
)

type Point struct {
	x, y int
}

func main() {
	var iterations int = 24
	const draw_scale int = 2
	defer fmt.Println("Finished!")

	// fmt.Printf("Calculating turns with %v iterations...\n", iterations)
	
	// byte 1 = right turn
	// byte 2 = left turn
	// byte 0 = unassigned

	// This is how you define a slice. Slices are basically dynamic arrays/vectors.
	// You can dynamically allocate them using make([]type, length)

	var turns = []byte{1}

	for iterations > 0 {
		addition := slices.Clone(turns)
		slices.Reverse(addition) // Reverse copy
		for i, v := range addition { // Invert copy (1=2, 2=1)
			if v == 1 {
				addition[i] = 2
			} else if v == 2 {
				addition[i] = 1
			}
		}

		turns = slices.Concat(turns, []byte{1}, addition) // Add right turn and additions
		iterations--
	}

	// fmt.Printf("Turns: %v\n", len(turns))

	points := make([]Point, len(turns)+1)
	points[0] = Point{0, 0}

	var position Point = Point{0, 0}
	var facing byte = 1 // 1 up, 2 right, 3 down, 4 left

	var bounds_min Point = Point{0, 0}
	var bounds_max Point = Point{0, 0}

	for i, v := range turns {
		if v == 1 { // Turn facing based on turn
			facing++
		} else if v == 2 {
			facing--
		}

		if facing < 1 { // Wrap around facing
			facing = 4
		} else if facing > 4 {
			facing = 1
		}

		switch facing { // Move position
			case 1:
				position.y -= draw_scale
			case 2:
				position.x += draw_scale
			case 3:
				position.y += draw_scale
			case 4:
				position.x -= draw_scale
		}
		points[i+1] = position

		if position.x > bounds_max.x { // Update bounds
			bounds_max.x = position.x
		} else if position.x < bounds_min.x {
			bounds_min.x = position.x
		}
		if position.y > bounds_max.y {
			bounds_max.y = position.y
		} else if position.y < bounds_min.y {
			bounds_min.y = position.y
		}

	}

	// Adjust points to not be in negatives
	for i, _ := range points {
		points[i].x -= bounds_min.x
		points[i].y -= bounds_min.y
	}

	dc := gg.NewContext(bounds_max.x - bounds_min.x, bounds_max.y - bounds_min.y)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.MoveTo(float64(points[0].x), float64(points[0].y))

	for i := 0; i < len(points)-1; i++ {
		// 2 lines below cause weird antialiasing
		//dc.DrawLine(float64(points[i].x), float64(points[i].y), float64(points[i+1].x), float64(points[i+1].y))
		//dc.Stroke()

		if points[i].x == points[i+1].x { // Then the next point is somewhere vertically
			if points[i].y > points[i+1].y { // Then the next point is above current
				for points[i].y != points[i+1].y {
					points[i].y--			
					dc.SetPixel(points[i].x, points[i].y)
				}
			} else { // Then the next point is below current
				for points[i].y != points[i+1].y {
					points[i].y++		
					dc.SetPixel(points[i].x, points[i].y)
				}
			}
		} else { // Then the next point is somewhere horizontally
			if points[i].x > points[i+1].x { // Then the next point is left of current
				for points[i].x != points[i+1].x {
					points[i].x--
					dc.SetPixel(points[i].x, points[i].y)
				}
			} else { // Then the next point is right of current
				for points[i].x != points[i+1].x {
					points[i].x++
					dc.SetPixel(points[i].x, points[i].y)
				}
			}
		}
	}

	dc.SavePNG("output.png")
}
