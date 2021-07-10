package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/kameshsampath/hybrid-cloud-frontend-api/docs"
	"github.com/kameshsampath/hybrid-cloud-frontend-api/pkg/routes"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	//DDLTABLES  creates the workers and cloud workers table
	DDLTABLES = `
DROP TABLE IF EXISTS cloud_workers;
CREATE TABLE IF NOT EXISTS cloud_workers (
requestId text PRIMARY KEY NOT NULL,
workerId TEXT NOT NULL,
cloud TEXT NOT NULL,
requestsProcessed INTEGER TEXT DEFAULT 1,
response TEXT TEXT,
timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
`
)

var (
	dbDir    string
	db       *sql.DB
	err      error
	dbFile   string
	router   *gin.Engine
	httpPort = "8080"
)

func init() {
	if dbDir = os.Getenv("HYBRID_CLOUD_DB_DIR"); dbDir == "" {
		homedir, _ := os.UserHomeDir()
		dbDir = filepath.Join(homedir, ".hybrid-cloud-app")
	}
	if _, err := os.Stat(dbDir); err != nil {
		if err = os.Mkdir(dbDir, os.ModeDir); err != nil {
			panic(fmt.Sprintf("Error creating DB Dir %s", dbDir))
		}
	}
}

// @title Hybrid Cloud Demo Front API
// @version 1.0
// @description The front API that builds message to be processed by the backend api that will be spread across the clouds

// @contact.name Kamesh Sampath
// @contact.email kamesh.sampath@solo.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1/api
// @query.collection.format multi
// @schemes http https
func main() {

	if err == nil {
		dbFile = filepath.Join(dbDir, "app.db")
		db, err = sql.Open("sqlite3", dbFile)
		if err != nil {
			log.Fatalf("Error opening DB %s, reason %s", dbFile, err)
		}
		//TODO Graceful shutdown
		defer db.Close()

		_, err = db.Exec(DDLTABLES)
		if err != nil {
			log.Fatalf("Error initializing DB: %s", err)
		}
	}

	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}
	router = gin.Default()
	// this is liberal CORS settings only for demo
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	addRoutes()

	if hPort := os.Getenv("HTTP_LISTEN_PORT"); hPort != "" {
		httpPort = hPort
	}
	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%s", httpPort),
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listent: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down server ...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutdown", err)
	}

	log.Println("Server Exiting")
}

func addRoutes() {
	endpoints := routes.NewEndpoints()
	endpoints.DBConn = db

	v1 := router.Group("/v1/api")
	{
		//Health Endpoints accessible via /v1/api/health
		health := v1.Group("/health")
		{
			health.GET("/live", endpoints.Live)
			health.GET("/ready", endpoints.Ready)
		}
		//frontend endpoints
		frontend := v1.Group("/")
		{
			frontend.POST("/send-request", endpoints.SendRequest)
		}
		//workers endpoints
		workers := v1.Group("/workers")
		{
			workers.GET("/all", endpoints.Responses)
			workers.GET("/cloud", endpoints.CloudWorkerRequests)
		}
	}

	// the default path to get swagger json is :8080/swagger/docs.json
	// TODO enable/disable based on ENV variable
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
