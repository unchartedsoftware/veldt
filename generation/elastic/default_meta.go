package elastic

import (
	"fmt"

	"gopkg.in/olivere/elastic.v3"

	"github.com/unchartedsoftware/veldt"
	"github.com/unchartedsoftware/veldt/binning"
	"github.com/unchartedsoftware/veldt/util/json"
)

// DefaultMeta represents a meta data generator that produces default
// metadata with property types and extrema.
type DefaultMeta struct {
	Elastic
}

// NewDefaultMeta instantiates and returns a pointer to a new generator.
func NewDefaultMeta(host string, port string) veldt.MetaCtor {
	return func() (veldt.Meta, error) {
		m := &DefaultMeta{}
		m.Host = host
		m.Port = port
		return m, nil
	}
}

// Parse parses the provided JSON object and populates the structs attributes.
func (m *DefaultMeta) Parse(params map[string]interface{}) error {
	return nil
}

// Create generates metadata from the provided URI.
func (m *DefaultMeta) Create(uri string) ([]byte, error) {
	// get the raw mappings
	service, err := m.CreateMappingService(uri)
	if err != nil {
		return nil, err
	}
	// get the raw mappings
	mapping, err := service.Do()
	if err != nil {
		return nil, err
	}
	// get nested 'properties' attribute of mappings payload
	// NOTE: If running a `mapping` query on an aliased index, the mapping
	// response will be nested under the original index name. Since we are only
	// getting the mapping of a single index at a time, we can simply get the
	// 'first' and only node.
	_, index, ok := json.GetRandomChild(mapping)
	if !ok {
		return nil, fmt.Errorf("Unable to retrieve the mappings response for %s",
			uri)
	}
	// get mappings node
	mappings, ok := json.GetChildMap(index, "mappings")
	if !ok {
		return nil, fmt.Errorf("unable to parse `mappings` from mappings response for %s",
			uri)
	}
	// for each type, parse the mapping
	meta := make(map[string]interface{})
	for key, typ := range mappings {
		typeMeta, err := m.parseType(uri, typ)
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

func (m *DefaultMeta) getExtrema(uri string, field string) (*binning.Extrema, error) {
	// search
	search, err := m.CreateSearchService(uri)
	if err != nil {
		return nil, err
	}
	result, err := search.Aggregation("min",
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
		return nil, fmt.Errorf("min '%s' aggregation was not found in response for %s", field, uri)
	}
	max, ok := result.Aggregations.Max("max")
	if !ok {
		return nil, fmt.Errorf("max '%s' aggregation was not found in response for %s", field, uri)
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

func (m *DefaultMeta) getPropertyMeta(uri string, field string, typ string) (*PropertyMeta, error) {
	prop := &PropertyMeta{
		Type: typ,
	}
	// if field is ordinal, get the extrema
	if isOrdinal(typ) {
		extrema, err := m.getExtrema(uri, field)
		if err != nil {
			return nil, err
		}
		prop.Extrema = extrema
	}
	return prop, nil
}

func (m *DefaultMeta) parsePropertiesRecursive(meta map[string]PropertyMeta, uri string, p map[string]interface{}, path string) error {
	children, ok := json.GetChildMap(p)
	if !ok {
		return nil
	}

	for key, props := range children {
		subpath := key
		if path != "" {
			subpath = path + "." + key
		}
		subprops, ok := json.GetChild(props, "properties")
		if ok {
			// recurse further
			err := m.parsePropertiesRecursive(meta, uri, subprops, subpath)
			if err != nil {
				return err
			}
		} else {
			typ, ok := json.GetString(props, "type")
			// we don't support nested types
			if ok && typ != "nested" {

				prop, err := m.getPropertyMeta(uri, subpath, typ)
				if err != nil {
					return err
				}
				meta[subpath] = *prop

				// Parse out multi-field mapping
				fields, hasFields := json.GetChild(props, "fields")
				if hasFields {
					for fieldName := range fields {
						multiFieldPath := subpath + "." + fieldName
						prop, err = m.getPropertyMeta(uri, multiFieldPath, typ)
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

func (m *DefaultMeta) parseProperties(uri string, props map[string]interface{}) (map[string]PropertyMeta, error) {
	// create empty map
	meta := make(map[string]PropertyMeta)
	// parse recursively, appending to the map
	err := m.parsePropertiesRecursive(meta, uri, props, "")
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (m *DefaultMeta) parseType(uri string, typ map[string]interface{}) (map[string]PropertyMeta, error) {
	props, ok := json.GetChild(typ, "properties")
	if !ok {
		return nil, fmt.Errorf("Unable to parse `properties` from mappings response for type `%s` for %s",
			typ,
			uri)
	}
	// parse json mappings into the property map
	return m.parseProperties(uri, props)
}
