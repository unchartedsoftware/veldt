package util

import (
	"math"
)

// Fract returns the fractional part of a floating point value
func Fract(value float64) float64 {
	_, fract := math.Modf(value)
	return fract
}

// Round rounds the provided floating point value
func Round(f float64) float64 {
	return math.Floor(f + 0.5)
}
