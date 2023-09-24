package utils

import (
	"errors"

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

func ConvertIntPtrToInt32Ptr[T constraints.Integer](i *T) *int32 {
	if i == nil {
		return nil
	}
	v := int32(*i)
	return &v
}

func NullableStringToError(s string, valid bool) error {
	if !valid {
		return nil
	}
	return errors.New(s)
}
