package main

import (
	"fmt"
	"net/http"
	"time"
	"user_server/database"
	"user_server/handlers"
	"user_server/models"

	"github.com/gorilla/mux"
)

var StartTime time.Time

func main() {

	PORT := ":8080"

	database.Connection()
	database.DB.AutoMigrate(models.User{})

	router := mux.NewRouter()
	funtions_EndPoints(router.PathPrefix("/api").Subrouter())

	StartTime = time.Now()
	fmt.Println("----------INIT USER SERVER " + PORT + " at: " + StartTime.Format("15:04:05") + "----------")
	http.ListenAndServe(PORT, router)
}

func funtions_EndPoints(router *mux.Router) {

	router.HandleFunc("/add", handlers.AddUser_handler).Methods("POST")
	router.HandleFunc("/update", handlers.UpdateUser_handler).Methods("PUT")
	router.HandleFunc("/delete/{email}", handlers.DeleteUser_handler).Methods("DELETE")
	router.HandleFunc("/search/{id}", handlers.GetUserById_handler).Methods("GET")
	router.HandleFunc("/all", handlers.GetAllUsers_handler).Methods("GET")

	router.HandleFunc("/login", handlers.Login_handler).Methods("POST")
	router.HandleFunc("/recover/{email}", handlers.RecoverPassword_handler).Methods("GET")
	router.HandleFunc("/update", handlers.UpdatePassword_handler).Methods("POST")

	router.HandleFunc("/health", handlers.CheckHealth_handler).Methods("GET")
	router.HandleFunc("/health/live", handlers.CheckLive_handler).Methods("GET")
	router.HandleFunc("/health/ready", handlers.CheckReady_handler).Methods("GET")
}
