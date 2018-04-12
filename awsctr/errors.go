package awsctr

import "github.com/juju/errors"

var (
	LogicErr        = errors.New("awsctr: wrong logic")
	InvalidParamErr = errors.New("awsctr: invalid param")
	NotFoundErr     = errors.New("awsctr: not found")
)
