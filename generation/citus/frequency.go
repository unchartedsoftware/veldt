package citus

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/prism/tile"
)

// Frequency represents a tiling generator that produces heatmaps.
type Frequency struct {
	tile.Frequency
}

type FrequencyResult struct {
	Bucket int
	Value  float64
}

func (f *Frequency) AddAggs(query *Query) *Query {
	//Bounds extension (empty buckets) will be done in the go code when parsing results
	//Not 100% sure if we need to substract the min value from the frequency field to
	//set the initial bucket.
	//Ex:
	//	data starts at 3, with intervals of 5.
	//	Should the first bucket be 0-5 or 3-8?

	//Ignoring potential error. Should really be done in some kind of setup function.
	intervalNum, _ := strconv.ParseFloat(f.Interval, 64)
	intervalArg := query.AddParameter(intervalNum)
	queryString := fmt.Sprintf("(%s / %s * %s)", f.FrequencyField, intervalArg, intervalArg)
	query.GroupBy(queryString)
	query.Select(fmt.Sprintf("%s as bucket", queryString))
	query.Select("COUNT(*) as frequency")

	return query
}

func (f *Frequency) AddQuery(query *Query) *Query {
	//TODO: Need to cast the frequency fields to a numeric value most likely.
	parameter := ""
	if f.GTE != nil {
		parameter = query.AddParameter(f.GTE)
		query.Where(fmt.Sprintf("%s >= %s", f.FrequencyField, parameter))
	}
	if f.GT != nil {
		parameter = query.AddParameter(f.GT)
		query.Where(fmt.Sprintf("%s > %s", f.FrequencyField, parameter))
	}
	if f.LTE != nil {
		parameter = query.AddParameter(f.LTE)
		query.Where(fmt.Sprintf("%s <= %s", f.FrequencyField, parameter))
	}
	if f.LT != nil {
		parameter = query.AddParameter(f.LT)
		query.Where(fmt.Sprintf("%s < %s", f.FrequencyField, parameter))
	}
	return query
}

func (f *Frequency) GetBuckets(rows *pgx.Rows) ([]*FrequencyResult, error) {
	//Need to build all the buckets over the window since empty buckets are needed.
	results := make(map[int]float64)
	//Parse the results. Build a map to fill in the buckets, and get the min/max.
	min, max := math.MaxInt32, math.MinInt32
	for rows.Next() {
		var bucket int
		var frequency int
		err := rows.Scan(&bucket, &frequency)
		if err != nil {
			return nil, fmt.Errorf("Error parsing top terms: %v", err)
		}
		results[bucket] = float64(frequency)
		if bucket < min {
			min = bucket
		}
		if bucket > max {
			max = bucket
		}
	}

	//Define the window limits.
	windowStart, windowEnd := 0, 0
	if f.GT != nil {
		windowStart = castFrequency(f.GT)
	} else if f.GTE != nil {
		windowStart = castFrequency(f.GTE)
	} else {
		windowStart = min
	}
	if f.LT != nil {
		windowEnd = castFrequency(f.LT)
	} else if f.LTE != nil {
		windowEnd = castFrequency(f.LTE)
	} else {
		windowEnd = max
	}

	//Create the buckets.
	intervalNum, err := strconv.ParseFloat(f.Interval, 64)
	if err != nil {
		return nil, err
	}

	//May be off by 1 as result of type conversion.
	numberOfBuckets := int64((windowEnd - windowStart)) / int64(intervalNum)
	buckets := make([]*FrequencyResult, numberOfBuckets)
	for i, _ := range buckets {
		//If value is not in the map, 0 will be returned as default value.
		bucket := i + windowStart
		frequency := results[bucket]
		buckets[i] = &FrequencyResult{
			Bucket: bucket,
			Value:  frequency,
		}
	}
	return buckets, nil
}

func castFrequency(val interface{}) int {
	num, isNum := val.(float64)
	if isNum {
		return int(num)
	}

	//TODO: Figure out which types are allowed, and what to do if bad data is received.
	return -1
}

func castTime(val interface{}) interface{} {
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

func castTimeToString(val interface{}) string {
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
