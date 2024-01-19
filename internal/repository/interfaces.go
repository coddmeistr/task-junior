package repository

type SortOptions interface {
	GetOrderBy() string
}

type PaginateOptions interface {
	GetPage() uint64
	GetPerPage() uint64
}
