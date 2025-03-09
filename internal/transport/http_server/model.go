package http_server

import (
	"fmt"
)

type ExpenseRequest struct {
	Amount     int    `json:"amount"`
	CategoryID int    `json:"category_id"`
	Currency   string `json:"currency"`
	ActionDate string `json:"action_date"`
	Note       string `json:"note"`
}

func (e *ExpenseRequest) Validate() error {
	if e.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if e.CategoryID <= 0 {
		return fmt.Errorf("category_id must be greater than 0")
	}
	if e.ActionDate == "" {
		return fmt.Errorf("action_date must be not empty")
	}

	return nil
}

type ExpenseResponse struct {
	ID         int    `json:"id"`
	Amount     int    `json:"amount"`
	Currency   string `json:"currency"`
	CategoryID int    `json:"category_id"`
	ActionDate string `json:"action_date"`
	Note       string `json:"note"`
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *CategoryRequest) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name must be not empty")
	}
	if c.Description == "" {
		return fmt.Errorf("description must be not empty")
	}

	return nil
}

type CategoryResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
