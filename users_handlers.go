package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
