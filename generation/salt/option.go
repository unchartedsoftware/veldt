package salt

import (
	"errors"
)

// Option is a structure that indicates a value that may or may not be defined
type Option struct {
	defined bool
	value   interface{}
}

// IsDefined is true if the optional value is defined
func (opt Option) IsDefined() bool {
	return opt.defined
}

// IsEmpty is true if the optional value is not defined
func (opt Option) IsEmpty() bool {
	return !opt.defined
}

// OrElse returns the value of the optional, or the given default value if the value is not defined.
func (opt Option) OrElse(defaultValue interface{}) interface{} {
	if opt.defined {
		return opt.value
	}
	return defaultValue
}

// Value returns the value of the Option, or an error if the option is not defined
func (opt Option) Value() (interface{}, error) {
	if opt.defined {
		return opt.value, nil
	}
	return nil, errors.New("Undefined option")
}
