package cache

import "errors"

var (
	// ErrNil is key not exist
	ErrNil error = errors.New("key not exist")
	// ErrMarshal is marshal error
	ErrMarshal = errors.New("marshal error")
	// ErrUnMarshal is unmarshal error
	ErrUnMarshal = errors.New("unMarshal error")
)
