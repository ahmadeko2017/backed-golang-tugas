package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

func GetPaginationParams(c *gin.Context) (int, int) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// Hardening: Prevent Deep Paging (max offset 10000)
	// Querying deep offsets is performance intensive.
	if (page-1)*limit > 10000 {
		page = 10000 / limit
		if page < 1 {
			page = 1
		}
	}

	return page, limit
}
