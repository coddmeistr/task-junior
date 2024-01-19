package repository

import (
	"fmt"
)

type sortOptions struct {
	Field string
	Order string
}

type paginateOptions struct {
	Page    int
	PerPage int
}

func NewSortOptions(field, order string) SortOptions {
	return &sortOptions{
		Field: field,
		Order: order,
	}
}

func (options *sortOptions) GetOrderBy() string {
	return fmt.Sprintf("%s %s", options.Field, options.Order)
}

func NewPaginateOptions(page, perPage int) PaginateOptions {
	return &paginateOptions{
		Page:    page,
		PerPage: perPage,
	}
}

func (options *paginateOptions) GetPage() uint64 {
	return uint64(options.Page)
}

func (options *paginateOptions) GetPerPage() uint64 {
	return uint64(options.PerPage)
}
