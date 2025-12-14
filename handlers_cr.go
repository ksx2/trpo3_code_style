package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(request.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid email format"})
		return
	}

	if len(request.Password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password must be at least 8 characters"})
		return
	}

	if strings.TrimSpace(request.Name) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Name is required"})
		return
	}

	if request.Age < 18 || request.Age > 120 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Age must be between 18 and 120"})
		return
	}

	var existingID int
	err := h.db.QueryRow("SELECT id FROM users WHERE email = $1", request.Email).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "User with this email already exists"})
		return
	}

	hashedPassword := fmt.Sprintf("hashed_%s", request.Password)

	query := `INSERT INTO users (email, password_hash, name, age, created_at) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	var userID int
	err = h.db.QueryRow(query, request.Email, hashedPassword, request.Name, request.Age, time.Now()).Scan(&userID)
	if err != nil {
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create user"})
		return
	}

	response := map[string]interface{}{
		"id":        userID,
		"email":     request.Email,
		"name":      request.Name,
		"age":       request.Age,
		"created_at": time.Now().Format(time.RFC3339),
		"message":   "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User ID is required"})
		return
	}

	var id int
	var email, name string
	var age int
	var createdAt time.Time
	
	err := h.db.QueryRow(
		"SELECT id, email, name, age, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&id, &email, &name, &age, &createdAt)
	
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "User not found"})
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		return
	}

	response := map[string]interface{}{
		"id":         id,
		"email":      email,
		"name":       name,
		"age":        age,
		"created_at": createdAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
