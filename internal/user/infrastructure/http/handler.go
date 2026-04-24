package userhttp

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/application"
	"github.com/SahidAyala/Event-Streaming-User-Service-Example/internal/user/domain"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateEmailRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CreateUser godoc
// @Summary      Crear usuario
// @Description  Crea un nuevo usuario en el sistema
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body CreateUserRequest true "Datos del usuario"
// @Success      201 {object} UserResponse
// @Failure      400 {string} string "Solicitud inválida"
// @Router       /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.service.CreateUser(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := UserResponse{ID: u.ID, Email: u.Email, Username: u.Username}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

// GetUserById godoc
// @Summary      Obtener usuario por ID
// @Description  Retorna un usuario dado su ID
// @Tags         users
// @Produce      json
// @Param        id path string true "ID del usuario"
// @Success      200 {object} UserResponse
// @Failure      400 {string} string "ID faltante"
// @Failure      404 {string} string "Usuario no encontrado"
// @Router       /users/{id} [get]
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	u, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	res := UserResponse{ID: u.ID, Email: u.Email, Username: u.Username}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// UpdateEmail godoc
// @Summary      Actualizar email
// @Description  Actualiza el email de un usuario
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID del usuario"
// @Param        body body UpdateEmailRequest true "Nuevo email"
// @Success      200 {object} UserResponse
// @Failure      400 {string} string "Solicitud inválida"
// @Failure      404 {string} string "Usuario no encontrado"
// @Router       /users/{id}/email [patch]
func (h *Handler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	var req UpdateEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	u, err := h.service.UpdateEmail(r.Context(), id, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := UserResponse{ID: u.ID, Email: u.Email, Username: u.Username}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// UpdatePassword godoc
// @Summary      Actualizar contraseña
// @Description  Actualiza la contraseña de un usuario verificando la actual
// @Tags         users
// @Accept       json
// @Param        id path string true "ID del usuario"
// @Param        body body UpdatePasswordRequest true "Contraseñas"
// @Success      204
// @Failure      400 {string} string "Solicitud inválida"
// @Failure      404 {string} string "Usuario no encontrado"
// @Router       /users/{id}/password [patch]
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	var req UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePassword(r.Context(), id, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
