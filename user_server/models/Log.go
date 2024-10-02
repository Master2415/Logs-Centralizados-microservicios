package models

type Log struct {
	ID          int    `json:"id" gorm:"primaryKey, autoIncrement"` // Identificador único del log
	AppName     string `json:"app_name"`                            // Nombre de la aplicación que genera el log
	LogType     string `json:"log_type"`                            // Tipo de log (e.g., info, warning, error)
	Module      string `json:"module"`                              // Clase o módulo que genera el log
	LogDateTime string `json:"log_date_time"`                       // Fecha y hora de generación del log
	Summary     string `json:"summary"`                             // Resumen del log
	Description string `json:"description"`                         // Descripción detallada del log
}
