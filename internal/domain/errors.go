package domain

import "fmt"

var (
	ErrInvariant     = fmt.Errorf("invariant violation")
	ErrInvalidPeriod = fmt.Errorf("%w: invalid period", ErrInvariant)
)
