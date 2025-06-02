package middleware

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"

    "github.com/pkgzx/liliApi/src/internal/services"
)

type AuthMiddleware struct {
    authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
    return &AuthMiddleware{
        authService: authService,
    }
}

type contextKey string

const UserContextKey contextKey = "user"

type ErrorResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Error   string `json:"error,omitempty"`
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        // Obtener token del header Authorization
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            m.writeErrorResponse(w, http.StatusUnauthorized, "Authorization header required", "")
            return
        }

        // Verificar formato Bearer token
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            m.writeErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format", "")
            return
        }

        token := tokenParts[1]

        // Validar token
        claims, err := m.authService.ValidateToken(token)
        if err != nil {
            m.writeErrorResponse(w, http.StatusUnauthorized, "Invalid token", err.Error())
            return
        }

        // Agregar claims al contexto
        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        r = r.WithContext(ctx)

        // Continuar con el siguiente handler
        next(w, r)
    }
}

func (m *AuthMiddleware) writeErrorResponse(w http.ResponseWriter, statusCode int, message, errorDetail string) {
    response := ErrorResponse{
        Success: false,
        Message: message,
        Error:   errorDetail,
    }

    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}

// Función auxiliar para obtener el usuario del contexto
func GetUserFromContext(r *http.Request) (*services.TokenClaims, bool) {
    user, ok := r.Context().Value(UserContextKey).(*services.TokenClaims)
    return user, ok
}