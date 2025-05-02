package main

// @title Go Auth Demo API
// @version 1.0
// @description This is a simple authentication API with MongoDB.
// @contact.name Felipe R. Tasoniero
// @contact.url https://github.com/frtasoniero
// @host localhost:5005
// @BasePath /

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/frtasoniero/lab/go-auth/pwh-auth/docs"
	"github.com/frtasoniero/lab/go-auth/pwh-auth/internal/auth"
	"github.com/frtasoniero/lab/go-auth/pwh-auth/internal/db"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate swag init -g cmd/server/main.go

func main() {
	log.Println("Starting authentication service...")

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found or failed to load, using system env vars")
	}

	// Get environment variables
	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	if port == "" {
		port = "5005" // Default fallback
	}

	if mongoURI == "" {
		log.Fatal("MONGO_URI not set")
	}

	// Connect to MongoDB
	if err := db.Connect(mongoURI); err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	log.Println("Connected to MongoDB")

	// Register HTTP handlers for the endpoints
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	// Start the HTTP server
	log.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

type Credentials struct {
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"secret123"`
}

type JSONResponse struct {
	Message string `json:"message" example:"Login successful"`
}

// registerHandler handles user registration
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body Credentials true "User credentials"
// @Success 200 {object} JSONResponse
// @Failure 400 {object} JSONResponse
// @Failure 409 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /register [post]
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	r.Body.Close()

	if creds.Username == "" || creds.Password == "" {
		respondJSON(w, http.StatusBadRequest, "Username and password required")
		return
	}

	creds.Username = strings.ToLower(strings.TrimSpace(creds.Username))

	// Check if user already exists
	existingUser, err := auth.GetUserByUsername(creds.Username)
	if err == nil && existingUser != nil {
		respondJSON(w, http.StatusConflict, "User already exists")
		return
	}

	// Optional: check if error is not "not found" — skip if you already handle it
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		log.Println("Database error during lookup:", err)
		respondJSON(w, http.StatusInternalServerError, "Server error")
		return
	}

	// Hash password and save
	hashedPassword, err := auth.HashPassword(creds.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		respondJSON(w, http.StatusInternalServerError, "Server error")
		return
	}

	newUser := auth.User{
		Username:     creds.Username,
		PasswordHash: hashedPassword,
	}

	if err := auth.SaveUser(newUser); err != nil {
		log.Println("Saving user error:", err)
		respondJSON(w, http.StatusInternalServerError, "Registration failed")
		return
	}

	log.Printf("Registered new user: %s", creds.Username)
	respondJSON(w, http.StatusOK, "Registration successful")
}

// loginHandler handles user login
// @Summary Login with credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body Credentials true "User credentials"
// @Success 200 {object} JSONResponse
// @Failure 400 {object} JSONResponse
// @Failure 401 {object} JSONResponse
// @Failure 500 {object} JSONResponse
// @Router /login [post]
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	r.Body.Close()

	if creds.Username == "" || creds.Password == "" {
		respondJSON(w, http.StatusBadRequest, "Username and password required")
		return
	}

	creds.Username = strings.ToLower(strings.TrimSpace(creds.Username))

	user, err := auth.GetUserByUsername(creds.Username)
	if err != nil {
		log.Println("User lookup error:", err)
		respondJSON(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if !auth.CheckPasswordHash(creds.Password, user.PasswordHash) {
		respondJSON(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	log.Printf("User logged in: %s", creds.Username)
	respondJSON(w, http.StatusOK, "Login successful")
}

func respondJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JSONResponse{Message: message})
}
