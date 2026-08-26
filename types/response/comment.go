package api_response

import (
	"github.com/devjoemedia/scrumpilot-go-api/models"
	"github.com/devjoemedia/scrumpilot-go-api/types"
)

type GetCommentsResponse struct {
	Success    bool             `json:"success"`
	Status     int              `json:"status"`
	Message    string           `json:"message"`
	Comments   []models.Comment `json:"comments"`
	Pagination types.Pagination `json:"pagination"`
}

type GetCommentResponse struct {
	Success bool            `json:"success"`
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Comment *models.Comment `json:"comment"`
}

type CreateCommentResponse struct {
	Success bool            `json:"success"`
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Comment *models.Comment `json:"comment"`
}

type UpdateCommentResponse struct {
	Success bool            `json:"success"`
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Comment *models.Comment `json:"comment"`
}

type DeleteCommentResponse struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
