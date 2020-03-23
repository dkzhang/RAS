package myUtils

import "math"

func FloatToInt(x float64) int {
	return int(math.Floor(x + 0.5))
}
