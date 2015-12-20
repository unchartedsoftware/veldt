package param

import (
	"errors"
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
	topicsStr := json.GetStringDefault(params, "topics", "")
	topics := strings.Split(topicsStr, ",")
	if len(topics) == 0 {
		return nil, errors.New("Topic parameters missing from tiling request")
	}
	return &Topic{
		Text:   json.GetStringDefault(params, "text", "text"),
		Topics: topics,
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
	for _, str := range p.Topics {
		ts = append(ts, str)
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
