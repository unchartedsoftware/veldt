package json

import (
	"encoding/json"
)

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

// GetFloat returns a float property under the given key.
func GetFloat(json map[string]interface{}, path ...string) (float64, bool) {
	v, ok := Get(json, path...)
	if !ok {
		return 0, false
	}
	flt, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return flt, true
}

// GetFloatDefault returns a float property under the given key, if it doesn't
// exist, it will return the provided default.
func GetFloatDefault(json map[string]interface{}, def float64, path ...string) float64 {
	v, ok := GetFloat(json, path...)
	if ok {
		return v
	}
	return def
}

// GetInt returns an int property under the given key.
func GetInt(json map[string]interface{}, path ...string) (int, bool) {
	v, ok := Get(json, path...)
	if !ok {
		return 0, false
	}
	flt, ok := v.(float64)
	if !ok {
		return 0, false
	}
	return int(flt), true
}

// GetIntDefault returns a float property under the given key, if it doesn't
// exist, it will return the provided default.
func GetIntDefault(json map[string]interface{}, def int, path ...string) int {
	v, ok := GetInt(json, path...)
	if ok {
		return v
	}
	return def
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

// GetFloatArray returns a []float64 property under the given key.
func GetFloatArray(json map[string]interface{}, path ...string) ([]float64, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	flts := make([]float64, len(vs))
	for i, v := range vs {
		flt, ok := v.(float64)
		if !ok {
			return nil, false
		}
		flts[i] = flt
	}
	return flts, true
}

// GetIntArray returns an []int64 property under the given key.
func GetIntArray(json map[string]interface{}, path ...string) ([]int, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	ints := make([]int, len(vs))
	for i, v := range vs {
		flt, ok := v.(float64)
		if !ok {
			return nil, false
		}
		ints[i] = int(flt)
	}
	return ints, true
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

// GetBoolArray returns an []bool property under the given key.
func GetBoolArray(json map[string]interface{}, path ...string) ([]bool, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	bools := make([]bool, len(vs))
	for i, v := range vs {
		val, ok := v.(bool)
		if !ok {
			return nil, false
		}
		bools[i] = val
	}
	return bools, true
}

// GetRandomChild returns the first key found in the object that is a nested
// json object.
func GetRandomChild(json map[string]interface{}, path ...string) (string, map[string]interface{}, bool) {
	child, ok := GetChild(json, path...)
	if !ok {
		return "", nil, false
	}
	if len(child) == 0 {
		return "", nil, false
	}
	var value map[string]interface{}
	var key string
	for k, v := range child {
		val, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		value = val
		key = k
		break
	}
	return key, value, true
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

// Marshal marhsals JSON into a byte slice, convenience wrapper for the native
// package so no need to import both and get a name collision.
func Marshal(j interface{}) ([]byte, error) {
	return json.Marshal(j)
}

// Unmarshal unmarshals JSON and returns a newly instantiated map.
func Unmarshal(data []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if nil != err {
		return nil, err
	}
	return m, nil
}

// UnmarshalArray unmarshals an array of JSON and returns a newly instantiated
// array of maps.
func UnmarshalArray(data []byte) ([]map[string]interface{}, error) {
	var arr []map[string]interface{}
	err := json.Unmarshal(data, &arr)
	if nil != err {
		return nil, err
	}
	return arr, nil
}

// Copy will copy the JSON data deeply by value, this process involves
// marshalling and then unmarshalling the data.
func Copy(j map[string]interface{}) (map[string]interface{}, error) {
	bytes, err := Marshal(j)
	if err != nil {
		return nil, err
	}
	return Unmarshal(bytes)
}
