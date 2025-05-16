package utils

import "math"

func Round(value float64, count int) float64 {
	pow := math.Pow(10, float64(count))
	return math.Round(value*pow) / pow
}
