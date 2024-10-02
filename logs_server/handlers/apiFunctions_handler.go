package handlers

import (
	"encoding/json"
	"fmt"
	"logs_server/communication"
	"logs_server/database"
	"logs_server/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetLogs_handler(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	logType := r.URL.Query().Get("logType")

	logs, err := database.GetLogs(page, pageSize, startDate, endDate, logType)
	if err != nil {
		http.Error(w, "Error: Failed to get logs", http.StatusInternalServerError)
		return
	}

	if logs == nil {
		http.Error(w, "Error: No information to show", http.StatusNotFound)
		log := models.Log{
			AppName:     "LOGS-API",
			LogType:     "ERROR",
			Module:      "GET-LOGS",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Not found info to show",
			Description: "Entered values that did not return info",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	formattedJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "LOGS-API",
		LogType:     "INFO",
		Module:      "GET-LOGS",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Get logs successfully",
		Description: "Logs has been geted successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func GetLogsByApp_handler(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	logType := r.URL.Query().Get("logType")

	params := mux.Vars(r)
	logs, err := database.GetLogsByApp(params["appName"], page, pageSize, startDate, endDate, logType)
	if err != nil {
		http.Error(w, "Error: Failed to get logs", http.StatusInternalServerError)
		return
	}

	if logs == nil {
		http.Error(w, "Error: No information to show", http.StatusNotFound)
		log := models.Log{
			AppName:     "LOGS-API",
			LogType:     "ERROR",
			Module:      "GET-LOGS-BY-APP",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Not found info to show",
			Description: "Entered values that did not return info",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	formattedJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		http.Error(w, "Error: Failed to format JSON", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "LOGS-API",
		LogType:     "INFO",
		Module:      "GET-LOGS-BY-APP",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Get logs successfully",
		Description: "Logs has been geted successfully",
	}
	communication.Communicate().SendLog(&log)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(formattedJSON)
}

func AddLogs_handler(w http.ResponseWriter, r *http.Request) {
	var newLog models.Log
	err := json.NewDecoder(r.Body).Decode(&newLog)
	if err != nil {
		http.Error(w, "Error: Cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	if newLog.AppName == "" || newLog.LogType == "" || newLog.Module == "" ||
		newLog.Summary == "" || newLog.Description == "" {
		http.Error(w, "Error: All info are obligatory", http.StatusBadRequest)
		log := models.Log{
			AppName:     "LOGS-API",
			LogType:     "ERROR",
			Module:      "ADD-LOG",
			LogDateTime: time.Now().Format(time.RFC3339),
			Summary:     "Failed to create a log",
			Description: "No mandatory information was entered",
		}
		communication.Communicate().SendLog(&log)
		return
	}

	newLog.LogDateTime = time.Now().Format(time.RFC3339)
	logID, err := database.AddLog(&newLog)
	if err != nil {
		http.Error(w, "Error: Could not create LOG", http.StatusInternalServerError)
		return
	}

	log := models.Log{
		AppName:     "LOGS-API",
		LogType:     "INFO",
		Module:      "ADD-LOG",
		LogDateTime: time.Now().Format(time.RFC3339),
		Summary:     "Log added successfully",
		Description: "Log has been added successfully",
	}
	communication.Communicate().SendLog(&log)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "OK: LOG with id %d was created.\n", logID)
}
