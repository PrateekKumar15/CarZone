package routes

import (
	"github.com/gorilla/mux"
)

// setupAuthRoutes configures all authentication-related routes
func (r *Router) setupAuthRoutes(router *mux.Router) {
	// POST /auth/register - Register a new user account
	router.HandleFunc("/auth/register", r.AuthHandler.RegisterHandler).Methods("POST", "OPTIONS")

	// POST /auth/login - Authenticate user and receive access token
	router.HandleFunc("/auth/login", r.AuthHandler.LoginHandler).Methods("POST", "OPTIONS")

	// GET /auth/logout - Logout user (invalidate session)
	router.HandleFunc("/auth/logout", r.AuthHandler.LogoutHandler).Methods("GET", "OPTIONS")
}
