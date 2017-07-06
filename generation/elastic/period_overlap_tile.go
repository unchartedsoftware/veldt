package elastic

import (
	"encoding/binary"
	"fmt"
	"time"
	"math"

	//ejson "encoding/json"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/geometry"
	"github.com/unchartedsoftware/veldt/util/json"
)

const (
	millisInHour        = 60 * 60 * 1000
	millisInDay         = 24 * millisInHour
	millisInWeek        = 7  * millisInDay
	millisInMonth       = 31 * millisInDay
	millisInYear        = 365 * millisInDay
)

// PeriodOverlapTile represents an elasticsearch implementation of the period overlap tile.
type PeriodOverlapTile struct {
	Elastic
	Bivariate
	globalBounds *geometry.Bounds
	period string
}

// NewPeriodOverlapTile instantiates and returns a new tile struct.
func NewPeriodOverlapTile(host, port string) veldt.TileCtor {
	return func() (veldt.Tile, error) {
		fmt.Println("PeriodOverlapTile Tile const!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		h := &PeriodOverlapTile{}
		h.Host = host
		h.Port = port
		return h, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (p *PeriodOverlapTile) Parse(params map[string]interface{}) error {
	fmt.Println("tile params", params);
	//Need global bounds for the query creation
	p.globalBounds = &geometry.Bounds{}
	p.globalBounds.Parse(params)
	period, ok := json.GetString(params, "xField")
	if !ok {
		return fmt.Errorf("`period` parameter missing from tile")
	}
	p.period = period
	fmt.Println("period", p.period);

	return p.Bivariate.Parse(params)
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (p *PeriodOverlapTile) Create(uri string, coord *binning.TileCoord, query veldt.Query) ([]byte, error) {
	fmt.Println("PeriodOverlapTile Tile create!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// create search service
	search, err := p.CreateSearchService(uri)
	if err != nil {
		return nil, err
	}

	// create root query
	q, err := p.CreateQuery(query)
	if err != nil {
		return nil, err
	}

	//Param testing
	//bounds := p.Bivariate.TileBounds(coord)
	//fmt.Println("tile min max", bounds.MinX(), bounds.MaxX(), bounds.MinY(), bounds.MaxY());
	//fmt.Println("Parse", p.globalBounds.MinX(), p.globalBounds.MaxX());





	// add tiling query
	//q.Must(p.Bivariate.GetQuery(coord))	//Heatmap way
	q.Must(p.GetQuery(coord))		//New way
	// set the query
	search.Query(q)

	//qs1, _ := q.Source()
	//marshalledQ1, _ := ejson.MarshalIndent(qs1, "", "  ")
	//fmt.Println("Period overlap Query1!!!:", string(marshalledQ1))


	// get aggs
	aggs := p.Bivariate.GetAggs(coord)
	// set the aggregation
	search.Aggregation("x", aggs["x"])

	// send query
	res, err := search.Do()
	if err != nil {
		return nil, err
	}

	fmt.Println("Create tile, about to get bins")
	// get bins
	//bins, err := p.Bivariate.GetBins(coord, &res.Aggregations)
	bins, err := p.getBins(coord, &res.Aggregations)
	if err != nil {
		return nil, err
	}

	fmt.Println("Create tile, got the bins, creating byte array")

	// convert to byte array
	bits := make([]byte, len(bins)*4)
	for i, bin := range bins {
		if bin != nil {
			binary.LittleEndian.PutUint32(
				bits[i*4:i*4+4],
				uint32(bin.DocCount))
		}
	}
	return bits, nil
}

func (p *PeriodOverlapTile) getPeriodMillis() int64 {
	if p.period == "Hourly"  {
		return millisInHour
	} else if p.period == "Daily" {
		return millisInDay
	} else if p.period == "Weekly" {
		return millisInWeek
	} else if p.period == "Monthly" {
		return millisInMonth
	} else if p.period == "Yearly" {
		return millisInYear
	} else {
		return int64(0)
	}
}

func (p *PeriodOverlapTile) isTileRangeGtePeriod(tileMinX int64, tileMaxX int64) bool {
	tileMinXTime := time.Unix(0, tileMinX*int64(time.Millisecond)).UTC()
	tileMaxXTime := time.Unix(0, tileMaxX*int64(time.Millisecond)).UTC()

	tileRangeDuration := int64(tileMaxXTime.Sub(tileMinXTime).Seconds() * 1000)
	periodMillis := p.getPeriodMillis()
	return tileRangeDuration >= periodMillis
}

func (p *PeriodOverlapTile) getInitRangeStartX(tileMinX int64) int64 {
	millisInPeriod := p.getPeriodMillis()
	rangeMult := math.Floor(float64((tileMinX - int64(p.globalBounds.MinX())) / millisInPeriod))
	initRangeStartX := tileMinX - (int64(rangeMult) * millisInPeriod)
	return initRangeStartX
}

func (p *PeriodOverlapTile) getTileRangeOffset(tileMinX int64) int64 {
	offset := math.Mod(float64(tileMinX - int64(p.globalBounds.MinX())), millisInMonth)

	return int64(math.Floor(offset))
}

func (p *PeriodOverlapTile) multiCalenderMonthQuery(offset int64, tileMinX int64, tileMaxX int64) bool {
	janFirst := time.Date(2015, time.Month(1), 1, 0, 0, 0, 0, time.UTC)	//Year doesn't matter
	rangeStartTimeMillis := janFirst.Unix()*1000 + offset
	rangeEndTimeMillis := rangeStartTimeMillis + tileMaxX - tileMinX

	return time.Unix(0, rangeStartTimeMillis * int64(time.Millisecond)).UTC().Month() != time.Unix(0, rangeEndTimeMillis * int64(time.Millisecond)).UTC().Month()
}

func (p *PeriodOverlapTile) adjustEndDateMonth(startDateMillis int64, endDateMillis int64, multiCalendarMonth bool) (int64) {
	startDateTime := time.Unix(0, startDateMillis * int64(time.Millisecond)).UTC()
	endDateTime := time.Unix(0, endDateMillis * int64(time.Millisecond)).UTC()
	endDateMonth := startDateTime.Month()
	if multiCalendarMonth && startDateTime.Month() != endDateTime.Month() {
		endMonthInt := int(endDateMonth) + 1
		if endMonthInt > int(time.December) {
			endMonthInt = 1
		}
		endDateMonth = time.Month(endMonthInt)
	}

	for endDateTime.Month() != endDateMonth {
		endDateTime = endDateTime.Add(time.Hour * 24 * -1)
	}

	return endDateTime.Unix()*1000
}

func (p *PeriodOverlapTile) clampDates(startDate int64, endDate int64) (int64, int64) {
	return int64(math.Max(float64(startDate), p.globalBounds.MinX())), int64(math.Min(float64(endDate), p.globalBounds.MaxX()));
}

//TODO Do I take into account data sets that cross a calendar year?
func (p *PeriodOverlapTile) getQueryMonthly(tileMinX int64, tileMaxX int64) *elastic.BoolQuery {
	dateRangeQuery := elastic.NewBoolQuery()

	tileRangeOffset := p.getTileRangeOffset(tileMinX)

	mapMinXTime := time.Unix(0, int64(p.globalBounds.MinX()) * int64(time.Millisecond)).UTC()
	queryRangeBaseTime := time.Date(mapMinXTime.Year(), time.Month(int(mapMinXTime.Month()) - 1), 1, 0, 0, 0, 0, time.UTC)
	multiCalendarMonthQuery := p.multiCalenderMonthQuery(tileRangeOffset, tileMinX, tileMaxX)

	currRangeStartTimeUnix := queryRangeBaseTime.Unix()*1000 + tileRangeOffset
	currRangeEndTimeUnix := p.adjustEndDateMonth(currRangeStartTimeUnix, currRangeStartTimeUnix + (tileMaxX - tileMinX), multiCalendarMonthQuery)

	currRangeStartTimeUnix, currRangeEndTimeUnix = p.clampDates(currRangeStartTimeUnix, currRangeEndTimeUnix);

	//TODO Consider data set that is greater than 1 year
	for i := 1; currRangeStartTimeUnix < int64(p.globalBounds.MaxX()); i++ {
		if currRangeEndTimeUnix > int64(p.globalBounds.MinX()) {
			rangeQuery := elastic.NewRangeQuery(p.Bivariate.XField).
				Gte(currRangeStartTimeUnix).
				Lt(currRangeEndTimeUnix)
			dateRangeQuery = dateRangeQuery.Should(rangeQuery)
		}
		queryRangeBaseTime = queryRangeBaseTime.AddDate(0,1,0)
		currRangeStartTimeUnix = queryRangeBaseTime.Unix()*1000 + tileRangeOffset
		currRangeEndTimeUnix = p.adjustEndDateMonth(currRangeStartTimeUnix, currRangeStartTimeUnix + (tileMaxX - tileMinX), multiCalendarMonthQuery)
		currRangeStartTimeUnix, currRangeEndTimeUnix = p.clampDates(currRangeStartTimeUnix, currRangeEndTimeUnix);
	}

	return dateRangeQuery
}


func (p *PeriodOverlapTile) getQuery(coord *binning.TileCoord) elastic.Query {
	bounds := p.Bivariate.TileBounds(coord)
	tileMinX := int64(bounds.MinX())
	tileMaxX := int64(bounds.MaxX())
	dateRangeQuery := elastic.NewBoolQuery()

	if p.isTileRangeGtePeriod(tileMinX, tileMaxX) {
		dateRangeQuery = elastic.NewBoolQuery().
			Must(elastic.NewRangeQuery(p.Bivariate.XField).
			Gte(tileMinX).
			Lt(tileMaxX))
	} else {
		if p.period == "Monthly" {
			dateRangeQuery = p.getQueryMonthly(tileMinX, tileMaxX)
		} else {
			tileRange := tileMaxX - tileMinX
			currRangeStartX := p.getInitRangeStartX(tileMinX)
			currRangeEndX := currRangeStartX + tileRange

			for i := 1; currRangeStartX < int64(p.globalBounds.MaxX()); i++ {
				rangeQuery := elastic.NewRangeQuery(p.Bivariate.XField).
					Gte(currRangeStartX).
					Lt(currRangeEndX)
				dateRangeQuery = dateRangeQuery.Should(rangeQuery)
				currRangeStartX += p.getPeriodMillis()
				currRangeEndX += p.getPeriodMillis()
				currRangeStartX, currRangeEndX = p.clampDates(currRangeStartX, currRangeEndX);
			}
		}
	}


	query := elastic.NewBoolQuery()
	query.Must(dateRangeQuery)
	query.Must(elastic.NewRangeQuery(p.Bivariate.YField).
		Gte(int64(bounds.MinY())).
		Lt(int64(bounds.MaxY())))
	return query
}

func (b *PeriodOverlapTile) clampBin(bin int64) int {
	if bin > int64(b.Resolution)-1 {
		return b.Resolution - 1
	}
	if bin < 0 {
		return 0
	}
	return int(bin)
}

func (p *PeriodOverlapTile) getXBinForBounds(coord *binning.TileCoord, x int64, left int64, right int64) int {
	fx := float64(x)
	var bin int64
	if left > right {
		bin = int64(float64(p.Resolution) - ((fx - float64(right)) / p.BinSizeX(coord)))
	} else {
		bin = int64((fx - float64(left)) / p.BinSizeX(coord))
	}
	return p.clampBin(bin)
}

//TODO for monthly pick a Month that has 31 days, use that as the base month to convert all the bucket dates
func (p *PeriodOverlapTile) getBins(coord *binning.TileCoord, aggs *elastic.Aggregations) ([]*elastic.AggregationBucketHistogramItem, error) {
	fmt.Println("Getting bins")
	bounds := p.Bivariate.TileBounds(coord)
	tileMinX := int64(bounds.MinX())
	tileMaxX := int64(bounds.MaxX())

	//Get tile range
	tileRange := tileMaxX - tileMinX
	tileMinXTime := time.Unix(0, tileMinX*int64(time.Millisecond)).UTC()

	if p.period == "Monthly" {
		tileRangeOffset := p.getTileRangeOffset(tileMinX)
		tileMinX = time.Date(tileMinXTime.Year(), time.July, 1, 0, 0, 0, 0, time.UTC).Unix() * 1000
		tileMinX += tileRangeOffset
		tileMaxX = tileMinX + tileRange
	}

	// parse aggregations
	xAgg, ok := aggs.Histogram("x")
	if !ok {
		return nil, fmt.Errorf("histogram aggregation `x` was not found")
	}


	// allocate bins buffer
	bins := make([]*elastic.AggregationBucketHistogramItem, p.Bivariate.Resolution*p.Bivariate.Resolution)

	fmt.Println("about to fill bins")

	// fill bins buffer
	for _, xBucket := range xAgg.Buckets {

		xBucketKey := xBucket.Key

		fmt.Println("xBucketKey", xBucketKey)

		/*
		 * Change the bucket date so it falls in our tile date range according to the period. Eg: if the tile date range is
		 * Fri Sep 11 01:59:55 UTC 2015 to Sun Sep 13 23:00:19 UTC 2015 and the bucket date is Sat Jul 4 10:23:55 UTC 2015
		 * with a weekly period then the we want the bucket date to be Sat Sep 12 10:23:55 UTC 2015
		 */

		//Change the x key so it's in the proper range
		var xKeys []int64

		//Get a date in tile range that matches the same day of the week as our bucket date
		matchTime := tileMinXTime
		xBucketTime := time.Unix(0, xBucketKey*int64(time.Millisecond)).UTC()

		if p.period == "Hourly" {
			matchTime = time.Date(matchTime.Year(), matchTime.Month(), matchTime.Day(), matchTime.Hour(), 0, 0, 0, time.UTC)
			bucketTimeOfHour := xBucketKey - time.Date(xBucketTime.Year(), xBucketTime.Month(), xBucketTime.Day(), xBucketTime.Hour(), 0, 0, 0, time.UTC).Unix()*1000
			//Add the bucket time of day to our stripped match date.
			xBucketKey = matchTime.Unix()*1000 + bucketTimeOfHour
		} else if p.period == "Daily" {
			matchTime = time.Date(matchTime.Year(), matchTime.Month(), matchTime.Day(), 0, 0, 0, 0, time.UTC)
			bucketTimeOfDay := xBucketKey - time.Date(xBucketTime.Year(), xBucketTime.Month(), xBucketTime.Day(), 0, 0, 0, 0, time.UTC).Unix()*1000
			//Add the bucket time of day to our stripped match date.
			xBucketKey = matchTime.Unix()*1000 + bucketTimeOfDay
		} else if p.period == "Weekly" {
			//Add a day to the tile min date until with have a match
			for matchTime.Weekday() != xBucketTime.Weekday() {
				matchTime = matchTime.Add(time.Hour * 24)
			}
			//Strip of the hours, mins, secs etc from out match date
			matchTime = time.Date(matchTime.Year(), matchTime.Month(), matchTime.Day(), 0, 0, 0, 0, time.UTC)
			//Get the hours, mins, secs etc from our bucket date
			bucketTimeOfDay := xBucketKey - time.Date(xBucketTime.Year(), xBucketTime.Month(), xBucketTime.Day(), 0, 0, 0, 0, time.UTC).Unix()*1000
			//Add the bucket time of day to our stripped match date.
			xBucketKey = matchTime.Unix()*1000 + bucketTimeOfDay
		} else if p.period == "Monthly" {
			//Convert bucketKey date month to Jan, since it has 31 days
			xBucketKey = time.Date(
				xBucketTime.Year(),
				time.July,
				xBucketTime.Day(),
				xBucketTime.Hour(),
				xBucketTime.Minute(),
				xBucketTime.Second(),
				xBucketTime.Nanosecond(),
				time.UTC).Unix()*1000

			if xBucketKey < tileMinX {
				xBucketKey = time.Date(
					xBucketTime.Year(),
					time.August,
					xBucketTime.Day(),
					xBucketTime.Hour(),
					xBucketTime.Minute(),
					xBucketTime.Second(),
					xBucketTime.Nanosecond(),
					time.UTC).Unix()*1000
			}

			if xBucketKey < tileMinX || xBucketKey > tileMaxX {
				xBucketKey = -1
			}
		}

		if xBucketKey >= 0 {
			xKeys = append(xKeys, xBucketKey)

			//Add additional dates if the tile range is bigger than the period
			millisInPeriod := p.getPeriodMillis()
			for i := xBucketKey + millisInPeriod; i < tileMaxX; i += millisInPeriod {
				xKeys = append(xKeys, i)
			}
		}

		fmt.Println("Got xKeys, getting yKeys")

		for _, xKey := range xKeys {
			var xBin int
			if p.period == "Monthly" {
				//Get the xBin for our converted tile range
				xBin = p.getXBinForBounds(coord, xKey, tileMinX, tileMaxX)
			} else {
				xBin = p.Bivariate.GetXBin(coord, float64(xKey))  //binning.GetXBin(xKey)
			}

			yAgg, ok := xBucket.Histogram("y")

			if !ok {
				return nil, fmt.Errorf("histogram aggregation `y` was not found")
			}
			for _, yBucket := range yAgg.Buckets {
				y := yBucket.Key
				yBin := p.Bivariate.GetYBin(coord, float64(y))
				index := xBin + p.Bivariate.Resolution * yBin

				// encode count
				//bins[index] += yBucket
				//TODO shouldn't we be aggregatig doc count like we did in prism?
				bins[index] = yBucket
			}
		}

		fmt.Println("Filled bin for xKey")
	}

	return bins, nil
}
