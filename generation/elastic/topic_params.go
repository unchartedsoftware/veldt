package elastic

import (
	"fmt"
	"strings"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/prism/generation/tile"
	"github.com/unchartedsoftware/prism/util/json"
)

// TopicParams represents params for extracting particular topics.
type TopicParams struct {
	Text	string
	Topics  []string
}

// NewTopicParams parses the params map returns a pointer to the param struct.
func NewTopicParams(tileReq *tile.Request) *TopicParams {
	params := tileReq.Params
	topicsStr := json.GetStringDefault(params, "topics", "")
	topics := strings.Split(topicsStr, ",")
	if len(topics) == 0 {
		return nil
	}
	return &TopicParams{
		Text: json.GetStringDefault(params, "text", "text"),
		Topics: topics,
	}
}

// GetHash returns a string hash of the parameter state.
func (p *TopicParams) GetHash() string {
	return fmt.Sprintf("%s:%s",
		p.Text,
		strings.Join(p.Topics, ":"))
}

// GetTopicQuery returns an elastic query.
func (p *TopicParams) GetTopicQuery() *elastic.TermsQuery {
	ts := make([]interface{}, len(p.Topics))
	for _, str := range p.Topics {
		ts = append(ts, str)
	}
	return elastic.NewTermsQuery(p.Text, ts...)
}

// GetTopicAggregations returns an elastic aggregation.
func (p *TopicParams) GetTopicAggregations() map[string]*elastic.FilterAggregation {
	aggs := make(map[string]*elastic.FilterAggregation, len(p.Topics))
	// add all filter aggregations
	for _, topic := range p.Topics {
		aggs[topic] = elastic.NewFilterAggregation().
			Filter(elastic.NewTermQuery(p.Text, topic))
	}
	return aggs
}
