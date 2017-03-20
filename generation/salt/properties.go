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
