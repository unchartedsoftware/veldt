package json

// Node represents a single json node as a map[string]interface{}
type Node map[string]interface{}

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

// GetChildren returns a map of all child nodes.
func GetChildren(json Node) map[string]Node {
	children := make(map[string]Node, len(json))
	for k, v := range json {
		c, ok := v.(map[string]interface{})
		if ok {
			children[k] = c
		}
	}
	return children
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
	v, ok := json[key]
	if !ok {
		return def
	}
	val, ok := v.(string)
	if !ok {
		return def
	}
	return val
}

// GetNumber returns a float property under the given key.
func GetNumber(json Node, key string) (float64, bool) {
	v, ok := json[key]
	if !ok {
		return 0, false
	}
	val, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return val, true
}

// GetNumberDefault returns a float property under the given key, if it doesn't
// exist, it will return the provided default.
func GetNumberDefault(json Node, key string, def float64) float64 {
	v, ok := json[key]
	if !ok {
		return def
	}
	val, ok := v.(float64)
	if !ok {
		return def
	}
	return val
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
