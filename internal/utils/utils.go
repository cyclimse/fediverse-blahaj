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
