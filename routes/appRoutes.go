package routes

import (
	"hells/controllers"
	"hells/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	// Authentication Routes
	router.HandleFunc("/register", controllers.Register).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/reset-password", controllers.ResetPassword).Methods("POST")

	// User Routes
	userRoutes := router.PathPrefix("/users").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware)
	userRoutes.HandleFunc("", controllers.ListUsers).Methods("GET")
	userRoutes.HandleFunc("/{id}", controllers.GetUser).Methods("GET")
	userRoutes.HandleFunc("/{id}", middleware.RBACMiddleware("Admin")(controllers.UpdateUser)).Methods("PUT")

	// Post Routes
	// postRoutes := router.PathPrefix("/posts").Subrouter()
	// postRoutes.Use(middleware.AuthMiddleware)
	// postRoutes.HandleFunc("", middleware.RBACMiddleware("Editor")(controllers.CreatePost)).Methods("POST")
	// postRoutes.HandleFunc("", controllers.ListPosts).Methods("GET")
	// postRoutes.HandleFunc("/{id}", controllers.GetPost).Methods("GET")
	// postRoutes.HandleFunc("/{id}", middleware.RBACMiddleware("Editor")(controllers.UpdatePost)).Methods("PUT")
}
