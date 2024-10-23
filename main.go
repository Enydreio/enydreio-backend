package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
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

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			ScanApps()
			<-ticker.C
		}
	}()

	db.AutoMigrate(&Application{})
	http.HandleFunc("/", ServeWebsite)
	http.HandleFunc("/api/create-application", CreateApplication)
	http.HandleFunc("/api/delete-application", DeleteApplication)
	http.HandleFunc("/api/list-applications", ListApplications)
	http.HandleFunc("/api/edit-application", EditApplication)

	PrintLogo()

	err := http.ListenAndServe(":8080", nil)
	if err != nil || errDb != nil {
		return
	}

}
func PrintLogo() {
	fmt.Print(" _____                _          _       \n|  ___|              | |        (_)      \n| |__ _ __  _   _  __| |_ __ ___ _  ___  \n|  __| '_ \\| | | |/ _` | '__/ _ \\ |/ _ \\ \n| |__| | | | |_| | (_| | | |  __/ | (_) |\n\\____/_| |_|\\__, |\\__,_|_|  \\___|_|\\___/ \n             __/ |                       \n            |___/  \n")
}
