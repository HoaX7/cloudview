package main

//
// Responsible for handling global APIServer struct which gives
// access to controllers & database.
//
// Responsible for handling routes.
// This file was written with golang best practices in mind.
//

import (
	"cloudview/app/src/api/controllers"
	"cloudview/app/src/api/middleware/logger"
	router "cloudview/app/src/api/routes"
	"cloudview/app/src/cache"
	"cloudview/app/src/database"

	"cloudview/app/src/helpers"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
)

type APIServer struct {
	listenAddr  string
	controllers *controllers.Controller
	DB          *database.DB
}

/*
*
Store all required structs in a global struct to be able to access
them where necessary.

The reason for this is, We only need to establish DB connection once
and not every time a controller/model is called.
This is the best practices used in golang to connect to desired DB.
*/
func NewAPIServer(listenAddr string) (*APIServer, error) {
	db, err := database.NewDB()
	if err != nil {
		return nil, err
	}
	return &APIServer{
		listenAddr:  listenAddr,
		controllers: controllers.InitControllers(),
		DB:          db,
	}, nil
}

func (s *APIServer) Run() {
	corsWhitelist := helpers.GoDotEnvVariable("CORS_WHITELIST")
	routes := router.RegisterRoutes(s)
	routes.Use(logRequests)

	// cors
	c := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(corsWhitelist, ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowCredentials: true,
	})

	routes.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Route not found"))
	})
	handler := c.Handler(routes)
	// Attach the middleware to the main handler

	// Init cache
	cache.Init()

	// Start the server
	logger.Logger.Log("Starting server on port: " + s.listenAddr)
	logger.Logger.Error(http.ListenAndServe(":"+s.listenAddr, handler))
}

func (s *APIServer) GetControllers() *controllers.Controller {
	return s.controllers
}

func (s *APIServer) GetDB() *database.DB {
	return s.DB
}

type responseRecorder struct {
	http.ResponseWriter
	status int
}

// Recording the status code for outgoing responses to log info
// see `logRequests`
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Logging incoming - outgoing requests
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		logger.Logger.Log(r.Method, r.RequestURI, r.RemoteAddr, recorder.status, time.Since(start))
	})
}
