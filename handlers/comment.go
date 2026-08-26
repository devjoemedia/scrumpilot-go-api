package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/devjoemedia/scrumpilot-go-api/database"
	"github.com/devjoemedia/scrumpilot-go-api/middleware"
	"github.com/devjoemedia/scrumpilot-go-api/models"
	"github.com/devjoemedia/scrumpilot-go-api/types"
	api_response "github.com/devjoemedia/scrumpilot-go-api/types/response"
	"github.com/devjoemedia/scrumpilot-go-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CreateComment godoc
// @Summary      Create a new comment
// @Description  Create a new comment item
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security 		 BearerAuth
// @Param        body  body     models.CreateCommentRequest  true  "Comment object"
// @Success      200   {object} api_response.CreateCommentResponse
// @Failure      400   {string} string      "Invalid JSON"
// @Router       /api/v1/comments [post]
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid input")
		return
	}

	userID, _, err := middleware.GetUserIDAndEmailFromRequest(r)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to get user ID and email")
		return
	}

	comment := models.Comment{
		Comment:  req.Comment,
		UserID:   userID,
		TicketID: *req.TicketID,
	}

	ctx := r.Context()
	result := database.DB.WithContext(ctx).Create(&comment)

	if result.Error != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	// IMPORTANT: Reload with relations
	if err := database.DB.WithContext(ctx).
		Preload("User").
		First(&comment, comment.ID).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to load comment relations")
		return
	}

	response := api_response.CreateCommentResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Comment created successfully",
		Comment: &comment,
	}

	utils.JSON(w, http.StatusOK, response)
}

// GetComments godoc
// @Summary      Get comments with pagination and filters
// @Description  Retrieve comments with optional pagination (page, size) and filter by user_id
// @Tags         comments
// @Produce      json
// @Security 		 BearerAuth
// @Param        page        query    int    false  "Page number (default: 1)"
// @Param        size        query    int    false  "Page size (default: 10, max: 100)"
// @Param        user_id     query    int      false  "Filter by user_id"
// @Param        priority    query    string   false  "Filter by priority"
// @Param        assignee_id query    int      false  "Filter by assignee_id"
// @Param        reporter_id query    int      false  "Filter by reporter_id"
// @Success      200 {object} api_response.GetCommentsResponse
// @Router       /api/v1/comments [get]
func GetComments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Pagination
	page := 1
	size := 10

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}
	if s, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil && s > 0 {
		size = s
	}

	if size > 100 {
		size = 100
	}

	// Filters
	ticketIDStr := r.URL.Query().Get("ticket_id")

	// Initialize comments slice
	var comments []models.Comment

	// Build Base Query
	query := database.DB.WithContext(ctx).Model(&models.Comment{})

	// Apply filter
	if ticketIDStr != "" {
		query = query.Where("ticket_id = ?", ticketIDStr)
	}

	// Count total comments
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Execute query
	if err := query.
		Preload("User").
		Offset((page - 1) * size).
		Limit(size).
		Find(&comments).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return response
	utils.JSON(w, http.StatusOK, api_response.GetCommentsResponse{
		Success:  true,
		Status:   http.StatusOK,
		Message:  "success",
		Comments: comments,
		Pagination: types.Pagination{
			Total: int(total),
			Size:  size,
			Page:  page,
		},
	})
}

// GetCommentByID godoc
// @Summary      Get comment by ID
// @Description  Fetch a specific comment by its ID
// @Tags         comments
// @Produce      json
// @Security 		 BearerAuth
// @Param        id   path    int  true  "Comment ID"
// @Success      200  {object} api_response.GetCommentResponse
// @Failure      400  {string} string     "Invalid ID"
// @Failure      404  {string} string     "Comment not found"
// @Router       /api/v1/comments/{id} [get]
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Comment Slice
	var comment models.Comment
	if err := database.DB.WithContext(ctx).Preload("User").First(&comment, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Comment not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Return response
	utils.JSON(w, http.StatusOK, api_response.GetCommentResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Comment retrieved successfully",
		Comment: &comment,
	})
}

// UpdateComment godoc
// @Summary      Update an existing comment
// @Description  Update a comment item by ID
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security 		 BearerAuth
// @Param        id     path    int  true  "Comment ID"
// @Param        body   body    models.UpdateCommentRequest  true  "Comment object"
// @Success      200    {object} api_response.UpdateCommentResponse
// @Failure      400    {string} string  "Invalid JSON"
// @Failure      404    {string} string  "Comment not found"
// @Router       /api/v1/comments/{id} [patch]
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req models.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Comment Slice
	var comment models.Comment
	if err := database.DB.WithContext(ctx).First(&comment, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Comment not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Update fields
	if req.Comment != "" {
		comment.Comment = req.Comment
	}

	// save updates
	result := database.DB.WithContext(ctx).Save(&comment).Error
	if result != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update comment")
		return
	}

	// Response
	response := api_response.UpdateCommentResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Comment updated successfully",
		Comment: &comment,
	}

	utils.JSON(w, http.StatusOK, response)
}

// DeleteComment godoc
// @Summary      Delete comment by ID
// @Description  Delete a specific comment by its ID
// @Tags         comments
// @Produce      json
// @Security 		 BearerAuth
// @Param        id   path    int  true  "Comment ID"
// @Success      200  {object} api_response.DeleteCommentResponse
// @Failure      400  {string} string             "Invalid ID"
// @Failure      404  {string} string             "Comment not found"
// @Router       /api/v1/comments/{id} [delete]
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Delete comment
	if err := database.DB.WithContext(ctx).Delete(&models.Comment{}, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Comment not found")
			return
		}

		utils.Error(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	// Response
	response := api_response.DeleteCommentResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Comment deleted successfully",
	}

	utils.JSON(w, http.StatusOK, response)
}
