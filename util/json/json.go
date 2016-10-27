package json

import (
	"bytes"
	"sort"
	"strconv"
)

// Set sets the value under a given path, creating intermediate nodes along the
// way if they do not exist.
func Set(json map[string]interface{}, v interface{}, path ...string) {
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

// Get returns an interface{} under the given path.
func Get(json map[string]interface{}, path ...string) (interface{}, bool) {
	child := json
	last := len(path) - 1
	var val interface{} = child
	for index, key := range path {
		// does a child exists?
		v, ok := child[key]
		if !ok {
			return nil, false
		}
		// is it the target?
		if index == last {
			val = v
			break
		}
		// if not, does it have children to traverse?
		c, ok := v.(map[string]interface{})
		if !ok {
			return nil, false
		}
		child = c
	}
	return val, true
}

// Exists returns true if something exists under the provided path.
func Exists(json map[string]interface{}, path ...string) bool {
	_, ok := Get(json, path...)
	return ok
}

// GetChild returns the child under the given path.
func GetChild(json map[string]interface{}, path ...string) (map[string]interface{}, bool) {
	c, ok := Get(json, path...)
	if !ok {
		return nil, false
	}
	child, ok := c.(map[string]interface{})
	if !ok {
		return nil, false
	}
	return child, true
}

// GetRandomChild returns the first key found in the object that is a nested
// json object.
func GetRandomChild(json map[string]interface{}, path ...string) (map[string]interface{}, bool) {
	child, ok := GetChild(json, path...)
	if !ok {
		return nil, false
	}
	if len(child) == 0 {
		return nil, false
	}
	for _, v := range child {
		val, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		child = val
		break
	}
	return child, true
}

// GetChildOrEmpty returns the child under the given path, if it doesn't
// exist, it will return the provided default.
func GetChildOrEmpty(json map[string]interface{}, path ...string) map[string]interface{} {
	v, ok := GetChild(json, path...)
	if ok {
		return v
	}
	return make(map[string]interface{})
}

// GetString returns a string property under the given path.
func GetString(json map[string]interface{}, path ...string) (string, bool) {
	v, ok := Get(json, path...)
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
func GetStringDefault(json map[string]interface{}, def string, path ...string) string {
	v, ok := GetString(json, path...)
	if ok {
		return v
	}
	return def
}

// GetBool returns a bool property under the given key.
func GetBool(json map[string]interface{}, path ...string) (bool, bool) {
	v, ok := Get(json, path...)
	if !ok {
		return false, false
	}
	val, ok := v.(bool)
	if !ok {
		return false, false
	}
	return val, true
}

// GetBoolDefault returns a bool property under the given key, if it doesn't
// exist, it will return the provided default.
func GetBoolDefault(json map[string]interface{}, def bool, path ...string) bool {
	v, ok := GetBool(json, path...)
	if ok {
		return v
	}
	return def
}

// GetNumber returns a float property under the given key.
func GetNumber(json map[string]interface{}, path ...string) (float64, bool) {
	v, ok := Get(json, path...)
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
func GetNumberDefault(json map[string]interface{}, def float64, path ...string) float64 {
	v, ok := GetNumber(json, path...)
	if ok {
		return v
	}
	return def
}

// GetArray returns an []interface{} property under the given key.
func GetArray(json map[string]interface{}, path ...string) ([]interface{}, bool) {
	v, ok := Get(json, path...)
	if !ok {
		return nil, false
	}
	val, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	return val, true
}

// GetChildArray returns a []map[string]interface{} from the given path.
func GetChildArray(json map[string]interface{}, path ...string) ([]map[string]interface{}, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	nodes := make([]map[string]interface{}, len(vs))
	for i, v := range vs {
		val, ok := v.(map[string]interface{})
		if !ok {
			return nil, false
		}
		nodes[i] = val
	}
	return nodes, true
}

// GetChildMap returns a map[string]map[string]interface{} of all child nodes
// under the given path.
func GetChildMap(json map[string]interface{}, path ...string) (map[string]map[string]interface{}, bool) {
	sub, ok := GetChild(json, path...)
	if !ok {
		return nil, false
	}
	children := make(map[string]map[string]interface{}, len(sub))
	for k, v := range sub {
		c, ok := v.(map[string]interface{})
		if ok {
			children[k] = c
		}
	}
	return children, true
}

// GetNumberArray returns an []float64 property under the given key.
func GetNumberArray(json map[string]interface{}, path ...string) ([]float64, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	flts := make([]float64, len(vs))
	for i, v := range vs {
		val, ok := v.(float64)
		if !ok {
			return nil, false
		}
		flts[i] = val
	}
	return flts, true
}

// GetStringArray returns an []string property under the given key.
func GetStringArray(json map[string]interface{}, path ...string) ([]string, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	strs := make([]string, len(vs))
	for i, v := range vs {
		val, ok := v.(string)
		if !ok {
			return nil, false
		}
		strs[i] = val
	}
	return strs, true
}

// GetHash returns a deterministic hash of unmarshalled json
func GetHash(data interface{}) string {
	var buffer bytes.Buffer
	getHash(data, &buffer)
	return buffer.String()
}

func getHash(data interface{}, buffer *bytes.Buffer) {
	switch data.(type) {
	case map[string]interface{}:
		jMap := data.(map[string]interface{})
		var keys []string
		for k := range jMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		buffer.WriteString("{")
		for _, k := range keys {
			buffer.WriteString(k)
			buffer.WriteString(":")
			getHash(jMap[k], buffer)
			buffer.WriteString(",")
		}
		buffer.WriteString("}")
	case []interface{}:
		jArray := data.([]interface{})
		for _, k := range jArray {
			getHash(k, buffer)
			buffer.WriteString(",")
		}
	case float64:
		val := data.(float64)
		buffer.WriteString(strconv.FormatFloat(val, 'f', 6, 64))
	case string:
		val := data.(string)
		buffer.WriteString(val)
	case bool:
		val := data.(bool)
		buffer.WriteString(strconv.FormatBool(val))
	}
}
