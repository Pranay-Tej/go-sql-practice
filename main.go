package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-sql-practice/config"
	"go-sql-practice/internal/database"
	"go-sql-practice/projects"
	"go-sql-practice/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	err = db.Ping()
	if err != nil {
		log.Fatal("error connecting to db")
	}

	dbQueries := database.New(db)

	apiConfig := config.ApiConfig{
		Db: dbQueries,
	}
	fmt.Printf("started server on port: %v\n", port)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handleHello)

	r.Post("/users", users.HandleCreateUser(&apiConfig))
	r.Get("/users", users.HandleGetAllUsers(&apiConfig))
	r.Get("/users/{id}", users.HandleGetUserById(&apiConfig))
	r.Post("/projects", projects.HandleCreateProject(&apiConfig))
	r.Get("/projects", projects.HandleGetAllProjects(&apiConfig))
	r.Get("/projects/{id}", projects.HandleGetProjectById(&apiConfig))
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}
