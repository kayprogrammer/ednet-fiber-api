package config

import (
	"context"
	"log"
	"math"

	"github.com/gofiber/fiber/v2"
)

type PaginationResponse[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	ItemsCount int `json:"items_count"`
	TotalPages int `json:"total_pages"`
	Limit      int `json:"limit"`
}

// PaginateModel - Generic Pagination for any Ent Query
func PaginateModel[T any, Q interface {
	Count(context.Context) (int, error)
	Limit(int) Q
	Offset(int) Q
	All(context.Context) ([]T, error)
}](c *fiber.Ctx, query Q) *PaginationResponse[T] {
	ctx := c.Context()

	// Parse query parameters
	page := c.QueryInt("page", 1)     // Default to page 1
	limit := c.QueryInt("limit", 100) // Default to 100 items per page
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 100
	}

	// Get total count of items
	itemsCount, err := query.Count(ctx)
	if err != nil {
		log.Println("Error getting total count:", err)
		return nil
	}

	// Calculate pagination values
	skip := (page - 1) * limit
	totalPages := int(math.Ceil(float64(itemsCount) / float64(limit)))

	// Fetch paginated items
	items, err := query.
		Limit(limit).
		Offset(skip).
		All(ctx)

	if err != nil {
		log.Println("Error fetching paginated items:", err)
		return nil
	}

	// Return paginated response
	return &PaginationResponse[T]{
		Items:      items,
		Page:       page,
		ItemsCount: itemsCount,
		TotalPages: totalPages,
		Limit:      limit,
	}
}
