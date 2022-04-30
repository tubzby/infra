package copy

import (
	"errors"
	"reflect"
)

// SameType .
func SameType(dst, src interface{}) error {
	var (
		dstv = reflect.Indirect(reflect.ValueOf(dst))
		srcv = reflect.Indirect(reflect.ValueOf(src))
	)

	if dstv.Type() != srcv.Type() {
		return errors.New("dst, src not the same type")
	}

	if !dstv.CanSet() {
		return errors.New("dst can't be set")
	}

	dstv.Set(srcv)
	return nil
}
