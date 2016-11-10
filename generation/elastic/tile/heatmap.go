package tile

import (
	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/tile"
)

const (
	xAggName = "x"
	yAggName = "y"
)

// Heatmap represents a tiling generator that produces heatmaps.
type Heatmap struct {
	generation.Heatmap
}

func (h *Heatmap) ApplyQuery(query *elastic.BoolQuery) error {
	h.Bivariate.ApplyQuery(query))
	return nil
}

func (h *Heatmap) ApplyAgg(search *elastic.SearchService) error {
	h.Bivariate.ApplyAgg(query))
	return nil
}

func (h *Heatmap) ParseRes(res *elastic.SearchResult) ([]byte, error) {
	bins, err := h.Bivariate.GetBins(query))
	if err != nil {
		return nil, err
	}
	buffer := make([]float64, len(bins))
	for i, bin := range bins {
		buffer[i] = bin.DocCount
	}
	return g.Float64ToBytes(buffer), nil
}
