package util

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

const (
	defaultPageSize = 10
	maxPageSize     = 100
)

func ParsePaginationParams(ctx fiber.Ctx) (int, int, map[string]string) {
	page := 1
	pageSize := defaultPageSize
	filters := make(map[string]string)

	if p := ctx.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if s := ctx.Query("page_size"); s != "" {
		if parsedPageSize, err := strconv.Atoi(s); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
			if pageSize > maxPageSize {
				pageSize = maxPageSize
			}
		}
	}

	// Collect additional filter parameters
	for key, value := range ctx.Queries() {
		if key != "page" && key != "page_size" {
			filters[key] = value
		}
	}

	return page, pageSize, filters
}
