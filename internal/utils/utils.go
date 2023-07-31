package utils

import (
	"golang.org/x/exp/constraints"
)

func IntPtrToVal[T constraints.Integer](ptr *T) T {
	if ptr == nil {
		return T(0)
	}
	return *ptr
}

func IntValToPtr[T constraints.Integer](val T, valid bool) *T {
	if !valid {
		return nil
	}
	return &val
}
