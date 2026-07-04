package util

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func ParsePaginationParams(ctx fiber.Ctx) (int, int, map[string]string) {
	page := 1
	pageSize := 10
	filters := make(map[string]string)

	if p := ctx.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil {
			page = parsedPage
		}
	}
	if s := ctx.Query("page_size"); s != "" {
		if parsedPageSize, err := strconv.Atoi(s); err == nil {
			pageSize = parsedPageSize
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
