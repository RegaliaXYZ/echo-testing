package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"bp-echo-test/internal/database"
	"bp-echo-test/internal/utils"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port       int
	db         database.Service
	gcp_client utils.GoogleServiceInterface
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	root_bucket := os.Getenv("ROOT_BUCKET")
	if root_bucket == "" {
		panic("ROOT_BUCKET not set")
	}
	NewServer := &Server{
		port:       port,
		db:         database.New(),
		gcp_client: utils.NewGoogleService(root_bucket),
	}

	auth_token := os.Getenv("AUTH_TOKEN")
	if auth_token == "" {
		panic("AUTH_TOKEN not set")
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(auth_token),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
