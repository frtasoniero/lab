package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/frtasoniero/lab/go-auth/pwt-auth/internal/auth"
	"github.com/frtasoniero/lab/go-auth/pwt-auth/internal/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	log.Println("Starting authentication service...")

	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found or failed to load, using system env vars")
	}

	// Get environment variables
	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default fallback
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

	// Start the HTTP server
	log.Printf("Server listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JSONResponse struct {
	Message string `json:"message"`
}

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

	hashedPassword, err := auth.HashPassword(creds.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		respondJSON(w, http.StatusInternalServerError, "Server error")
		return
	}

	user := auth.User{
		Username:     creds.Username,
		PasswordHash: hashedPassword,
	}

	err = auth.SaveUser(user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			respondJSON(w, http.StatusConflict, "User already exists")
		} else {
			log.Println("Error saving user:", err)
			respondJSON(w, http.StatusInternalServerError, "Registration failed")
		}
		return
	}

	log.Printf("New user registered: %s", creds.Username)
	respondJSON(w, http.StatusOK, "Registration successful")
}

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
