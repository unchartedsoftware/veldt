package salt

import (
	"fmt"
	"strings"
)

// Routines to make it easy to set property values for use by salt

// setProperty sets the specified (multi-leveled) key in the given property map
func setProperty(key string, value interface{}, props map[string]interface{}) error {
	if 0 == len(key) {
		return fmt.Errorf("Null key given to SetKey")
	}

	subKeys := strings.Split(key, ".")
	if 1 == len(subKeys) {
		// Direct property
		props[key] = value
		return nil
	}

	// Deep property requested
	localKey := subKeys[0]
	descendentKey := strings.Join(subKeys[1:], ".")

	localValue, ok := props[localKey]
	if !ok {
		localValue = make(map[string]interface{})
		props[localKey] = localValue
	}
	subConfig, ok := localValue.(map[string]interface{})
	if !ok {
		return fmt.Errorf("Value found under key %s [%v] was not a map (full  key was %s)", localKey, localValue, key)
	}
	return setProperty(descendentKey, value, subConfig)
}

// getProperty gets the value of the specified (multi-leveled) key from the given property map
func getProperty(key string, props map[string]interface{}) (interface{}, error) {
	if 0 == len(key) {
		// whole property set requested
		return props, nil
	}

	subKeys := strings.Split(key, ".")
	if 1 == len(subKeys) {
		// Direct property requested
		return props[key], nil
	}

	// Deep property requested
	localKey := subKeys[0]
	descendentKey := strings.Join(subKeys[1:], ".")

	localValue, ok := props[localKey]
	if !ok {
		return nil, fmt.Errorf("No sub-properties found under key %s (full key was %s)", localKey, key)
	}
	subConfig, ok := localValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Value found under key %s [%v] was not a map (full key was %s)", localKey, localValue, key)
	}

	return getProperty(descendentKey, subConfig)
}

// getFloat32Property gets the value of the specified multi-leveld) key from the given property map, as a float32 value
func getFloat32Property (key string, props map[string]interface{}) (float32, error) {
	rawValue, err := getProperty(key, props)
	if nil != err {
		return 0.0, err
	}

	switch n := rawValue.(type) {
	case int:
		return float32(n), nil
	case int8:
		return float32(n), nil
	case int16:
		return float32(n), nil
	case int32:
		return float32(n), nil
	case int64:
		return float32(n), nil
	case uint:
		return float32(n), nil
	case uint8:
		return float32(n), nil
	case uint16:
		return float32(n), nil
	case uint32:
		return float32(n), nil
	case uint64:
		return float32(n), nil
	case float32:
		return float32(n), nil
	case float64:
		return float32(n), nil
	default:
		return 0.0, fmt.Errorf("Can't convert %v to float", rawValue)
	}
}

// propertiesEqual determines if two property maps are equivalent
func propertiesEqual (a, b map[string]interface{}) bool {
	// Check keys
	if len(a) != len(b) {
		return false
	}
	for k, valA := range a {
		valB, ok := b[k]
		if !ok {
			return false
		}
		subMapA, isSubMapA := valA.(map[string]interface{})
		subMapB, isSubMapB := valB.(map[string]interface{})
		if isSubMapA && isSubMapB {
			if !propertiesEqual(subMapA, subMapB) {
				return false
			}
		} else if isSubMapA || isSubMapB {
			return false
		} else if valA != valB {
			return false
		}
	}
	return true
}
