// Package pagination provides page/per_page parsing and response metadata
// shared by all list endpoints.
package pagination

import "strconv"

// Defaults and bounds for list endpoints.
const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Params are sanitized pagination inputs.
type Params struct {
	Page    int
	PerPage int
}

// Parse builds Params from raw query values, applying defaults and bounds.
// Invalid or out-of-range values silently fall back to defaults so that list
// endpoints never fail on pagination input.
func Parse(pageRaw, perPageRaw string) Params {
	page, err := strconv.Atoi(pageRaw)
	if err != nil || page < 1 {
		page = DefaultPage
	}
	perPage, err := strconv.Atoi(perPageRaw)
	if err != nil || perPage < 1 {
		perPage = DefaultPerPage
	}
	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}
	return Params{Page: page, PerPage: perPage}
}

// Offset returns the SQL offset for the current page.
func (p Params) Offset() int { return (p.Page - 1) * p.PerPage }

// Limit returns the SQL limit for the current page.
func (p Params) Limit() int { return p.PerPage }

// Meta is the pagination block embedded in response envelopes.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta computes response metadata for a total row count.
func NewMeta(p Params, total int64) Meta {
	pages := int((total + int64(p.PerPage) - 1) / int64(p.PerPage))
	if pages < 1 {
		pages = 1
	}
	return Meta{Page: p.Page, PerPage: p.PerPage, Total: total, TotalPages: pages}
}
