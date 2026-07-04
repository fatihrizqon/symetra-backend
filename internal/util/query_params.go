package util

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

var reservedKeys = map[string]bool{
	"page":      true,
	"page_size": true,
	"search":    true,
	"sort":      true,
	"order":     true,
}

type QueryParams struct {
	Page     int
	PageSize int
	Search   string
	Fields   []string
	SortBy   string
	SortDir  string
	Filters  map[string][]string
}

func ParseQueryParams(ctx fiber.Ctx, fields []string) *QueryParams {
	qp := &QueryParams{
		Page:     1,
		PageSize: 10,
		Fields:   fields,
		Filters:  make(map[string][]string),
	}

	if p := ctx.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			qp.Page = v
		}
	}

	if s := ctx.Query("page_size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			qp.PageSize = v
		}
	}

	qp.Search = ctx.Query("search")

	if sort := ctx.Query("sort"); sort != "" {
		qp.SortBy = sort
	}

	if ctx.Query("order") == "desc" {
		qp.SortDir = "desc"
	} else {
		qp.SortDir = "asc"
	}

	args := ctx.RequestCtx().QueryArgs()
	args.VisitAll(func(key, value []byte) {
		k := string(key)
		if !reservedKeys[k] {
			for _, part := range strings.Split(string(value), ",") {
				if p := strings.TrimSpace(part); p != "" {
					qp.Filters[k] = append(qp.Filters[k], p)
				}
			}
		}
	})

	return qp
}

// fiber:context-methods migrated
