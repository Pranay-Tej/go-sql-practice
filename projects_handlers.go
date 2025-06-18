package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Pranay-Tej/go-sql-practice/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *ApiConfig) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Name   string    `json:"name"`
		UserId uuid.UUID `json:"user_id"` // TODO: get this from token
	}
	input := Input{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding input: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := apiConfig.db.CreateProject(r.Context(), database.CreateProjectParams{
		Name: input.Name,
		UserID: uuid.NullUUID{
			UUID:  input.UserId,
			Valid: true,
		},
	})

	if err != nil {
		log.Printf("error creating project: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	projectJson, err := json.Marshal(project)
	if err != nil {
		log.Printf("error encoding json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(projectJson)
}

func (apiConfig *ApiConfig) handleGetAllProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := apiConfig.db.GetProjects(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonProjects := make([]Project, 0)
	for _, project := range projects {
		jsonProjects = append(jsonProjects, Project{
			Id:        project.ID,
			CreatedAt: project.CreatedAt,
			UpdatedAt: project.UpdatedAt,
			Name:      project.Name,
			UserId:    project.UserID.UUID,
		})
	}
	data, err := json.Marshal(jsonProjects)
	if err != nil {
		log.Printf("error encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
