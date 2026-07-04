package util

import (
	"fmt"

	"github.com/fatihrizqon/gofiber-microservice/internal/delivery/http/response"
)

// GenerateMeta builds the pagination meta block for list responses.
// It reconstructs query strings from qp so that all active filters,
// sort, and search values are preserved in the next/prev/first/last links.
func GenerateMeta(baseURL string, qp *QueryParams, totalCount int) response.Meta {
	page := qp.Page
	if page < 1 {
		page = 1
	}

	totalPages := (totalCount + qp.PageSize - 1) / qp.PageSize

	buildURL := func(p int) string {
		query := fmt.Sprintf("?page=%d&page_size=%d", p, qp.PageSize)
		if qp.Search != "" {
			query += "&search=" + qp.Search
		}
		if qp.SortBy != "" {
			query += "&sort=" + qp.SortBy + "&order=" + qp.SortDir
		}
		for key, values := range qp.Filters {
			for _, value := range values {
				query += fmt.Sprintf("&%s=%s", key, value)
			}
		}
		return baseURL + query
	}

	currentPage := buildURL(page)
	prevPage := getPageURL(page-1, buildURL)
	nextPage := getPageURL(page+1, buildURL, (page*qp.PageSize) < totalCount)

	startIndex, endIndex := calculateIndices(page, qp.PageSize, totalCount)

	return response.Meta{
		Search:     qp.Search,
		Info:       fmt.Sprintf("Showing %d to %d from %d item(s).", startIndex, endIndex, totalCount),
		Page:       page,
		TotalCount: totalCount,
		TotalPages: totalPages,
		PageSize:   qp.PageSize,
		Links: response.Links{
			CurrentPage: currentPage,
			FirstPage:   buildURL(1),
			LastPage:    buildURL(totalPages),
			NextPage:    nextPage,
			PrevPage:    prevPage,
		},
	}
}

func getPageURL(page int, buildURL func(int) string, condition ...bool) *string {
	if len(condition) == 0 || condition[0] {
		url := buildURL(page)
		return &url
	}
	return nil
}

func calculateIndices(page, pageSize, totalCount int) (int, int) {
	startIndex := (page-1)*pageSize + 1
	endIndex := min(startIndex+pageSize-1, totalCount)
	return startIndex, endIndex
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
