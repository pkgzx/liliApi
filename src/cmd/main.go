package main

import (
	"log"
	"net/http"

	database "github.com/pkgzx/liliApi/src/internal/db"
	"github.com/pkgzx/liliApi/src/internal/handlers"
	"github.com/pkgzx/liliApi/src/internal/middleware"
	"github.com/pkgzx/liliApi/src/internal/routes"
	"github.com/pkgzx/liliApi/src/internal/services"
	"github.com/pkgzx/liliApi/src/pkg/config"
	"github.com/pkgzx/liliApi/src/pkg/repository"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Conectar a la base de datos
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inicializar repositorios
	userRepo := repository.NewUserRepository(db.DB)

	// Inicializar servicios
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userService, cfg.JWT.Secret)

	// Inicializar middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Inicializar handlers
	userHandler := handlers.NewUserHandler(userService, authService)

	// Configurar rutas
	router := routes.NewRouter(authMiddleware)
	mux := router.SetupRoutes(userHandler)

	// Servidor
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: mux,
	}

	log.Printf("Server starting on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
