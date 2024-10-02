package health

import (
	"logs_server/communication"
	"logs_server/models"
	"logs_server/database"
	"time"
)

func LiveCheck() models.GeneralCheck {

	databaseCheck := database.CheckDatabaseLive()
	comunicationCheck := communication.Communicate().CheckCommunicationLive()

	databaseStatus := serviceCheking("Database", databaseCheck, "Live")
	comunicationStatus := serviceCheking("Comunicaction", comunicationCheck, "Live")

	servicesChek := []models.ServiceCheck{}
	servicesChek = append(servicesChek, databaseStatus, comunicationStatus)
	var serviceStatus string

	if databaseStatus.Status == "UP" && comunicationStatus.Status == "UP" {
		serviceStatus = "UP"
	} else {
		serviceStatus = "DOWN"
	}

	report := models.GeneralCheck{
		Status:  serviceStatus,
		Checks:  servicesChek,
		Version: "1.0.0",
		Uptime:  time.Since(database.StartTime).String(),
	}

	return report
}

func ReadyCheck() models.GeneralCheck {

	databaseCheck := database.CheckDatabaseReady()
	comunicationCheck := communication.Communicate().CheckCommunicationReady()

	databaseStatus := serviceCheking("Database", databaseCheck, "Ready")
	comunicationStatus := serviceCheking("Comunicaction", comunicationCheck, "Ready")

	servicesChek := []models.ServiceCheck{}
	servicesChek = append(servicesChek, databaseStatus, comunicationStatus)
	var serviceStatus string

	if databaseStatus.Status == "UP" && comunicationStatus.Status == "UP" {
		serviceStatus = "UP"
	} else {
		serviceStatus = "DOWN"
	}

	report := models.GeneralCheck{
		Status:  serviceStatus,
		Checks:  servicesChek,
		Version: "1.0.0",
		Uptime:  time.Since(database.StartTime).String(),
	}

	return report
}

func serviceCheking(serviceName string, serviceCheck bool, toEvaluate string) models.ServiceCheck {
	fromTime := time.Now()

	checkStatus := "DOWN"
	serviceStatus := "DOWN"

	if serviceCheck {
		serviceStatus = "READY"
	}

	if serviceStatus == "READY" {
		checkStatus = "UP"
	}

	return models.ServiceCheck{
		Data: models.CheckData{
			From:   fromTime,
			Status: serviceStatus,
		},
		Name:   serviceName + " " + toEvaluate + " conection check.",
		Status: checkStatus,
	}
}
