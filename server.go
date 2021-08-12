package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/kameshsampath/hybrid-cloud-frontend-api/docs"
	"github.com/kameshsampath/hybrid-cloud-frontend-api/pkg/data"
	"github.com/kameshsampath/hybrid-cloud-frontend-api/pkg/routes"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const ()

var (
	err            error
	dbFile         string
	router         *gin.Engine
	httpPort       = "8080"
	db             *sql.DB
	httpListenPort = "8080"
	pgHost         = "localhost"
	pgPort         = "5432"
	pgUser         = "demo"
	pgPassword     = "pa55Word!"
	pgDatabase     = "demodb"
)

// @title Hybrid Cloud Demo Front API
// @version 1.0
// @description The front API that builds message to be processed by the backend api that will be spread across the clouds

// @contact.name Kamesh Sampath
// @contact.email kamesh.sampath@solo.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1/api
// @query.collection.format multi
// @schemes http https
func main() {

	if h := os.Getenv("POSTGRES_HOST"); h != "" {
		pgHost = h
	}
	if p := os.Getenv("POSTGRES_PORT"); p != "" {
		pgPort = p
	}
	if h := os.Getenv("POSTGRES_USER"); h != "" {
		pgUser = h
	}
	if p := os.Getenv("POSTGRES_PASSWORD"); p != "" {
		pgPassword = p
	}
	if d := os.Getenv("POSTGRES_DB"); d != "" {
		pgDatabase = d
	}
	pgsqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDatabase)
	db, err = sql.Open("postgres", pgsqlInfo)
	if err != nil {
		log.Fatalf("Error opening DB %s, reason %s", dbFile, err)
	}
	//TODO Graceful shutdown
	defer db.Close()

	_, err = db.Exec(data.DDLTABLES)
	if err != nil {
		log.Fatalf("Error initializing DB: %s", err)
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
