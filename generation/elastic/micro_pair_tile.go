package elastic

import (
	"fmt"
	"encoding/json"

	"github.com/unchartedsoftware/prism"
	"github.com/unchartedsoftware/prism/binning"
	"github.com/unchartedsoftware/prism/tile"
)

// MicroPairTile represents an elasticsearch implementation of the microPair tile.
type MicroPairTile struct {
	MicroTile
	BivariatePair // No Parse()
	tile.MicroPair
}

// NewMicroTile instantiates and returns a new tile struct.
func NewMicroPairTile(host, port string) prism.TileCtor {
	return func() (prism.Tile, error) {
		m := &MicroPairTile{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

// Parse parses the provided JSON object and populates the tiles attributes.
func (m *MicroPairTile) Parse(params map[string]interface{}) error {
	// 1. tile.BivariatePair.Bivariate.Parse  ( xField, yField, bottom, etc, Resolution)
	err := m.BivariatePair.BivariatePair.Parse(params)
	if err != nil {
		return err
	}

	fmt.Printf("<><><> micropairtile: Parsing TopHits...%v\n", params)
	err = m.TopHits.Parse(params)
	if err != nil {
		return err
	}

	err = m.MicroTile.Micro.Parse(params)
	if err != nil {
		return err
	}

	fmt.Printf("><><>< m.BivariatePair %v \n", m)
	// parse includes
	m.TopHits.IncludeFields = m.MicroPair.ParseIncludes(
		m.TopHits.IncludeFields,
		m.BivariatePair.BivariatePair.Bivariate.XField,
		m.BivariatePair.BivariatePair.Bivariate.YField,
		m.BivariatePair.BivariatePair.X2Field,
		m.BivariatePair.BivariatePair.Y2Field)

	fmt.Printf("<><><> micropairtile: Parse completed. m.BivariatePair.X2Field=%s \n", m.BivariatePair.X2Field)
	return nil
}

func PrintJSON(x interface{}) {
	d, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Errorf("PrintJSON failed with error: %s\n", err)
		return
	}
	fmt.Printf("%s\n", d)
	return
}

// Create generates a tile from the provided URI, tile coordinate and query
// parameters.
func (m *MicroPairTile) Create(uri string, coord *binning.TileCoord, query prism.Query) ([]byte, error) {
	// get client
	client, err := NewClient(m.Host, m.Port)
	if err != nil {
		return nil, err
	}
	// create search service
	search := client.Search().
		Index(uri).
		Size(0)


	// create root query
	q, err := m.CreateQuery(query)
	if err != nil {
		return nil, err
	}
	// add tiling query
	q.Must(m.BivariatePair.GetQuery(coord))
	// set the query
	search.Query(q)

	// get aggs
	aggs := m.TopHits.GetAggs()
	// set the aggregation
	search.Aggregation("top-hits", aggs["top-hits"])

	//source, _ := q.Source()
	fmt.Printf("<><><> micropairtile: sending query: \n")
	//PrintJSON(source)

	// send query
	res, err := search.Do()
	if err != nil {
		fmt.Printf("<><><> micropairtile: query error %s\n", err)
		return nil, err
	}

	// get top hits
	hits, err := m.TopHits.GetTopHits(&res.Aggregations)
	if err != nil {
		return nil, err
	}

	// convert to point array
	points := make([]float32, len(hits)*4)
	fmt.Printf("<><><> micropairtile: hits %d\n", len(hits))
	for i, hit := range hits {
		// get hit x/y in tile coords
		x, y, ok := m.BivariatePair.BivariatePair.GetXY(hit)
		if !ok {
			fmt.Printf("couldn't GetXY\n");
			continue
		}
		x2, y2, ok := m.BivariatePair.BivariatePair.GetX2Y2(hit)
		if !ok {
			fmt.Printf("couldn't GetX2Y2\n");
			continue
		}
		// add to point array
		points[i*4] = x
		points[i*4+1] = y
		points[i*4+2] = x2
		points[i*4+3] = y2
	}
	// encode and return results
	return m.MicroPair.Encode(hits, points)
}
