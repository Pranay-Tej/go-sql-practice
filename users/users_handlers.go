package users

import (
	"encoding/json"
	"log"
	"net/http"

	"go-sql-practice/config"
	"go-sql-practice/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandleCreateUser(apiConfig *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Input struct {
			Username string `json:"username"`
		}
		input := Input{}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Printf("error decoding input: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := apiConfig.Db.CreateUser(r.Context(), input.Username)

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
}

func HandleGetAllUsers(apiConfig *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := apiConfig.Db.GetUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonUsers := make([]models.User, 0)

		for _, user := range users {
			jsonUsers = append(jsonUsers, models.User{
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
}

func HandleGetUserById(apiConfig *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		userId, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := apiConfig.Db.GetUserById(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
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
}
