package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"bp-echo-test/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		db:   database.New(),
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
