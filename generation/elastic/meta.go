package elastic

import (
	"encoding/json"
	"fmt"

	"github.com/unchartedsoftware/prism/generation/meta"
	jsonutil "github.com/unchartedsoftware/prism/util/json"
)

// PropertyMeta represents the meta data for a single property.
type PropertyMeta struct {
	Type    string   `json:"type"`
	Extrema *Extrema `json:"extrema,omitempty"`
}

func getPropertyMeta(endpoint string, index string, field string, typ string) (*PropertyMeta, error) {
	p := PropertyMeta{
		Type: typ,
	}
	// if field is 'numeric', get the extrema
	if typ == "long" || typ == "double" || typ == "date" {
		extrema, err := GetExtrema(endpoint, index, field)
		if err != nil {
			return nil, err
		}
		p.Extrema = extrema
	}
	return &p, nil
}

func parsePropertiesRecursive(meta map[string]PropertyMeta, endpoint string, index string, p map[string]interface{}, path string) error {
	for key, props := range jsonutil.GetChildren(p) {
		subpath := key
		if path != "" {
			subpath = path + "." + key
		}
		subprops, hasProps := jsonutil.GetChild(props, "properties")
		if hasProps {
			// recurse further
			err := parsePropertiesRecursive(meta, endpoint, index, subprops, subpath)
			if err != nil {
				return err
			}
		} else {
			typ, hasType := jsonutil.GetString(props, "type")
			// we don't support nested types
			if hasType && typ != "nested" {
				prop, err := getPropertyMeta(endpoint, index, subpath, typ)
				if err != nil {
					return err
				}
				meta[subpath] = *prop
			}
		}
	}
	return nil
}

func parseProperties(endpoint string, index string, props map[string]interface{}) (map[string]PropertyMeta, error) {
	// create empty map
	meta := make(map[string]PropertyMeta)
	err := parsePropertiesRecursive(meta, endpoint, index, props, "")
	if err != nil {
		return nil, err
	}
	return meta, nil
}

// GetMeta returns the meta data for a given index.
func GetMeta(metaReq *meta.MetaRequest) ([]byte, error) {
	// get the raw mappings
	mapping, err := GetMapping(metaReq.Endpoint, metaReq.Index)
	if err != nil {
		return nil, err
	}
	// get nested 'properties' attribute of mappings payload
	props, ok := jsonutil.GetChild(mapping, metaReq.Index, "mappings", "datum", "properties")
	if !ok {
		return nil, fmt.Errorf("Unable to parse properties from mappings response for %s/%s",
			metaReq.Endpoint,
			metaReq.Index)
	}
	// parse json mappings into the property map
	meta, err := parseProperties(metaReq.Endpoint, metaReq.Index, props)
	if err != nil {
		return nil, err
	}
	return json.Marshal(meta)
}
