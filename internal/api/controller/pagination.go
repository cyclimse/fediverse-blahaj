package controller

import "golang.org/x/exp/constraints"

const (
	defaultPageSize = 30
	minimumPageSize = 1
	maximumPageSize = 100 // to prevent abuses

	defaultPage = 1
	minimumPage = defaultPage
)

func validatePage[T constraints.Integer](page *T) T {
	if page == nil {
		return defaultPage
	}
	if *page < T(1) {
		return 1
	}
	return *page
}

func validatePageSize[T constraints.Integer](pageSize *T) T {
	if pageSize == nil {
		return defaultPageSize
	}
	if *pageSize < minimumPageSize {
		return minimumPageSize
	}
	if *pageSize > maximumPageSize {
		return maximumPageSize
	}
	return *pageSize
}
