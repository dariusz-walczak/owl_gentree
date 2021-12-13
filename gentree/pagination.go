package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strconv"
)

const (
	minPageSize = 10
	maxPageSize = 100
)

type paginationData struct {
	PageIdx     int
	PageSize    int
	TotalCnt    int
	minPageSize int
	maxPageSize int
}

func composePageUrl(baseUrl url.URL, pageIdx int, pageSize int) string {
	query := baseUrl.Query()
	query.Set("page", strconv.Itoa(pageIdx))
	query.Set("limit", strconv.Itoa(pageSize))
	baseUrl.RawQuery = query.Encode()

	return baseUrl.String()
}

func (p *paginationData) getJson(baseUrl url.URL) gin.H {
	log.Trace("Entry checkpoint")

	json := gin.H{}

	if p.PageIdx-1 >= 0 {
		json["prev_url"] = composePageUrl(baseUrl, p.PageIdx-1, p.PageSize)
	}

	pageCnt := int(maxInt(p.TotalCnt-1, 0) / p.PageSize)

	if p.PageIdx+1 <= pageCnt {
		json["next_url"] = composePageUrl(baseUrl, p.PageIdx+1, p.PageSize)
	}

	return json
}

func (p *paginationData) validate() error {
	if p.PageIdx < 0 {
		return AppError{
			errInvalidArgument,
			fmt.Sprintf("The page index is negative (%d)", p.PageIdx)}
	}

	if (p.PageSize < p.minPageSize) || (p.PageSize > p.maxPageSize) {
		return AppError{
			errInvalidArgument,
			fmt.Sprintf("The page size (%d) is out of bounds ([%d, %d])",
				p.PageSize, p.minPageSize, p.maxPageSize)}
	}

	return nil
}
