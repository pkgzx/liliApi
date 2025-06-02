package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkgzx/liliApi/src/internal/middleware"
	"github.com/pkgzx/liliApi/src/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

func NewUserHandler(userService *services.UserService, authService *services.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// Estructuras para requests y responses
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserData struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type UserResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	User    *UserData `json:"user,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// Login handler con JWT
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validar datos requeridos
	if req.Username == "" || req.Password == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Username and password are required", "")
		return
	}

	// Autenticar usuario y generar token
	result, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", "")
		return
	}

	// Respuesta exitosa con token
	response := LoginResponse{
		Token: result.Token,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Refresh token endpoint
func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	type RefreshRequest struct {
		Token string `json:"token"`
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if req.Token == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Token is required", "")
		return
	}

	newToken, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Cannot refresh token", err.Error())
		return
	}

	type RefreshResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	response := RefreshResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Token:   newToken,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Profile endpoint protegido
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	// Obtener usuario del contexto (agregado por el middleware)
	userClaims, ok := middleware.GetUserFromContext(r)
	if !ok {
		h.writeErrorResponse(w, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	response := UserResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		User: &UserData{
			ID:       userClaims.UserID,
			Username: userClaims.Username,
			FullName: userClaims.FullName,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Crear usuario
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validar datos requeridos
	if req.Username == "" || req.Password == "" || req.FullName == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Username, password and full_name are required", "")
		return
	}

	// Validaciones básicas
	if len(req.Username) < 3 {
		h.writeErrorResponse(w, http.StatusBadRequest, "Username must be at least 3 characters", "")
		return
	}

	if len(req.Password) < 6 {
		h.writeErrorResponse(w, http.StatusBadRequest, "Password must be at least 6 characters", "")
		return
	}

	// Crear usuario
	user, err := h.userService.CreateUser(req.Username, req.Password, req.FullName)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			h.writeErrorResponse(w, http.StatusConflict, "Username already exists", "")
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	// Respuesta exitosa
	response := UserResponse{
		Success: true,
		Message: "User created successfully",
		User: &UserData{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Función auxiliar para escribir respuestas de error
func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message, errorDetail string) {
	response := ErrorResponse{
		Success: false,
		Message: message,
		Error:   errorDetail,
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
