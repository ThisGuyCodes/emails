package emails

import (
	"errors"
	"fmt"
)

var (
	ErrInvalid       = errors.New("")
	ErrInvalidLocal  = fmt.Errorf("%wlocal part is invalid", ErrInvalid)
	ErrInvalidDomain = fmt.Errorf("%wdomain part is invalid", ErrInvalid)
	ErrNoAt          = fmt.Errorf("%wno @", ErrInvalid)
)
