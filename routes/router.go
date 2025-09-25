package routes

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	authHandler "github.com/PrateekKumar15/CarZone/handler/auth"
	bookingHandler "github.com/PrateekKumar15/CarZone/handler/booking"
	carHandler "github.com/PrateekKumar15/CarZone/handler/car"
	"github.com/PrateekKumar15/CarZone/middleware"
)

// Router holds all the handler dependencies
type Router struct {
	AuthHandler    *authHandler.AuthHandler
	CarHandler     *carHandler.CarHandler
	BookingHandler *bookingHandler.BookingHandler
}

// NewRouter creates a new router instance with handler dependencies
func NewRouter(authHandler *authHandler.AuthHandler, carHandler *carHandler.CarHandler, bookingHandler *bookingHandler.BookingHandler) *Router {
	return &Router{
		AuthHandler:    authHandler,
		CarHandler:     carHandler,
		BookingHandler: bookingHandler,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add OpenTelemetry middleware for tracing
	router.Use(otelmux.Middleware("CarZone"))

	// Setup public routes (no authentication required)
	r.setupPublicRoutes(router)

	// Setup protected routes (authentication required)
	r.setupProtectedRoutes(router)

	// Setup monitoring routes
	r.setupMonitoringRoutes(router)

	return router
}

// setupPublicRoutes configures routes that don't require authentication
func (r *Router) setupPublicRoutes(router *mux.Router) {
	// Create a subrouter for public routes
	public := router.PathPrefix("/").Subrouter()

	// Authentication routes
	r.setupAuthRoutes(public)
}

// setupProtectedRoutes configures routes that require authentication
func (r *Router) setupProtectedRoutes(router *mux.Router) {
	// Create a subrouter for protected routes
	protected := router.PathPrefix("/").Subrouter()

	// Apply authentication middleware to all protected routes
	protected.Use(middleware.AuthMiddleware)
	protected.Use(middleware.MetricMiddleware)

	// Setup resource-specific routes
	r.setupCarRoutes(protected)
	r.setupBookingRoutes(protected)
}

// setupMonitoringRoutes configures monitoring and metrics routes
func (r *Router) setupMonitoringRoutes(router *mux.Router) {
	// Prometheus metrics endpoint (usually public for monitoring systems)
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
}
