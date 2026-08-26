package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/devjoemedia/scrumpilot-go-api/database"
	"github.com/devjoemedia/scrumpilot-go-api/models"
	"github.com/devjoemedia/scrumpilot-go-api/types"
	api_response "github.com/devjoemedia/scrumpilot-go-api/types/response"
	"github.com/devjoemedia/scrumpilot-go-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CreateTodo godoc
// @Summary      Create a new todo
// @Description  Create a new todo item
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security 		 BearerAuth
// @Param        body  body     models.Todo  true  "Todo object"
// @Success      200   {object} api_response.CreateTodoResponse
// @Failure      400   {string} string      "Invalid JSON"
// @Router       /api/v1/todos [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Validation failed")
		return
	}

	// Map to database model (GORM fills ID, timestamps)
	todo := models.Todo{
		Title:       req.Title,
		Description: req.Description,
		IsCompleted: req.IsCompleted,
	}

	ctx := r.Context()
	result := database.DB.WithContext(ctx).Create(&todo)
	if result.Error != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to create todo")
		return
	}

	response := api_response.CreateTodoResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Todo created successfully",
		Todo:    &todo,
	}

	utils.JSON(w, http.StatusOK, response)
}

// GetTodos godoc
// @Summary      Get todos with pagination and filters
// @Description  Retrieve todos with optional pagination (page, size) and filter by is_completed
// @Tags         todos
// @Produce      json
// @Security 		 BearerAuth
// @Param        page        query    int    false  "Page number (default: 1)"
// @Param        size        query    int    false  "Page size (default: 10, max: 100)"
// @Param        is_completed query   bool   false  "Filter by completion status"
// @Success      200 {object} api_response.GetTodosResponse
// @Router       /api/v1/todos [get]
func GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Pagination
	page := 1
	size := 10

	// Parse query params with defaults
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
	is_completed := r.URL.Query().Get("is_completed")

	// Initialize todos slice
	var todos []models.Todo

	// Build Base Query
	query := database.DB.WithContext(ctx).Model(&models.Todo{})

	// Apply filter
	if is_completed != "" {
		if parsedValue, err := strconv.ParseBool(is_completed); err == nil {
			query.Where("is_completed = ?", parsedValue)
		}
	}

	// Pagination
	total := int64(0)
	if err := query.Count(&total).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to fetch todos")
		return
	}

	if err := query.
		Offset((page - 1) * size).
		Limit(size).
		Find(&todos).Error; err != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to fetch todos")
		return
	}

	// Response with pagination metadata
	response := api_response.GetTodosResponse{
		Message: "Todos retrieved successfully",
		Success: false,
		Status:  http.StatusOK,
		Todos:   todos,
		Pagination: types.Pagination{
			Page:  page,
			Size:  size,
			Total: int(total),
			Pages: int((total + int64(size) - 1) / int64(size)),
		},
	}

	utils.JSON(w, http.StatusOK, response)
}

// GetTodoByID godoc
// @Summary      Get todo by ID
// @Description  Fetch a specific todo by its ID
// @Tags         todos
// @Produce      json
// @Security 		 BearerAuth
// @Param        id   path    int  true  "Todo ID"
// @Success      200  {object} api_response.GetTodoResponse
// @Failure      400  {string} string     "Invalid ID"
// @Failure      404  {string} string     "Todo not found"
// @Router       /api/v1/todos/{id} [get]
func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var todo models.Todo
	if err := database.DB.WithContext(ctx).First(&todo, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Todo not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Database error")
		return
	}

	response := api_response.GetTodoResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Todo retrieved successfully",
		Todo:    &todo,
	}

	utils.JSON(w, http.StatusOK, response)
}

// UpdateTodo godoc
// @Summary      Update an existing todo
// @Description  Update a todo item by ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security 		 BearerAuth
// @Param        id     path    int  true  "Todo ID"
// @Param        body   body    models.UpdateTodoRequest  true  "Todo object"
// @Success      200    {object} api_response.UpdateTodoResponse
// @Failure      400    {string} string  "Invalid JSON"
// @Failure      404    {string} string  "Todo not found"
// @Router       /api/v1/todos/{id} [patch]
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	// Get todo ID from URL params
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate input
	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Validation failed")
		return
	}

	ctx := r.Context()
	var todo models.Todo
	result := database.DB.WithContext(ctx).First(&todo, id)
	if result.Error != nil {
		utils.Error(w, http.StatusNotFound, "Todo not found")
		return
	}

	// Update fields
	todo.Title = req.Title
	todo.Description = req.Description
	todo.IsCompleted = req.IsCompleted

	result = database.DB.WithContext(ctx).Save(&todo)
	if result.Error != nil {
		utils.Error(w, http.StatusInternalServerError, "Failed to update todo")
		return
	}

	response := api_response.UpdateTodoResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Todo updated successfully",
		Todo:    &todo,
	}

	utils.JSON(w, http.StatusOK, response)
}

// DeleteTodo godoc
// @Summary      Delete todo by ID
// @Description  Delete a specific todo by its ID
// @Tags         todos
// @Produce      json
// @Security 		 BearerAuth
// @Param        id   path    int  true  "Todo ID"
// @Success      200  {object} api_response.DeleteTodoResponse
// @Failure      400  {string} string             "Invalid ID"
// @Failure      404  {string} string             "Todo not found"
// @Router       /api/v1/todos/{id} [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := database.DB.WithContext(ctx).Delete(&models.Todo{}, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(w, http.StatusNotFound, "Todo not found")
			return
		}
		utils.Error(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	response := api_response.DeleteTodoResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Todo deleted successfully",
	}

	utils.JSON(w, http.StatusOK, response)
}
