package config

import (
	"math"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type PaginatedResponseDataSchema struct {
	Limit       uint `json:"limit" example:"100"`
	CurrentPage uint `json:"current_page" example:"1"`
	LastPage    uint `json:"last_page" example:"50"`
}

func PaginateQueryset(queryset interface{}, fiberCtx *fiber.Ctx) (*PaginatedResponseDataSchema, any, *ErrorResponse) {
	currentPage := fiberCtx.QueryInt("page", 1)
	limit := fiberCtx.QueryInt("limit", 100)

	if currentPage < 1 {
		errData := RequestErr(ERR_INVALID_PAGE, "Invalid Page")
		return nil, nil, &errData
	}

	// Check if page size is provided as an argument
	querysetValue := reflect.ValueOf(queryset)
	itemsCount := querysetValue.Len()
	lastPage := math.Ceil(float64(itemsCount) / float64(limit))
	if lastPage == 0 {
		lastPage = 1
	}
	if currentPage > int(lastPage) {
		errData := RequestErr(ERR_INVALID_PAGE, "Page number is out of range")
		return nil, nil, &errData
	}

	startIndex := (currentPage - 1) * limit
	endIndex := startIndex + limit

	// Ensure startIndex is within bounds
	if startIndex < 0 {
		startIndex = 0
	}

	// Ensure endIndex is within bounds
	if endIndex > itemsCount {
		endIndex = itemsCount
	}

	paginatorData := PaginatedResponseDataSchema{
		Limit:     uint(limit),
		CurrentPage: uint(currentPage),
		LastPage:    uint(lastPage),
	}
	paginatedItems := querysetValue.Slice(startIndex, endIndex).Interface()
	return &paginatorData, paginatedItems, nil
}