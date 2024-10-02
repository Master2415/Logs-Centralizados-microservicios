package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	communication "user_server/communication"
	"user_server/health"
	"user_server/models"
)

func CheckHealth_handler(w http.ResponseWriter, r *http.Request) {
	liveInfo := health.LiveCheck()
	readyInfo := health.ReadyCheck()

	response := map[string]interface{}{
		"live":  liveInfo,
		"ready": readyInfo,
	}

	formattedJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "HEALT",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Get API health",
		Description: "Healt info has been geted successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func CheckLive_handler(w http.ResponseWriter, r *http.Request) {
	liveInfo := health.LiveCheck()

	response := map[string]interface{}{
		"live": liveInfo,
	}

	formattedJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "HEALT",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Get API Live",
		Description: "Live info has been geted successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func CheckReady_handler(w http.ResponseWriter, r *http.Request) {
	ReadyInfo := health.ReadyCheck()

	response := map[string]interface{}{
		"Ready": ReadyInfo,
	}

	formattedJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "USERS-API",
		LogType:     "INFO",
		Module:      "HEALT",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Get API Ready",
		Description: "Ready info has been geted successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}
