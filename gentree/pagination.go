package main

import (
	"fmt"
)

type pagePaginationQuery struct {
	Page  int `form:"page" binding:"min=0"`
	Limit int `form:"limit" binding:"isdefault|min=2,max=100"`
}

func (p *pagePaginationQuery) applyDefaults() {
	if p.Limit == 0 {
		p.Limit = 20
	}
}

const (
	minPageSize = 2
	maxPageSize = 100
)

func checkPaginationParams(pageIdx int, pageSize int) error {
	if pageIdx < 0 {
		return AppError{
			errInvalidArgument,
			fmt.Sprintf("The page index is negative (%d)", pageIdx)}
	}

	if (pageSize < minPageSize) || (pageSize > maxPageSize) {
		return AppError{
			errInvalidArgument,
			fmt.Sprintf("The page size (%d) is out of bounds ([%d, %d])",
				pageSize, minPageSize, maxPageSize)}
	}

	return nil
}
