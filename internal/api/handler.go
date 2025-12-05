package handler

import (
	"encoding/json"
	"go-sharding-basic/internal/api/response"
	"go-sharding-basic/internal/user"
	"net/http"
)

type Handler struct {
	service user.UserService
}

func NewHandler(service user.UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create-user", h.handleCreateUser)
	mux.HandleFunc("/get-user", h.handleGetUser)
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}

	var req createUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.service.CreateUser(r.Context(), req.User, req.Password)
	if err != nil {
		response.Error(w, http.StatusCreated, "Error saving user", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, "User created successfully", nil, http.StatusCreated)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		response.Error(w, http.StatusBadRequest, "The 'user' parameter is required", 400)
		return
	}

	u, err := h.service.GetUser(r.Context(), username)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error retrieving user", 500)
		return
	}

	if u == nil {
		response.Error(w, http.StatusNotFound, "User not found", 404)
		return
	}

	response.JSON(w, http.StatusOK, "User found", u, 200)
}
