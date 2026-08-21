package user

import domain "github.com/21mebrat/lost-found-platform/internal/domain/user"

type ListParams struct {
	Page      int            `json:"page"`
	Limit     int            `json:"limit"`
	Search    string         `json:"search"`
	Status    *domain.Status `json:"status"`
	SortBy    string         `json:"sort_by"`
	SortOrder string         `json:"sort_order"`
}

type PaginatedUsers struct {
	Users      []*domain.User `json:"users"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
}
