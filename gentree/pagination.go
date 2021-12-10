package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
)

const (
	minPageSize = 2
	maxPageSize = 100
)

type paginationData struct {
	PageIdx  int
	PageSize int
	TotalCnt int
}

func (p *paginationData) getJson(baseUrl url.URL) gin.H {
	//	prevPageIdx := maxInt(0, p.PageIdx-1)
	//	nextPageIdx := 0

	return gin.H{}
}

func (p *paginationData) validate() error {
	if p.PageIdx < 0 {
		return AppError{
			errInvalidArgument,
			fmt.Sprintf("The page index is negative (%d)", p.PageIdx)}
	}

	if (p.PageSize < minPageSize) || (p.PageSize > maxPageSize) {
		return AppError{
			errInvalidArgument,
			fmt.Sprintf("The page size (%d) is out of bounds ([%d, %d])",
				p.PageSize, minPageSize, maxPageSize)}
	}

	return nil
}
