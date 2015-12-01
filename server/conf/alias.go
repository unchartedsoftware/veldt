package conf

import (
	"fmt"
	"strings"
)

// AliasMap represents a set of aliased string values.
type AliasMap map[string]string

// String will return the string representation of the map.
func (a *AliasMap) String() string {
	var strs []string
	for key, val := range *a {
		strs = append(strs, key+": "+val)
	}
	return strings.Join(strs, ", ")
}

// CheckExisting returns an error if a key/value does not result in a 1:1 mapping.
func (a *AliasMap) CheckExisting(key string, val string) error {
	for k, v := range *a {
		if k == key {
			return fmt.Errorf("Key '%s' has already been aliased, duplicate aliases are not supported.", key)
		}
		if v == val {
			return fmt.Errorf("Value '%s' has already been aliased under '%s', duplicate alias mappings are not supported.", val, k)
		}
	}
	return nil
}

// Set will set an individual alias to the map.
func (a *AliasMap) Set(arg string) error {
	args := strings.Split(arg, "=")
	key := args[0]
	val := args[1]
	err := a.CheckExisting(key, val)
	if err != nil {
		return err
	}
	(*a)[key] = val
	return nil
}

// Unalias attempts to unalias a string, if an alias doesn't exist, returns unmodified argument.
func Unalias(alias string) string {
	value, ok := GetConf().Aliases[alias]
	if ok {
		return value
	}
	return alias
}

// Alias attempts to alias a string, if the value is not aliased, returns unmodified argument.
func Alias(value string) string {
	alias, ok := GetConf().InvAliases[value]
	if ok {
		return alias
	}
	return value
}

func invertAliases(aliases AliasMap) AliasMap {
	values := make(AliasMap)
	for k, v := range aliases {
		values[v] = k
	}
	return values
}
