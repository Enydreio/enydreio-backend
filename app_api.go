package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

func ServeWebsite(w http.ResponseWriter, r *http.Request) {
	website := "dist"
	fallbackFile := "index.html"
	requestedFile := filepath.Join(website, r.URL.Path)
	_, err := os.Stat(requestedFile)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(website, fallbackFile))
		return
	}

	http.ServeFile(w, r, requestedFile)
}
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var app Application

	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result := db.Create(&app)
	if result.Error != nil {
		http.Error(w, "Failed to create application", http.StatusInternalServerError)
		return
	}
}
func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var app Application
	result := db.Where("id = ?", id).Delete(&app)
	if result.Error != nil {
		http.Error(w, "Failed to delete application", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func ListApplications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var applications []Application
	result := db.Find(&applications)
	if result.Error != nil {
		http.Error(w, "Failed to retrieve applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(applications)
	if err != nil {
		http.Error(w, "Failed to encode applications", http.StatusInternalServerError)
		return
	}
}
func EditApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var app Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result := db.Save(&app)
	if result.Error != nil {
		http.Error(w, "Failed to update application", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(app)
}
