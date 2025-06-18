package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Pranay-Tej/go-sql-practice/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ApiConfig struct {
	db *database.Queries
}

type User struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

type Project struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	UserId    uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL env not set")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("error connecting to db")
	}

	dbQueries := database.New(db)

	apiConfig := ApiConfig{
		db: dbQueries,
	}
	fmt.Printf("started server on port: %v\n", port)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handleHello)

	r.Post("/users", apiConfig.handleCreateUser)
	r.Get("/users", apiConfig.handleGetAllUsers)
	r.Get("/users/{id}", apiConfig.handleGetUserById)
	r.Post("/projects", apiConfig.handleCreateProject)
	r.Get("/projects", apiConfig.handleGetAllProjects)
	r.Get("/projects/{id}", apiConfig.handleGetProjectById)
	http.ListenAndServe(":"+port, r)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
