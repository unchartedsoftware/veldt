package util

import (
    "math"
)

// Fract returns the fractional part of a floating point value
func Fract( value float64 ) float64 {
	_, fract := math.Modf( value )
	return fract
}
