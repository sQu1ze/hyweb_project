package router

import (
	_ "hyweb-api/internal/config"
	"hyweb-api/internal/handler"
	"hyweb-api/internal/middleware"
	"hyweb-api/internal/repository"
	"hyweb-api/internal/service"

	"github.com/gin-gonic/gin"

	_ "hyweb-api/docs/swagger" // swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize repositories and services
	userRepo := repository.NewUserRepository(repository.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Health Check
	r.GET("/api/v1/health", userHandler.HealthCheck)

	// Swagger API Docs
	// Use gin-swagger middleware to serve the API docs
	// @BasePath /api/v1
	// @title Gin-Gonic API
	// @version 1.0
	// @description This is a sample server for a Gin-Gonic API.
	// @contact.name API Support
	// @contact.url http://www.swagger.io/support
	// @contact.email support@swagger.io
	// @license.name Apache 2.0
	// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
	// @host localhost:8080
	// @schemes http
	// @securityDefinitions.apikey BearerAuth
	// @in header
	// @name Authorization
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)

			auth := users.Use(middleware.JWTAuthMiddleware())
			{
				auth.PUT("/changePassword", userHandler.ChangePassword)
			}
		}
	}

	return r
}
