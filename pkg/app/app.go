package app

import (
	"code-exec/pkg"
	"code-exec/pkg/dependencies/anchor"
	"code-exec/pkg/dependencies/rust"
	"code-exec/pkg/dependencies/typescript"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var _ = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type App struct {
	env    string
	engine *gin.Engine
	deps   pkg.Dependencies
}

func NewApp(
	env string,
	repo pkg.Repository,
	rpcEngine pkg.RpcEngine,
) *App {
	engine := gin.New()
	engine.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	gin.SetMode(gin.DebugMode)

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Content-Disposition", "Sec-Websocket-Protocol"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	engine.Use(cors.New(config))

	return &App{
		env:    env,
		engine: engine,
		deps: pkg.Dependencies{
			Repo:          repo,
			RpcEngine:     rpcEngine,
			TsRuntime:     typescript.NewRuntime(),
			RustRuntime:   rust.NewRuntime(),
			AnchorRuntime: anchor.NewRuntime(),
		},
	}
}

func (a *App) Run() {
	a.AttachStandardRoutes()
	a.AttachCodeExecRoutes()

	// Server configurations for access across go routines
	server := &http.Server{
		Addr:    getPort(),
		Handler: a.engine,
	}

	// Execute a go routine so we can start the server and wait for closing signal
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections")
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Create a context to deliver possible errors
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 60*time.Second)
	defer shutdownRelease()

	// Shutdown the server gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}

func (a *App) AttachStandardRoutes() {
	a.engine.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		log.Printf("Environment variable PORT=\"%s\"", port)
		return ":" + port
	}
	log.Println("Environment variable PORT is undefined. Using port :8080 by default")
	return ":8081"
}
