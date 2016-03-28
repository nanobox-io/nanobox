package validate

import (
	"fmt"
)

type (
	validatorFunc func() error

	validationError []error
)

var (
	validators = map[string]validatorFunc{}
)

func Register(name string, validator validatorFunc) {
	validators[name] = validator
}

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

func (self *validationError) Add(err error) {
	tmp := validationError(append([]error(*self), err))
	self = &tmp
}

func (self validationError) Error() string {
	str := ""
	for _, err := range self {
		str += fmt.Sprintf("%s\n", err)
	}
	return str
}
