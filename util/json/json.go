package json

import (
	"strconv"
)

// Node represents a single json node as a map[string]interface{}
type Node map[string]interface{}

// Set sets the value under a given path, creating intermediate nodes along the
// way if they do not exist.
func Set(json Node, v interface{}, path ...string) {
	child := json
	last := len(path) - 1
	for _, key := range path[:last] {
		v, ok := child[key]
		if !ok {
			v = make(map[string]interface{})
		}
		c, ok := v.(map[string]interface{})
		if !ok {
			c = make(map[string]interface{})
		}
		child[key] = c
		child = c
	}
	child[path[last]] = v
}

// GetChild returns the child under the given path.
func GetChild(json Node, path ...string) (Node, bool) {
	child := json
	for _, key := range path {
		v, ok := child[key]
		if !ok {
			return nil, false
		}
		c, ok := v.(map[string]interface{})
		if !ok {
			return nil, false
		}
		child = c
	}
	return child, true
}

// GetChildren returns a map of all child nodes under a given path.
func GetChildren(json Node, path ...string) (map[string]Node, bool) {
	sub, ok := GetChild(json, path...)
	if !ok {
		return nil, false
	}
	children := make(map[string]Node, len(sub))
	for k, v := range sub {
		c, ok := v.(map[string]interface{})
		if ok {
			children[k] = c
		}
	}
	return children, true
}

// GetString returns a string property under the given key.
func GetString(json Node, key string) (string, bool) {
	v, ok := json[key]
	if !ok {
		return "", false
	}
	val, ok := v.(string)
	if !ok {
		return "", false
	}
	return val, true
}

// GetStringDefault returns a string property under the given key, if it doesn't
// exist, it will return the provided default.
func GetStringDefault(json Node, key string, def string) string {
	v, ok := GetString(json, key)
	if ok {
		return v
	}
	return def
}

// GetNumber returns a float property under the given key.
func GetNumber(json Node, key string) (float64, bool) {
	v, ok := json[key]
	if !ok {
		return 0, false
	}
	val, ok := v.(float64)
	if !ok {
		// if it is a string value, cast it to float64
		strval, ok := v.(string)
		if ok {
			val, err := strconv.ParseFloat(strval, 64)
			if err == nil {
				return val, true
			}
		}
		return 0, false
	}
	return val, true
}

// GetNumberDefault returns a float property under the given key, if it doesn't
// exist, it will return the provided default.
func GetNumberDefault(json Node, key string, def float64) float64 {
	v, ok := GetNumber(json, key)
	if ok {
		return v
	}
	return def
}

// GetArray returns an []interface{} property under the given key.
func GetArray(json Node, key string) ([]interface{}, bool) {
	v, ok := json[key]
	if !ok {
		return nil, false
	}
	val, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	return val, true
}

// GetNumberArray returns an []float64 property under the given key.
func GetNumberArray(json Node, key string) ([]float64, bool) {
	v, ok := json[key]
	if !ok {
		return nil, false
	}
	val, ok := v.([]float64)
	if !ok {
		return nil, false
	}
	return val, true
}

// GetStringArray returns an []string property under the given key.
func GetStringArray(json Node, key string) ([]string, bool) {
	v, ok := json[key]
	if !ok {
		return nil, false
	}
	val, ok := v.([]string)
	if !ok {
		return nil, false
	}
	return val, true
}
