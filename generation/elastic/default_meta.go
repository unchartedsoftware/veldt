package elastic

import (
	"encoding/json"
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	jsonutil "github.com/unchartedsoftware/veldt/util/json"
)

// DefaultMeta represents a meta data generator that produces default
// metadata with property types and extrema.
type DefaultMeta struct {
	Host string
	Port string
}

// NewDefaultMeta instantiates and returns a pointer to a new generator.
func NewDefaultMeta(host string, port string) veldt.MetaCtor {
	return func() (veldt.Meta, error) {
		return &DefaultMeta{
			Host: host,
			Port: port,
		}, nil
	}
}

// Parse parses the provided JSON object and populates the structs attributes.
func (g *DefaultMeta) Parse(params map[string]interface{}) error {
	return nil
}

// Create generates metadata from the provided URI.
func (g *DefaultMeta) Create(uri string) ([]byte, error) {
	client, err := NewClient(g.Host, g.Port)
	if err != nil {
		return nil, err
	}
	// get the raw mappings
	mapping, err := client.GetMapping().Index(uri).Do()
	if err != nil {
		return nil, err
	}
	// get nested 'properties' attribute of mappings payload
	// NOTE: If running a `mapping` query on an aliased index, the mapping
	// response will be nested under the original index name. Since we are only
	// getting the mapping of a single index at a time, we can simply get the
	// 'first' and only node.
	index, ok := jsonutil.GetRandomChild(mapping)
	if !ok {
		return nil, fmt.Errorf("Unable to retrieve the mappings response for %s",
			uri)
	}
	// get mappings node
	mappings, ok := jsonutil.GetChildMap(index, "mappings")
	if !ok {
		return nil, fmt.Errorf("unable to parse `mappings` from mappings response for %s",
			uri)
	}
	// for each type, parse the mapping
	meta := make(map[string]interface{})
	for key, typ := range mappings {
		typeMeta, err := parseType(client, uri, typ)
		if err != nil {
			return nil, err
		}
		meta[key] = typeMeta
	}
	// return
	return json.Marshal(meta)
}

// PropertyMeta represents the meta data for a single property.
type PropertyMeta struct {
	Type    string           `json:"type"`
	Extrema *binning.Extrema `json:"extrema,omitempty"`
}

func isOrdinal(typ string) bool {
	return typ == "long" ||
		typ == "integer" ||
		typ == "short" ||
		typ == "byte" ||
		typ == "double" ||
		typ == "float" ||
		typ == "date"
}

func getExtrema(client *elastic.Client, index string, field string) (*binning.Extrema, error) {
	// query
	result, err := client.
		Search(index).
		Size(0).
		Aggregation("min",
			elastic.NewMinAggregation().
				Field(field)).
		Aggregation("max",
			elastic.NewMaxAggregation().
				Field(field)).
		Do()
	if err != nil {
		return nil, err
	}
	// parse aggregations
	min, ok := result.Aggregations.Min("min")
	if !ok {
		return nil, fmt.Errorf("min '%s' aggregation was not found in response for %s", field, index)
	}
	max, ok := result.Aggregations.Max("max")
	if !ok {
		return nil, fmt.Errorf("max '%s' aggregation was not found in response for %s", field, index)
	}
	// if the mapping exists, but no documents have the attribute, the min / max
	// are null
	if min.Value == nil || max.Value == nil {
		return nil, nil
	}
	return &binning.Extrema{
		Min: *min.Value,
		Max: *max.Value,
	}, nil
}

func getPropertyMeta(client *elastic.Client, index string, field string, typ string) (*PropertyMeta, error) {
	p := PropertyMeta{
		Type: typ,
	}
	// if field is 'ordinal', get the extrema
	if isOrdinal(typ) {
		extrema, err := getExtrema(client, index, field)
		if err != nil {
			return nil, err
		}
		p.Extrema = extrema
	}
	return &p, nil
}

func parsePropertiesRecursive(meta map[string]PropertyMeta, client *elastic.Client, index string, p map[string]interface{}, path string) error {
	children, ok := jsonutil.GetChildMap(p)
	if !ok {
		return nil
	}

	for key, props := range children {
		subpath := key
		if path != "" {
			subpath = path + "." + key
		}
		subprops, hasProps := jsonutil.GetChild(props, "properties")
		if hasProps {
			// recurse further
			err := parsePropertiesRecursive(meta, client, index, subprops, subpath)
			if err != nil {
				return err
			}
		} else {
			typ, hasType := jsonutil.GetString(props, "type")
			// we don't support nested types
			if hasType && typ != "nested" {

				prop, err := getPropertyMeta(client, index, subpath, typ)
				if err != nil {
					return err
				}
				meta[subpath] = *prop

				// Parse out multi-field mapping
				fields, hasFields := jsonutil.GetChild(props, "fields")
				if hasFields {
					for fieldName := range fields {
						multiFieldPath := subpath + "." + fieldName
						prop, err = getPropertyMeta(client, index, multiFieldPath, typ)
						if err != nil {
							return err
						}
						meta[multiFieldPath] = *prop
					}
				}
			}
		}
	}

	return nil
}

func parseProperties(client *elastic.Client, index string, props map[string]interface{}) (map[string]PropertyMeta, error) {
	// create empty map
	meta := make(map[string]PropertyMeta)
	err := parsePropertiesRecursive(meta, client, index, props, "")
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func parseType(client *elastic.Client, index string, typ map[string]interface{}) (map[string]PropertyMeta, error) {
	props, ok := jsonutil.GetChild(typ, "properties")
	if !ok {
		return nil, fmt.Errorf("Unable to parse `properties` from mappings response for type `%s` for %s",
			typ,
			index)
	}
	// parse json mappings into the property map
	return parseProperties(client, index, props)
}
