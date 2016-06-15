// Package validate ...
package validate

import (
	"fmt"
)

var (
	validators = map[string]validatorFunc{}
)

type (
	validatorFunc   func() error
	validationError []error
)

// Register ...
func Register(name string, validator validatorFunc) {
	validators[name] = validator
}

// Add ...
func (vError *validationError) Add(err error) {
	tmp := validationError(append([]error(*vError), err))
	vError = &tmp
}

// Check ...
func Check(checks ...string) error {
	ve := validationError{}
	for _, check := range checks {
		valFunc, ok := validators[check]
		if ok {
			if err := valFunc(); err != nil {
				ve.Add(err)
			}
		}
	}
	if len(ve) != 0 {
		return ve
	}
	return nil
}

// Error ...
func (vError validationError) Error() (str string) {

	//
	for _, err := range vError {
		str += fmt.Sprintf("%s\n", err)
	}

	return str
}
