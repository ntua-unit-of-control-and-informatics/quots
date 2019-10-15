package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"quots/config"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	auth "quots/authentication"
	db "quots/database"
	httphandlers "quots/httphandlers"
	middleware "quots/middleware"
)

type IApp interface {
	GetConfig()
}

type App struct {
	config *config.Config
}

// func optionsHandler(h http.Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == "OPTIONS" {
// 			headers := w.Header()
// 			headers.Add("Access-Control-Allow-Origin", "*")
// 			headers.Add("Vary", "Origin")
// 			headers.Add("Vary", "Access-Control-Request-Method")
// 			headers.Add("Vary", "Access-Control-Request-Headers")
// 			headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
// 			headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
// 		} else {
// 			h.ServeHTTP(w, r)
// 		}
// 	}
// }

//docker run -it -p 8000:8000 --name quo --env MONGO_URL="mongodb://host.docker.internal:27017" quots

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token, Authorization, Client-id, Client-secret, Total, total")
	headers.Add("Access-Control-Allow-Methods", "GET, PUT, DELETE, POST,OPTIONS")
	json.NewEncoder(w)
	// return w
}

func main() {
	fmt.Println("Starting the application...")
	// conf := config.Init()
	// app := &App{config: conf}
	r := mux.NewRouter()
	db.NewDB()
	auth.NewAuth()
	r.Methods("OPTIONS").HandlerFunc(optionsHandler)

	r.HandleFunc("/users", middleware.AuthMiddleware(httphandlers.CreateUser)).Methods("POST")
	r.HandleFunc("/users", middleware.AuthMiddleware(httphandlers.GetUsers)).Queries("min", "{min}").Queries("max", "{max}").Queries("email", "{email}").Methods("GET")
	r.HandleFunc("/users/{id}", middleware.AuthMiddleware(httphandlers.GetUser)).Methods("GET")
	r.HandleFunc("/users/{id}", middleware.AuthMiddleware(httphandlers.DeleteUser)).Methods("DELETE")
	r.HandleFunc("/users/{id}", middleware.AuthMiddleware(httphandlers.UpdateEmailAndName)).Methods("PUT")
	// r.HandleFunc("/users/email", middleware.AuthMiddleware(httphandlers.GetUserByEmail)).Queries("email", "{email}").Methods("GET")
	r.HandleFunc("/users/credits", middleware.AuthMiddleware(httphandlers.UpdateCredits)).Methods("PUT")
	r.HandleFunc("/users/{id}/quots", middleware.AuthMiddleware(httphandlers.UserQuotes)).Queries("appid", "{appid}").Queries("usage", "{usage}").Queries("size", "{size}").Methods("GET")
	r.HandleFunc("/app/{id}", middleware.AuthMiddleware(httphandlers.GetApplication)).Methods("GET")
	r.HandleFunc("/app/{id}", middleware.AuthMiddleware(httphandlers.DeleteApplication)).Methods("DELETE")
	r.HandleFunc("/app/{id}/secret", middleware.AuthMiddleware(httphandlers.UpdateApplicationSecret)).Methods("PUT")
	// r.HandleFunc("/app", optionsHandler).Methods("OPTIONS")
	r.HandleFunc("/app", middleware.AuthMiddleware(httphandlers.CreateApplication)).Methods("POST")
	r.HandleFunc("/app", middleware.AuthMiddleware(httphandlers.UpdateApplication)).Methods("PUT")
	// r.HandleFunc("/app", optionsHandler).Queries("min", "{min}").Queries("max", "{max}").Methods("OPTIONS")
	r.HandleFunc("/app", middleware.AuthMiddleware(httphandlers.GetApplications)).Queries("min", "{min}").Queries("max", "{max}").Methods("GET")
	r.HandleFunc("/apikey", auth.AuthenicateHandler).Queries("username", "{username}").Queries("password", "{password}").Methods("GET")
	r.HandleFunc("/apikey/valid", auth.ValidateTokenHandler).Queries("apikey", "{apikey}").Methods("GET")
	r.HandleFunc("/apikey/refresh", middleware.AuthMiddleware(auth.RefreshTokenHandler)).Queries("apikey", "{apikey}").Methods("PUT")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/dist/quots-ui/")))
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"*"})
	exposedHeaders := handlers.ExposedHeaders([]string{"*"})
	log.Fatal(http.ListenAndServe(":8002", handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods, exposedHeaders, handlers.IgnoreOptions())(r)))
}
