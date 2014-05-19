package form

import (
	"errors"
	"strconv"
)

// TODO better error interface

type Rule struct {
	Name       string
	Parameters string
	Function   func(string) error
}

func MinLength(parameters string) *Rule {
	minlen, _ := strconv.Atoi(parameters)
	return &Rule{
		Name:       "minlength",
		Parameters: parameters,
		Function: func(value string) error {
			if len(value) >= minlen {
				return nil
			} else {
				return errors.New("minlength")
			}
		},
	}
}

func MaxLength(parameters string) *Rule {
	maxlen, _ := strconv.Atoi(parameters)
	return &Rule{
		Name:       "maxlength",
		Parameters: parameters,
		Function: func(value string) error {
			if len(value) <= maxlen {
				return nil
			} else {
				return errors.New("maxlength")
			}
		},
	}
}
