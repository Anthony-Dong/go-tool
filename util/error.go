package util

import "github.com/juju/errors"

var (
	NilError = func(name string) error {
		return errors.Errorf("the %s is nil", name)
	}
)
