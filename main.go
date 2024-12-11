package main

import (
	"flag"
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

var dsn = fmt.Sprintf(os.Getenv("DATABASE_DATA"))
var db, errDb = gorm.Open(postgres.Open(dsn), &gorm.Config{})

func main() {
	PrintLogo()

	kubernetesFlag := flag.Bool("kubernetes", false, "Scan kubernetes apps")
	externFlag := flag.Bool("extern", false, "Kubernetes Cluster is external")
	dockerFlag := flag.Bool("docker", false, "Scan docker apps")
	intervalFlag := flag.Int("interval", 0, "Interval in minutes to wait between actions")

	flag.Parse()
	ticker := time.NewTicker(1 * time.Minute)
	if *intervalFlag > 0 {
		ticker = time.NewTicker(time.Duration(*intervalFlag) * time.Minute)
	} else {
		fmt.Println("No interval specified, using default 1 Minute")
	}
	go func() {
		for {
			if *kubernetesFlag {
				ScanKubeApps(*externFlag)
			}
			if *dockerFlag {
				ScanDockerApps()
			}
			<-ticker.C
		}
	}()

	db.AutoMigrate(&Application{})
	http.HandleFunc("/", ServeWebsite)
	http.HandleFunc("/api/create-application", CreateApplication)
	http.HandleFunc("/api/delete-application", DeleteApplication)
	http.HandleFunc("/api/list-applications", ListApplications)
	http.HandleFunc("/api/edit-application", EditApplication)

	err := http.ListenAndServe(":8080", nil)
	if err != nil || errDb != nil {
		return
	}

}
func PrintLogo() {
	fmt.Print(" _____                _          _       \n" +
		"|  ___|              | |        (_)      \n" +
		"| |__ _ __  _   _  __| |_ __ ___ _  ___  \n" +
		"|  __| '_ \\| | | |/ _` | '__/ _ \\ |/ _ \\ \n" +
		"| |__| | | | |_| | (_| | | |  __/ | (_) |\n" +
		"\\____/_| |_|\\__, |\\__,_|_|  \\___|_|\\___/ \n" +
		"             __/ |                       " +
		"\n            |___/  \n")
}
