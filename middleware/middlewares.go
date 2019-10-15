package middleware

import (
	"encoding/json"
	"net/http"
	auth "quots/authentication"
	db "quots/database"
	models "quots/models"

	"strings"
)

var IAppdb db.IApplicationDao
var dbAppImpl = &db.ApplicationDao{}

var IUsersdb db.IUsersDao
var dbUserImpl = &db.UsersDao{}

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiAll := r.Header.Get("Authorization")
		if apiAll == "" {
			err := models.ErrorReport{
				Message: "Not authorized",
				Status:  http.StatusUnauthorized,
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}
		apiKeyAr := strings.Split(apiAll, " ")
		authType := apiKeyAr[0]
		if authType == "Bearer" {
			apiKey := apiKeyAr[1]
			_, err := auth.ValidateToken(apiKey)
			if err != nil {
				err := models.ErrorReport{
					Message: err.Error(),
					Status:  http.StatusUnauthorized,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
		}
		if authType == "QUOTSAPP" {
			clientId := r.Header.Get("app-id")
			clientSecret := r.Header.Get("app-secret")
			app, err := dbAppImpl.GetApplicationBiId(clientId)
			if err != nil {
				err := models.ErrorReport{
					Message: err.Error(),
					Status:  http.StatusUnauthorized,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
			if app.AppSecret != clientSecret {
				err := models.ErrorReport{
					Message: "Wrong app secret",
					Status:  http.StatusUnauthorized,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
			if app.Enabled == false {
				err := models.ErrorReport{
					Message: "Application is disabled",
					Status:  http.StatusUnauthorized,
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
			// origin := r.Header.Get("Origin")
			// index := indexOfUrls(origin, app)
			// if index == -1 {
			// 	err := models.ErrorReport{
			// 		Message: "Not available origin for the Application",
			// 		Status:  http.StatusUnauthorized,
			// 	}
			// 	w.WriteHeader(http.StatusUnauthorized)
			// 	json.NewEncoder(w).Encode(err)
			// 	return
			// }
		}
		h.ServeHTTP(w, r)
	})
}

func indexOfUrls(element string, app models.Application) int {
	for k, v := range app.BaseURLS {
		if element == v {
			return k
		}
		if "*" == v {
			return 100
		}
	}
	return -1 //not found.
}
