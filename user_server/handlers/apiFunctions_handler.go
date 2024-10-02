package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	communication "user_server/communication"
	"user_server/database"
	"user_server/models"
	"user_server/security"

	"github.com/gorilla/mux"
)

func AddUser_handler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error: Invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Password == "" || user.Email == "" {
		http.Error(w, "Error: User, Password & Email info are obligatory", http.StatusBadRequest)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "ADD-USER",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to create a user",
			Description: "No mandatory information was entered",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	createdUser, err := database.AddUser(user)
	if err != nil {
		http.Error(w, "Error: This email is already in use", http.StatusConflict)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "ADD-USER",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to create a user",
			Description: "An email in use was entered",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	tokenString := security.GenerateToken(createdUser)
	if tokenString == "" {
		http.Error(w, "Error: Failed to sign token", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "ADD-USER",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Create a user",
		Description: "A user has been created successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, tokenString)
}

func UpdateUser_handler(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error: Problem to covert JSON", http.StatusBadRequest)
		return
	}

	if !security.IsValidToken(r, strconv.Itoa(user.Id)) {
		http.Error(w, "Error: Invalid token", http.StatusUnauthorized)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-USER",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update user",
			Description: "Invalid token provided",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	updatedUser, err := database.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error: Email to update is already in use", http.StatusConflict)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-USER",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update user",
			Description: "Email to update is already in use",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	formattedJSON, err := json.MarshalIndent(updatedUser, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "UPDATE-USER",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "User updated successfully",
		Description: "A user has been updated successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func DeleteUser_handler(w http.ResponseWriter, r *http.Request) {

	info := mux.Vars(r)
	email := info["email"]

	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Error: Email not found", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "DELETE-USER",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to delete user",
			Description: "Email not found in the database",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	if !security.IsValidToken(r, strconv.Itoa(user.Id)) {
		http.Error(w, "Error: Invalid token", http.StatusUnauthorized)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "DELETE-USER",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to delete user",
			Description: "Invalid token provided",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	err = database.DeleteUser(user)
	if err != nil {
		http.Error(w, "Error: Failed to delete user", http.StatusNotFound)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "DELETE-USER",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "User deleted successfully",
		Description: "User was deleted successfully",
	}
	communication.Communicate().SendLog(&log)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "\nOK: User with ID %d was deleted.\n", user.Id)
}

func GetUserById_handler(w http.ResponseWriter, r *http.Request) {
	if !security.IsValidToken(r, "") {
		http.Error(w, "Error: Invalid token", http.StatusUnauthorized)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "GET-USER-BY-ID",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to get user by ID",
			Description: "Invalid token provided",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	params := mux.Vars(r)
	user, err := database.GetUserById(params["id"])
	if err != nil {
		http.Error(w, "Error: Failed to search user id ("+params["id"]+")", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "GET-USER-BY-ID",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to get user by ID",
			Description: "Failed to retrieve user ID from database",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	if user == nil {
		http.Error(w, "Error: User not found", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "GET-USER-BY-ID",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to get user by ID",
			Description: "User ID not found in the database",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	formattedJSON, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "GET-USER-BY-ID",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "User retrieved successfully",
		Description: "User ID retrieved successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func GetAllUsers_handler(w http.ResponseWriter, r *http.Request) {

	if !security.IsValidToken(r, "") {
		http.Error(w, "Error: Invalid token", http.StatusUnauthorized)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "GET-ALL-USERS",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to get all users",
			Description: "Invalid token provided",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	users, err := database.GetUsers(page, pageSize)
	if err != nil {
		http.Error(w, "Error: Failed to get users", http.StatusUnauthorized)
		return
	}

	formattedJSON, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "GET-ALL-USERS",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "All users retrieved successfully",
		Description: "All users were retrieved successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func Login_handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Error: Method not allowed", http.StatusMethodNotAllowed)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "LOGIN",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to log in",
			Description: "Method not allowed",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	_, err := database.SearchUser(&user)

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Error: Email & Password info are obligatory", http.StatusBadRequest)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "LOGIN",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to log in",
			Description: "Missing required email or password",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	tokenString := security.GenerateToken(&user)
	if tokenString == "" {
		http.Error(w, "Error: Failed to sign token", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Error: User not found", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "LOGIN",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to log in",
			Description: "User not found in the database",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "LOGIN",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Login successful",
		Description: "User logged in successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, tokenString)
}

func RecoverPassword_handler(w http.ResponseWriter, r *http.Request) {

	info := mux.Vars(r)
	email := info["email"]

	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Error: Email not found", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "RECOVER-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to recover password",
			Description: "Email not found in the database",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	if !security.IsValidToken(r, strconv.Itoa(user.Id)) {
		http.Error(w, "Error: Invalid token", http.StatusUnauthorized)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "RECOVER-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to recover password",
			Description: "Invalid token provided",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	password, err := database.RecoverPassword(info["email"])
	if err != nil {
		http.Error(w, "Error: User not found", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "RECOVER-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to recover password",
			Description: "Error occurred while recovering password",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "RECOVER-PASSWORD",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Password recovered successfully",
		Description: "Password was recovered successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "\nOK: Password was recuperated.\nPassword is: ")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(password))
	fmt.Fprintf(w, "\n")
}

func UpdatePassword_handler(w http.ResponseWriter, r *http.Request) {

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error: Invalid JSON", http.StatusBadRequest)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update password",
			Description: "Invalid JSON input",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Error: Password & Email info to update is obligatory", http.StatusBadRequest)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update password",
			Description: "Missing required email or password",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	newPassword := user.Password
	user, err = database.GetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, "Error: Email not found", http.StatusNotFound)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update password",
			Description: "Email not found in the database",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	if !security.IsValidToken(r, strconv.Itoa(user.Id)) {
		http.Error(w, "Error: Invalid token", http.StatusUnauthorized)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update password",
			Description: "Invalid token provided",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	user.Password = newPassword
	newPassword, err = database.UpdatePassword(user)
	if err != nil {
		http.Error(w, "Error: New password must be different from the current one", http.StatusBadRequest)
		log := models.Log{
			AppName:     "USERS-API",
			LogType:     "ERROR",
			Module:      "UPDATE-PASSWORD",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to update password",
			Description: "New password must be different from the current one",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "UPDATE-PASSWORD",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Password updated successfully",
		Description: "User password updated successfully",
	}
	communication.Communicate().SendLog(&log)

	message := "\nOK: Password was Updated.\nNew password is: " + newPassword

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}
