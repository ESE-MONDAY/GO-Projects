package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 1. Our Data Model (The Movie Struct)
type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

// 2. Our In-Memory Database (A Slice of Movie Structs)
var movies = []Movie{
	{ID: 1, Title: "Inception", Director: "Christopher Nolan", Year: 2010},
	{ID: 2, Title: "The Matrix", Director: "Lana Wachowski", Year: 1999},
}

// Track the next ID to allocate when creating new items
var nextID = 3

// 3. Prometheus Metric Definition
var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"path", "method", "status"},
)

func init() {
	prometheus.MustRegister(httpDuration)
}

// 4. Monitoring & Timing Middleware
func monitorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()

		// A simple custom response writer wrapper to capture the status code
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start).Seconds()

		// Extract clean route patterns (e.g., matching "/movies/{id}" instead of raw IDs)
		routePattern, _ := mux.CurrentRoute(r).GetPathTemplate()

		// Record metrics to Prometheus
		httpDuration.WithLabelValues(routePattern, r.Method, strconv.Itoa(wrappedWriter.statusCode)).Observe(duration)

		// Output Structured JSON Logs
		slog.Info("http_request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrappedWriter.statusCode),
			slog.Duration("latency", time.Since(start)),
		)
	})
}

// Helper struct to trap status codes for metrics
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	r := mux.NewRouter()
	r.Use(monitorMiddleware)

	// --- METRICS ROUTE ---
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	// --- CRUD ROUTES ---
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST") // Fixed the typo here!
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	slog.Info("Movie CRUD service booting up on port :8080...")

	// Caught the execution error explicitly so it screams if port 8080 is blocked
	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		slog.Error("Server failed to bind to port", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

// --- HANDLER FUNCTIONS ---

// READ ALL
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// READ SINGLE
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, item := range movies {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

// CREATE
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = nextID
	nextID++

	movies = append(movies, movie)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// UPDATE
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedMovie Movie
	_ = json.NewDecoder(r.Body).Decode(&updatedMovie)

	for i, item := range movies {
		if item.ID == id {
			updatedMovie.ID = item.ID
			movies[i] = updatedMovie
			json.NewEncoder(w).Encode(updatedMovie)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

// DELETE
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range movies {
		if item.ID == id {
			movies = append(movies[:i], movies[i+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}
