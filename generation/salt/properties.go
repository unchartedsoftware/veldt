package salt

import (
	"fmt"
	"strings"
)

const (
	boolArray = 37 + iota
	intArray
	int8Array
	int16Array
	int32Array
	int64Array
	uint8Array
	uint16Array
	uint32Array
	uint64Array
	float32Array
	float64Array
	stringArray
	interfaceArray
	notAnArray
)

// Routines to make it easy to set property values for use by salt

// setProperty sets the specified (multi-leveled) key in the given property map
func setProperty(key string, value interface{}, props map[string]interface{}) error {
	if 0 == len(key) {
		return fmt.Errorf("null key given to SetKey")
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
		return fmt.Errorf("value found under key %s [%v] was not a map (full  key was %s)", localKey, localValue, key)
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
		return nil, fmt.Errorf("no sub-properties found under key %s (full key was %s)", localKey, key)
	}
	subConfig, ok := localValue.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("value found under key %s [%v] was not a map (full key was %s)", localKey, localValue, key)
	}

	return getProperty(descendentKey, subConfig)
}

// getFloat32Property gets the value of the specified multi-leveld) key from the given property map, as a float32 value
func getFloat32Property(key string, props map[string]interface{}) (float32, error) {
	rawValue, err := getProperty(key, props)
	if err != nil {
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
		return 0.0, fmt.Errorf("can't convert %v to float", rawValue)
	}
}

// propertiesEqual determines if two property maps are equivalent
func propertiesEqual(a, b map[string]interface{}) bool {
	// Check keys
	if len(a) != len(b) {
		return false
	}
	for k, valA := range a {
		valB, ok := b[k]
		if !ok {
			return false
		}
		if !propertyElementsEqual(valA, valB) {
			return false
		}
	}
	return true
}

func propertyArraysEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, valA := range a {
		valB := b[i]
		if !propertyElementsEqual(valA, valB) {
			return false
		}
	}
	return true
}

func propertyElementsEqual(a, b interface{}) bool {
	mapA, isMapA := a.(map[string]interface{})
	mapB, isMapB := b.(map[string]interface{})
	if isMapA && isMapB {
		return propertiesEqual(mapA, mapB)
	} else if isMapA || isMapB {
		return false
	}

	aArrayType := arrayType(a)
	bArrayType := arrayType(b)
	if aArrayType == bArrayType {
		if notAnArray == aArrayType {
			return a == b
		}
		return compareArrays(aArrayType, a, b)
	}
	return false
}

func arrayType(a interface{}) int {
	_, isBoolArray := a.([]bool)
	if isBoolArray {
		return boolArray
	}
	_, isIntArray := a.([]int)
	if isIntArray {
		return intArray
	}
	_, isInt8Array := a.([]int8)
	if isInt8Array {
		return int8Array
	}
	_, isInt16Array := a.([]int16)
	if isInt16Array {
		return int16Array
	}
	_, isInt32Array := a.([]int32)
	if isInt32Array {
		return int32Array
	}
	_, isInt64Array := a.([]int64)
	if isInt64Array {
		return int64Array
	}
	_, isUInt8Array := a.([]uint8)
	if isUInt8Array {
		return uint8Array
	}
	_, isUInt16Array := a.([]uint16)
	if isUInt16Array {
		return uint16Array
	}
	_, isUInt32Array := a.([]uint32)
	if isUInt32Array {
		return uint32Array
	}
	_, isUInt64Array := a.([]uint64)
	if isUInt64Array {
		return uint64Array
	}
	_, isFloat32Array := a.([]float32)
	if isFloat32Array {
		return float32Array
	}
	_, isFloat64Array := a.([]float64)
	if isFloat64Array {
		return float64Array
	}
	_, isStringArray := a.([]string)
	if isStringArray {
		return stringArray
	}
	_, isInterfaceArray := a.([]interface{})
	if isInterfaceArray {
		return interfaceArray
	}
	return notAnArray
}

func compareArrays(arrayType int, a, b interface{}) bool {
	switch arrayType {
	case boolArray:
		aArray, _ := a.([]bool)
		bArray, _ := b.([]bool)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case intArray:
		aArray, _ := a.([]int)
		bArray, _ := b.([]int)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case int8Array:
		aArray, _ := a.([]int8)
		bArray, _ := b.([]int8)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case int16Array:
		aArray, _ := a.([]int16)
		bArray, _ := b.([]int16)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case int32Array:
		aArray, _ := a.([]int32)
		bArray, _ := b.([]int32)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case int64Array:
		aArray, _ := a.([]int64)
		bArray, _ := b.([]int64)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case uint8Array:
		aArray, _ := a.([]uint8)
		bArray, _ := b.([]uint8)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case uint16Array:
		aArray, _ := a.([]uint16)
		bArray, _ := b.([]uint16)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case uint32Array:
		aArray, _ := a.([]uint32)
		bArray, _ := b.([]uint32)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case uint64Array:
		aArray, _ := a.([]uint64)
		bArray, _ := b.([]uint64)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case float32Array:
		aArray, _ := a.([]float32)
		bArray, _ := b.([]float32)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case float64Array:
		aArray, _ := a.([]float64)
		bArray, _ := b.([]float64)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case stringArray:
		aArray, _ := a.([]string)
		bArray, _ := b.([]string)
		if len(aArray) != len(bArray) {
			return false
		}
		for i := 0; i < len(aArray); i++ {
			if aArray[i] != bArray[i] {
				return false
			}
		}
		return true
	case interfaceArray:
		aArray, _ := a.([]interface{})
		bArray, _ := b.([]interface{})
		return propertyArraysEqual(aArray, bArray)
	}

	// Should always be one of those cases
	Warnf("Unrecognized array type %v", arrayType)
	return false
}
