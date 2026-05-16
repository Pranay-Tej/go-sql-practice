package projects

import (
	"encoding/json"
	"log"
	"net/http"

	"go-sql-practice/config"
	"go-sql-practice/internal/database"
	"go-sql-practice/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func HandleCreateProject(apiConfig *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Input struct {
			Name   string    `json:"name"`
			UserId uuid.UUID `json:"user_id"` // TODO: get this from token
		}
		input := Input{}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			log.Printf("error decoding input: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		project, err := apiConfig.Db.CreateProject(r.Context(), database.CreateProjectParams{
			Name: input.Name,
			UserID: uuid.NullUUID{
				UUID:  input.UserId,
				Valid: true,
			},
		})

		if err != nil {
			log.Printf("error creating project: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		projectJson, err := json.Marshal(project)
		if err != nil {
			log.Printf("error encoding json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(projectJson)
	}
}

func HandleGetAllProjects(apiConfig *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := apiConfig.Db.GetProjects(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonProjects := make([]models.Project, 0)
		for _, project := range projects {
			jsonProjects = append(jsonProjects, models.Project{
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
}

func HandleGetProjectById(apiConfig *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		projectId, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		project, err := apiConfig.Db.GetProjectById(r.Context(), projectId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		projectJson, err := json.Marshal(project)
		if err != nil {
			log.Printf("error encoding json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(projectJson)
	}
}
