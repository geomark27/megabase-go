package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"megabaseGo/internal/config"
	"megabaseGo/internal/database"
	"megabaseGo/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🚀 Iniciando MegabaseGo API Server...")
	log.Println("⚡ Arquitectura: Handler + Service (Simplificada)")

	// 1. Cargar configuración
	cfg := config.LoadConfig()
	log.Printf("📊 Configuración cargada - Puerto: %s", cfg.ServerPort)

	// 2. Conectar a la base de datos
	_, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}
	defer database.CloseDB()
	log.Println("✅ Conexión a la base de datos establecida")

	// 3. Configurar Gin para producción si es necesario
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("🏭 Modo de producción activado")
	} else {
		log.Println("🔧 Modo de desarrollo activado")
	}

	// 4. Inicializar rutas
	router := routes.Setup()
	log.Println("🛣️  Rutas configuradas")

	// 5. Configurar servidor HTTP
	server := &http.Server{
		Addr:           ":" + cfg.ServerPort,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 6. Canal para manejar señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 7. Iniciar servidor en goroutine
	go func() {
		log.Printf("🌐 Servidor iniciado en http://localhost:%s", cfg.ServerPort)
		log.Printf("📚 Info de la API: http://localhost:%s/api/v1/info", cfg.ServerPort)
		log.Printf("❤️  Health check: http://localhost:%s/health", cfg.ServerPort)
		log.Println("")
		log.Println("📋 Endpoints disponibles:")
		log.Printf("   🎭 Roles:    http://localhost:%s/api/v1/roles", cfg.ServerPort)
		log.Printf("   👥 Usuarios: http://localhost:%s/api/v1/users", cfg.ServerPort)
		log.Println("")
		log.Println("✋ Presiona Ctrl+C para detener el servidor")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error iniciando el servidor: %v", err)
		}
	}()

	// 8. Esperar señal de terminación
	<-quit
	log.Println("🛑 Recibida señal de terminación, cerrando servidor...")

	// 9. Shutdown graceful del servidor
	log.Println("🔄 Cerrando conexiones activas...")
	
	// Da tiempo para que las conexiones activas se completen
	time.Sleep(1 * time.Second)
	
	log.Println("✅ Servidor cerrado correctamente")
	log.Println("👋 ¡Hasta luego!")
}