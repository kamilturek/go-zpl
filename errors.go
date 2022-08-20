package zpl

import (
	"errors"
	"fmt"
)

var ErrInvalidDensity = fmt.Errorf("invalid density: must be one of %v", allowedDensities())
var ErrInvalidOutputFormat = fmt.Errorf("invalid output format: must be one of %v", allowedOutputFormats())
var ErrNilInput = errors.New("nil input")
var ErrNilOutput = errors.New("nil output")
