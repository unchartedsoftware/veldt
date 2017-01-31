package json

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
	return numAsFloat64(v)
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
func GetInt(json map[string]interface{}, path ...string) (int64, bool) {
	v, ok := Get(json, path...)
	if !ok {
		return 0, false
	}
	return numAsInt64(v)
}

// GetIntDefault returns a float property under the given key, if it doesn't
// exist, it will return the provided default.
func GetIntDefault(json map[string]interface{}, def int64, path ...string) int64 {
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
		val, ok := numAsFloat64(v)
		if !ok {
			return nil, false
		}
		flts[i] = val
	}
	return flts, true
}

// GetIntArray returns an []int64 property under the given key.
func GetIntArray(json map[string]interface{}, path ...string) ([]int64, bool) {
	vs, ok := GetArray(json, path...)
	if !ok {
		return nil, false
	}
	flts := make([]int64, len(vs))
	for i, v := range vs {
		val, ok := numAsInt64(v)
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

func numAsFloat64(num interface{}) (float64, bool) {
	switch t := num.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int32:
		return float64(t), true
	case int64:
		return float64(t), true
	case uint:
		return float64(t), true
	case uint32:
		return float64(t), true
	case uint64:
		return float64(t), true
	}
	return 0, false
}

func numAsInt64(num interface{}) (int64, bool) {
	switch t := num.(type) {
	case float64:
		return int64(t), true
	case float32:
		return int64(t), true
	case int:
		return int64(t), true
	case int32:
		return int64(t), true
	case int64:
		return t, true
	case uint:
		return int64(t), true
	case uint32:
		return int64(t), true
	case uint64:
		return int64(t), true
	}
	return 0, false
}
