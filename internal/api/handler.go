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
		response.Error(w, http.StatusBadRequest, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = h.service.CreateUser(r.Context(), req.User, req.Password)
	if err != nil {
		response.Error(w, http.StatusCreated, "Error al guardar usuario", http.StatusInternalServerError)
		return
	}

	response.JSON(w, http.StatusCreated, "Usuario creado correctamente", nil, http.StatusCreated)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		response.Error(w, http.StatusBadRequest, "El parámetro 'user' es obligatorio", 400)
		return
	}

	u, err := h.service.GetUser(r.Context(), username)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error al obtener el usuario", 500)
		return
	}

	if u == nil {
		response.Error(w, http.StatusNotFound, "Usuario no encontrado", 404)
		return
	}

	response.JSON(w, http.StatusOK, "Usuario encontrado", u, 200)
}
