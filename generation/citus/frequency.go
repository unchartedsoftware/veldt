package citus

import (
	"fmt"
	"math"
	"strconv"

	"github.com/jackc/pgx"

	"github.com/unchartedsoftware/veldt/tile"
)

// Frequency represents a tiling generator that produces heatmaps.
type Frequency struct {
	tile.Frequency
}

// FrequencyResult represents a single frequency result bucket.
type FrequencyResult struct {
	Bucket int64
	Value  float64
}

// AddAggs adds the tiling aggregations to the provided query object.
func (f *Frequency) AddAggs(query *Query) *Query {
	//Bounds extension (empty buckets) will be done in the go code when parsing results
	//Not 100% sure if we need to subtract the min value from the frequency field to
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

// AddQuery adds the tiling query to the provided query object.
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

// GetBuckets returns the frequency buckets from the query results.
func (f *Frequency) GetBuckets(rows *pgx.Rows) ([]*FrequencyResult, error) {
	//Need to build all the buckets over the window since empty buckets are needed.
	results := make(map[int64]float64)
	//Parse the results. Build a map to fill in the buckets.
	for rows.Next() {
		var bucket int64
		var frequency int
		err := rows.Scan(&bucket, &frequency)
		if err != nil {
			return nil, fmt.Errorf("Error parsing top terms: %v", err)
		}
		results[bucket] = float64(frequency)
	}

	return f.CreateBuckets(results)
}

// CreateBuckets creates the frequency buckets, including the empty buckets as defined by the tile params.
func (f *Frequency) CreateBuckets(results map[int64]float64) ([]*FrequencyResult, error) {
	//Find the min & max buckets.
	min, max := int64(math.MaxInt64), int64(math.MinInt64)
	for k := range results {
		if k < min {
			min = k
		}
		if k > max {
			max = k
		}
	}

	//Define the window limits.
	windowStart, windowEnd := int64(0), int64(0)
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
	for i := range buckets {
		//If value is not in the map, 0 will be returned as default value.
		bucket := int64(float64(i)*intervalNum) + windowStart
		frequency := results[bucket]
		buckets[i] = &FrequencyResult{
			Bucket: bucket,
			Value:  frequency,
		}
	}
	return buckets, nil
}

func (f *Frequency) encodeResult(frequency []*FrequencyResult) []map[string]interface{} {
	buckets := make([]map[string]interface{}, len(frequency))
	for i, bucket := range frequency {
		buckets[i] = map[string]interface{}{
			"timestamp": bucket.Bucket,
			"count":     bucket.Value,
		}
	}
	return buckets
}

func castFrequency(val interface{}) int64 {
	numF, isNum := val.(float64)
	if isNum {
		return int64(numF)
	}
	numI, isNum := val.(int64)
	if isNum {
		return numI
	}
	// TODO: Figure out which types are allowed, and what to do if bad data is
	// received.
	return -1
}
