package authentication

import (
	"encoding/json"
	"net/http"

	models "quots/models"
)

type IAuthenticate interface {
	Authenticate(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
}

type AuthenicateImpl struct {
}

// var authImpl = &AuthenicateImpl

func AuthenicateHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	jwt, err := AuthenticateAdmin(username, password)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusForbidden,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
	} else {
		auth := AuthResponse{
			ApiKey: jwt,
		}
		json.NewEncoder(w).Encode(auth)
	}

}

func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	api_key := r.FormValue("apikey")
	if api_key != "null" {
		valid, err := ValidateToken(api_key)
		if err != nil || api_key == "" {
			err := models.ErrorReport{
				Message: err.Error(),
				Status:  http.StatusForbidden,
			}
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(err)
		} else {
			validation := ValidationResponse{
				ApiKey: api_key,
				Valid:  valid,
			}
			json.NewEncoder(w).Encode(validation)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(api_key)
	}
}

func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	api_key := r.FormValue("apikey")
	token, err := RenewToken(api_key)
	if err != nil {
		err := models.ErrorReport{
			Message: err.Error(),
			Status:  http.StatusForbidden,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
	} else {
		refreshed := RefreshedToken{
			ApiKey: token,
		}
		json.NewEncoder(w).Encode(refreshed)
	}
}
