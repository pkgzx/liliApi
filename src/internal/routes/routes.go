package routes

import (
	"net/http"

	"github.com/pkgzx/liliApi/src/internal/handlers"
	"github.com/pkgzx/liliApi/src/internal/middleware"
)

type Router struct {
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(authMiddleware *middleware.AuthMiddleware) *Router {
	return &Router{
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes(
	userHandler *handlers.UserHandler,
	// productHandler *handlers.ProductHandler,
	// orderHandler *handlers.OrderHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Configurar rutas de autenticación
	r.setupAuthRoutes(mux, userHandler)

	// Configurar otras rutas (cuando estén listas)
	// r.setupProductRoutes(mux, productHandler)
	// r.setupOrderRoutes(mux, orderHandler)

	return mux
}

// Rutas de autenticación
func (r *Router) setupAuthRoutes(mux *http.ServeMux, userHandler *handlers.UserHandler) {

	// Rutas públicas
	mux.HandleFunc("/api/auth/signin", userHandler.Login)
	mux.HandleFunc("/api/auth/signup", userHandler.CreateUser)
	mux.HandleFunc("/api/auth/refresh", userHandler.RefreshToken)

	// Rutas protegidas
	mux.HandleFunc("/api/auth/profile", r.authMiddleware.RequireAuth(userHandler.GetProfile))
}

// Rutas de productos (para cuando implementes el handler)
// func (r *Router) setupProductRoutes(mux *http.ServeMux, productHandler *handlers.ProductHandler) {
// 	// Rutas públicas de productos
// 	mux.HandleFunc("/api/products", productHandler.HandleProducts)
// 	mux.HandleFunc("/api/products/", productHandler.HandleProductByID)
// 	mux.HandleFunc("/api/categories", productHandler.HandleCategories)
//
// 	// Rutas protegidas de productos (solo para administradores)
// 	mux.HandleFunc("/api/admin/products", r.authMiddleware.RequireAuth(productHandler.AdminHandleProducts))
// 	mux.HandleFunc("/api/admin/products/", r.authMiddleware.RequireAuth(productHandler.AdminHandleProductByID))
// }

// Rutas de pedidos (para cuando implementes el handler)
// func (r *Router) setupOrderRoutes(mux *http.ServeMux, orderHandler *handlers.OrderHandler) {
// 	// Todas las rutas de pedidos requieren autenticación
// 	mux.HandleFunc("/api/orders", r.authMiddleware.RequireAuth(orderHandler.HandleOrders))
// 	mux.HandleFunc("/api/orders/", r.authMiddleware.RequireAuth(orderHandler.HandleOrderByID))
// 	mux.HandleFunc("/api/orders/status/", r.authMiddleware.RequireAuth(orderHandler.HandleOrdersByStatus))
// }
