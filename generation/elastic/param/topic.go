package param

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// Topic represents params for extracting particular topics.
type Topic struct {
	Text   string
	Topics []string
}

// NewTopic instantiates and returns a new topic parameter object.
func NewTopic(tileReq *tile.Request) (*Topic, error) {
	params := tileReq.Params
	topics, ok := json.GetString(params, "topics")
	if !ok {
		return nil, fmt.Errorf("Topic parameters missing from tiling request %s", tileReq.String())
	}
	return &Topic{
		Text:   json.GetStringDefault(params, "text", "text"),
		Topics: strings.Split(topics, ","),
	}, nil
}

// GetHash returns a string hash of the parameter state.
func (p *Topic) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Text,
		strings.Join(p.Topics, ":"))
}

// GetTopicQuery returns an elastic query.
func (p *Topic) GetTopicQuery() *elastic.TermsQuery {
	ts := make([]interface{}, len(p.Topics))
	for i, str := range p.Topics {
		ts[i] = str
	}
	return elastic.NewTermsQuery(p.Text, ts...)
}

// GetTopicAggregations returns an elastic aggregation.
func (p *Topic) GetTopicAggregations() map[string]*elastic.FilterAggregation {
	aggs := make(map[string]*elastic.FilterAggregation, len(p.Topics))
	// add all filter aggregations
	for _, topic := range p.Topics {
		aggs[topic] = elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery(p.Text, topic))
	}
	return aggs
}
