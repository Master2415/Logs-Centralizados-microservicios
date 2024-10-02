package database

import (
	"logs_server/models"
)

func GetLogs(page, pageSize int, startDate, endDate, logType string) ([]models.Log, error) {
	var logs []models.Log
	offset := (page - 1) * pageSize

	query := DB.Offset(offset).Limit(pageSize)

	if startDate != "" && endDate != "" {
		query = query.Where("log_date_time BETWEEN ? AND ?", startDate, endDate)
	}

	if logType != "" {
		query = query.Where("log_type = ?", logType)
	}

	err := query.Order("log_date_time DESC").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, nil
	}

	return logs, nil
}

func GetLogsByApp(appName string, page, pageSize int, startDate, endDate, logType string) ([]models.Log, error) {
	var logs []models.Log
	offset := (page - 1) * pageSize

	query := DB.Offset(offset).Limit(pageSize)

	query = query.Where("app_name = ?", appName)

	if startDate != "" && endDate != "" {
		query = query.Where("log_date_time BETWEEN ? AND ?", startDate, endDate)
	}

	if logType != "" {
		query = query.Where("log_type = ?", logType)
	}

	err := query.Order("log_date_time DESC").Find(&logs).Error
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, nil
	}

	return logs, nil
}

func AddLog(log *models.Log) (int, error) {
	result := DB.Create(log)

	if result.Error != nil {
		return 0, result.Error
	}

	return log.ID, nil
}

func VerifyApp(appName string) (bool, error) {
	var log models.Log
	err := DB.Where("app_name = ?", appName).First(&log).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
