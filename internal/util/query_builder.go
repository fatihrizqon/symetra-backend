package util

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ApplySearch appends OR LIKE conditions to db for every search term
// (separated by ";") across all searchable fields defined in qp.Fields.
func ApplySearch(db *gorm.DB, qp *QueryParams) *gorm.DB {
	if qp.Search == "" || len(qp.Fields) == 0 {
		return db
	}

	var conditions []string
	var values []interface{}

	for _, term := range strings.Split(qp.Search, ";") {
		term = strings.TrimSpace(term)
		if term == "" {
			continue
		}
		for _, field := range qp.Fields {
			conditions = append(conditions, "LOWER("+field+") LIKE LOWER(?)")
			values = append(values, "%"+term+"%")
		}
	}

	if len(conditions) == 0 {
		return db
	}

	return db.Where(strings.Join(conditions, " OR "), values...)
}

// ApplySort applies ORDER BY to db using a whitelist of allowed columns.
// allowedColumns maps the public-facing key (from ?sort=) to the actual
// SQL column expression (e.g. "users.created_at").
// Falls back to defaultColumn (with ASC direction) when the requested
// column is not whitelisted or when no sort is requested.
func ApplySort(db *gorm.DB, qp *QueryParams, allowedColumns map[string]string, defaultColumn string) *gorm.DB {
	col, ok := allowedColumns[qp.SortBy]
	if !ok || col == "" {
		col = defaultColumn
	}

	dir := "ASC"
	if strings.ToLower(qp.SortDir) == "desc" {
		dir = "DESC"
	}

	return db.Order(fmt.Sprintf("%s %s", col, dir))
}

// ApplyPagination appends LIMIT and OFFSET to db based on qp.Page and qp.PageSize.
func ApplyPagination(db *gorm.DB, qp *QueryParams) *gorm.DB {
	offset := (qp.Page - 1) * qp.PageSize
	return db.Limit(qp.PageSize).Offset(offset)
}
