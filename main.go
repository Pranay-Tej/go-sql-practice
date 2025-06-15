package main

import (
	"database/sql"
	"encoding/json"
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

func main() {
	godotenv.Load()

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", handleHello)

	r.Post("/users", apiConfig.handleCreateUser)
	r.Get("/users", apiConfig.handleGetAllUsers)

	http.ListenAndServe(":3000", r)
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func (apiConfig *ApiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Username string `json:"username"`
	}
	input := Input{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding input: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := apiConfig.db.CreateUser(r.Context(), input.Username)

	if err != nil {
		log.Printf("error creating user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("error encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}

func (apiConfig *ApiConfig) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiConfig.db.GetUsers(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonUsers := make([]User, 0)

	for _, user := range users {
		jsonUsers = append(jsonUsers, User{
			Id:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
		})
	}

	data, err := json.Marshal(jsonUsers)
	if err != nil {
		log.Printf("error encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
