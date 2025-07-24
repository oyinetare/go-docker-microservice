package api

import (
	"encoding/json"
	"net/http"

	"github.com/oyinetare/go-docker-microservice/repository"
)

type UserAPI struct {
	repo *repository.Repository
}

func NewUserAPI(repo *repository.Repository) *UserAPI {
	return &UserAPI{repo: repo}
}

func (api *UserAPI) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.repo.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (api *UserAPI) SearchUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "When searching for a user, the email must be specified, e.g: '/search?email=homer@thesimpsons.com'", http.StatusBadRequest)
		return
	}

	user, err := api.repo.GetUserByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
