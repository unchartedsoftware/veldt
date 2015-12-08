package elastic

import (
	"encoding/json"
	"fmt"
)

// PropertyMeta represents the meta data for a single property.
type PropertyMeta struct {
	Type string `json:"type"`
	Extrema *Extrema `json:"extrema"`
}

func getSubMap( m map[string]interface{}, k string ) (map[string]interface{}, bool) {
	v, ok := m[k]
	if !ok {
		return nil, false
	}
	sub, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return sub, true
}

func getPropertyType( m map[string]interface{}, k string ) (string, bool) {
	v, ok := m[k]
	if !ok {
		return "", false
	}
	typ, ok := v.(string)
	if !ok {
		return "", false
	}
	return typ, true
}

func getMapping(endpoint string, index string) (map[string]interface{}, error) {
	// get client
	client, err := getClient(endpoint)
	if err != nil {
		return nil, err
	}
	// get mapping
	result, err := client.
		GetMapping().
		Index(index).
		Do()
	if err != nil {
		return nil, err
	}
	json, err := json.Marshal(result)
	fmt.Println(string(json))
	return result, nil
}

func getTypeAndExtents(endpoint string, index string, field string, prop map[string]interface{}) (*PropertyMeta, error) {
	typ, ok := getPropertyType(prop, "type")
	if !ok {
		return nil, fmt.Errorf("No 'type' attribute for property '%s'", field)
	}
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

func parsePropertiesRecursive(meta map[string]*PropertyMeta, endpoint string, index string, props map[string]interface{}, path string) error {
	for key, val := range props {
		subpath := key
		if path != "" {
			subpath = path + "." + key
		}
		subprops, ok := val.(map[string]interface{})
		if ok {
			props, ok := getSubMap(subprops, "properties")
			if ok {
				err := parsePropertiesRecursive(meta, endpoint, index, props, subpath)
				if err != nil {
					return err
				}
			} else {
				prop, err := getTypeAndExtents(endpoint, index, subpath, subprops)
				if err != nil {
					return err
				}
				meta[subpath] = prop
			}
		}
	}
	return nil
}

func parseProperties(meta map[string]*PropertyMeta, endpoint string, index string, props map[string]interface{}) error {
	return parsePropertiesRecursive(meta, endpoint, index, props, "")
}

func getBasePropertiesMap( mapping map[string]interface{}, index string ) (map[string]interface{}, bool) {
	path := []string{ index, "mappings", "datum", "properties" }
	m := mapping
	for _, key := range path {
		submap, ok := getSubMap(m, key)
		if !ok {
			return nil, false
		}
		m = submap
	}
	return m, true
}

// GetMeta returns the meta data for a given index.
func GetMeta(endpoint string, index string) (map[string]*PropertyMeta, error) {
	// get the raw mappings
	mapping, err := getMapping(endpoint, index)
	if err != nil {
		return nil, err
	}
	// get nested 'properties' attribute of mappings payload
	props, ok := getBasePropertiesMap(mapping, index)
	if !ok {
		return nil, fmt.Errorf("Unable to parse properties from mappings response for %s/%s", endpoint, index)
	}
	// create empty map
	meta := make(map[string]*PropertyMeta)
	// parse json mappings into the property map
	err = parseProperties(meta, endpoint, index, props)
	if err != nil {
		return nil, err
	}
	return meta, nil
}
