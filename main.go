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
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Logo        string `json:"logo"`
}

var db *gorm.DB

func main() {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Vienna", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, errDb := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Application{})
	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/", fs)
	http.HandleFunc("/api/create-application", createApplication)

	err := http.ListenAndServe(":8080", nil)
	if err != nil || errDb != nil {
		return
	}
}
func createApplication(w http.ResponseWriter, r *http.Request) {
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
