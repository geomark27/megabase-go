package routes

import (
	"megabaseGo/internal/app/handlers"
	"megabaseGo/internal/app/middleware"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup configura todas las rutas de la aplicación
func Setup() *gin.Engine {
	// Crear router con configuración por defecto
	router := gin.Default()

	// ---- INICIO DEL AJUSTE ----

	// Configurar CORS de forma específica para permitir credenciales (cookies)
	config := cors.DefaultConfig()
	// 1. Permitir explícitamente el origen de tu frontend
	config.AllowOrigins = []string{os.Getenv("FRONT_URL")}
	// 2. Permitir que el navegador envíe y reciba cookies
	config.AllowCredentials = true 
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// ---- FIN DEL AJUSTE ----

	// Middleware de logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "MegabaseGo API is running",
			"version": "1.0.0",
		})
	})

	// Inicializar handlers y middleware
	userHandler := handlers.NewUserHandler()
	roleHandler := handlers.NewRoleHandler()
	authHandler := handlers.NewAuthHandler()
	authMiddleware := middleware.NewAuthMiddleware()

	// Grupo de rutas API v1
	v1 := router.Group("/api/v1")
	{
		// Rutas públicas de autenticación
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
		}

		// Rutas protegidas (requieren autenticación)
		protected := v1.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		{
			// Profile endpoints
			protected.GET("/profile", authHandler.GetProfile)
			protected.POST("/change-password", authHandler.ChangePassword)
			protected.GET("/check-auth", authHandler.CheckAuth)

			// Rutas para roles (requiere autenticación)
			roles := protected.Group("/roles")
			{
				roles.POST("", roleHandler.CreateRole)
				roles.GET("", roleHandler.GetRoles)
				roles.GET("/:id", roleHandler.GetRole)
				roles.PUT("/:id", roleHandler.UpdateRole)
				roles.DELETE("/:id", roleHandler.DeleteRole)
			}

			// Rutas para usuarios (requiere autenticación)
			users := protected.Group("/users")
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.GET("/check-username", userHandler.CheckUsernameAvailability)
				users.GET("/check-email", userHandler.CheckEmailAvailability)
			}

			// Grupo de rutas para ciudadanos
			citizens := protected.Group("/citizens")
			{
				citizenHandler := handlers.NewCitizenHandler()
				
				// CRUD básico
				citizens.GET("", citizenHandler.GetAllCitizens)
				citizens.POST("", citizenHandler.CreateCitizen)
				citizens.GET("/:id", citizenHandler.GetCitizenByID)
				citizens.PUT("/:id", citizenHandler.UpdateCitizen)
				citizens.DELETE("/:id", citizenHandler.DeleteCitizen)
				
				// Búsquedas específicas
				citizens.GET("/email/:email", citizenHandler.GetCitizenByEmail)
				citizens.GET("/identification/:numero", citizenHandler.GetCitizenByIdentification)
				citizens.GET("/razon-social/:razon", citizenHandler.GetCitizenByRazonSocial)
				
				// Verificaciones de disponibilidad
				citizens.GET("/check/identification/:numero", citizenHandler.CheckIdentificationAvailability)
				citizens.GET("/check/email/:email", citizenHandler.CheckEmailAvailability)
				citizens.GET("/check/razon-social/:razon", citizenHandler.CheckRazonSocialAvailability)
			}
		}

		// Ruta de información de la API (pública)
		v1.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"api_version": "1.0.0",
				"project":     "MegabaseGo",
				"architecture": gin.H{
					"pattern":     "Handler + Service (Simplified) + JWT Auth",
					"description": "Clean architecture with JWT authentication for rapid development",
				},
				"endpoints": gin.H{
					"auth": gin.H{
						"login":     "POST /api/v1/auth/login",
						"register":  "POST /api/v1/auth/register",
						"refresh":   "POST /api/v1/auth/refresh",
						"logout":    "POST /api/v1/auth/logout",
						"profile":   "GET /api/v1/profile (protected)",
						"check":     "GET /api/v1/check-auth (protected)",
						"password":  "POST /api/v1/change-password (protected)",
					},
					"roles": gin.H{
						"create": "POST /api/v1/roles (protected)",
						"list":   "GET /api/v1/roles (protected)",
						"get":    "GET /api/v1/roles/:id (protected)",
						"update": "PUT /api/v1/roles/:id (protected)",
						"delete": "DELETE /api/v1/roles/:id (protected)",
					},
					"users": gin.H{
						"create": "POST /api/v1/users (protected)",
						"list":   "GET /api/v1/users (protected)",
						"get":    "GET /api/v1/users/:id (protected)",
						"update": "PUT /api/v1/users/:id (protected)",
						"delete": "DELETE /api/v1/users/:id (protected)",
					},
				},
				"authentication": gin.H{
					"type":   "JWT Bearer Token",
					"header": "Authorization: Bearer <token>",
					"note":   "Include access token in Authorization header for protected routes",
				},
				"query_params": gin.H{
					"roles": gin.H{
						"include_inactive": "bool - Include inactive roles",
					},
					"users": gin.H{
						"include_inactive": "bool - Include inactive users",
						"role_id":          "int - Filter by role ID",
					},
				},
			})
		})
	}

	return router
}