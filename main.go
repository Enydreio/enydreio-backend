package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type Application struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Logo        string `json:"logo"`
}

var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Vienna", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
var db, errDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})

func main() {

	db.AutoMigrate(&Application{})
	http.HandleFunc("^(?!\\/api).*\n", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "dist/index.html")
	})
	http.HandleFunc("/api/create-application", CreateApplication)
	http.HandleFunc("/api/delete-application", DeleteApplication)
	http.HandleFunc("/api/list-applications", ListApplications)
	http.HandleFunc("/api/edit-application", EditApplication)

	err := http.ListenAndServe(":8080", nil)
	if err != nil || errDb != nil {
		return
	}
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
