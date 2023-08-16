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

func StringPtrToVal(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func BoolPtrToVal(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

func ValToPtr[T any](val T, valid bool) *T {
	if !valid {
		return nil
	}
	return &val
}
