package main

import (
	"fmt"
	"net/http"
	"time"

	"logs_server/communication"
	"logs_server/database"
	"logs_server/handlers"
	"logs_server/models"

	"github.com/gorilla/mux"
)

var StartTime time.Time

func main() {
	PORT := ":8081"

	database.Connection()
	database.DB.AutoMigrate(models.Log{})

	router := mux.NewRouter()
	functions_EndPoints(router.PathPrefix("/logs").Subrouter())

	go communication.ListenNotifications()

	StartTime = time.Now()
	fmt.Println("----------INIT LOG SERVER " + PORT + " at: " + StartTime.Format("15:04:05") + "----------")
	http.ListenAndServe(PORT, router)
}

func functions_EndPoints(router *mux.Router) {
	
	router.HandleFunc("/health", handlers.CheckHealth_handler).Methods("GET")
	router.HandleFunc("/health/live", handlers.CheckLive_handler).Methods("GET")
	router.HandleFunc("/health/ready", handlers.CheckReady_handler).Methods("GET")

	router.HandleFunc("/all", handlers.GetLogs_handler).Methods("GET")
	router.HandleFunc("/{appName}", handlers.GetLogsByApp_handler).Methods("GET")
	router.HandleFunc("/", handlers.AddLogs_handler).Methods("POST")
}
