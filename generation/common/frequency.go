package common

import (
	"fmt"

	"github.com/unchartedsoftware/prism/tile"
)

// Frequency represents a tiling generator that produces heatmaps.
type Frequency struct {
	tile.Frequency
}

func (f *Frequency) CastFrequency(val interface{}) int64 {
	numF, isNum := val.(float64)
	if isNum {
		return int64(numF)
	}
	numI, isNum := val.(int64)
	if isNum {
		return numI
	}

	//TODO: Figure out which types are allowed, and what to do if bad data is received.
	return -1
}

func (f *Frequency) CastTime(val interface{}) interface{} {
	num, isNum := val.(float64)
	if isNum {
		return int64(num)
	}
	str, isStr := val.(string)
	if isStr {
		return str
	}
	return val
}

func (f *Frequency) CastTimeToString(val interface{}) string {
	num, isNum := val.(float64)
	if isNum {
		// assume milliseconds
		return fmt.Sprintf("%dms\n", int64(num))
	}
	str, isStr := val.(string)
	if isStr {
		return str
	}
	return ""
}
