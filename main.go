package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	auth "euclia.xyz/accounts-api/authentication"
	db "euclia.xyz/accounts-api/database"
	email "euclia.xyz/accounts-api/emails"
	httphandlers "euclia.xyz/accounts-api/httphandlers"
	middleware "euclia.xyz/accounts-api/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type IApp interface {
	// GetConfig()
}

type App struct {
	// config *config.Config
}

var mid middleware.IAuth = &middleware.AuthImplementation{}

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

	auth.Init()
	// conf := config.Init()
	// app := &App{config: conf}
	r := mux.NewRouter()
	db.NewDB()
	httphandlers.InitQuots()
	email.Init()
	r.Methods("OPTIONS").HandlerFunc(optionsHandler)
	r.HandleFunc("/users", mid.AuthMiddleware(httphandlers.GetUser)).Queries("min", "{min}").Queries("max", "{max}").Queries("email", "{email}").Queries("id", "{id}").Methods("GET")
	r.HandleFunc("/users", mid.AuthMiddleware(httphandlers.UpdateUser)).Methods("PUT")
	r.HandleFunc("/users/{id}/credits", mid.AuthMiddleware(httphandlers.UpdateUserCredits)).Methods("PUT")
	r.HandleFunc("/users/{id}/organizations/{orgid}", mid.AuthMiddleware(httphandlers.UpdateUserOrganizations)).Methods("PUT")
	r.HandleFunc("/organizations", mid.AuthMiddleware(httphandlers.CreateOrganization)).Methods("POST")
	r.HandleFunc("/organizations", mid.AuthMiddleware(httphandlers.UpdateOrganization)).Methods("PUT")
	r.HandleFunc("/organizations", mid.AuthMiddleware(httphandlers.GetOrganizations)).Methods("GET")
	r.HandleFunc("/organizations/{id}", mid.AuthMiddleware(httphandlers.GetOrganizationById)).Methods("GET")
	r.HandleFunc("/organizations/{id}", mid.AuthMiddleware(httphandlers.DeleteOrganization)).Methods("DELETE")
	r.HandleFunc("/notification", mid.AuthMiddleware(httphandlers.CreateInvitationNotification)).Methods("POST")
	r.HandleFunc("/request", mid.AuthMiddleware(httphandlers.CreateRequest)).Methods("POST")
	r.HandleFunc("/invitation", mid.AuthMiddleware(httphandlers.CreateInvitation)).Methods("POST")
	r.HandleFunc("/invitation", mid.AuthMiddleware(httphandlers.UpdateInvitation)).Queries("answer", "{answer}").Methods("PUT")
	r.HandleFunc("/invitation", mid.AuthMiddleware(httphandlers.GetInvitation)).Queries("min", "{min}").Queries("max", "{max}").Queries("email", "{email}").Queries("viewed", "{viewed}").Methods("GET")
	// r.HandleFunc("/users", httphandlers.GetUser).Queries("id", "{id}").Queries("min", "{min}").Queries("max", "{max}").Queries("email", "{email}").Methods("GET")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"*"})
	exposedHeaders := handlers.ExposedHeaders([]string{"*"})
	log.Fatal(http.ListenAndServe(":8006", handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods, exposedHeaders, handlers.IgnoreOptions())(loggedRouter)))
}
