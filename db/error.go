package db

import "errors"

var (
	// ErrNil is record not exist
	ErrNil = errors.New("record not exist")
	// ErrParam is invalid parameter
	ErrParam = errors.New("invalid parameter")
	// ErrConnect is database connect error
	ErrConnect = errors.New("db connect error")
)
