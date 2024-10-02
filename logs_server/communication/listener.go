package communication

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"logs_server/database"
	"logs_server/models"

	"github.com/nats-io/nats.go"
)

func ListenNotifications() {

	NATS_HOSTS := os.Getenv("LS_NATS_SERVER")
	if NATS_HOSTS == "" {
		NATS_HOSTS = "localhost"
	}

	NATS_URL := "nats://" + NATS_HOSTS + ":4222"
	nc, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Fatalf("Error al conectar con NATS: %v", err)
	}
	defer nc.Close()

	subscribes := []string{"UsersServer", "LogsServer"}

	for _, subject := range subscribes {
		sub, err := nc.Subscribe(subject, func(m *nats.Msg) {
			handleMessage(m)
		})

		if err != nil {
			log.Printf("Error al suscribirse al canal %s: %v", subject, err)
		}
		defer sub.Unsubscribe()
	}

	select {}
}

func handleMessage(m *nats.Msg) {
	var logEntry models.Log
	err := json.Unmarshal(m.Data, &logEntry)
	if err != nil {
		log.Printf("Error al deserializar el mensaje: %v", err)
		return
	}

	err = database.DB.Create(&logEntry).Error
	if err != nil {
		log.Printf("Error al guardar en la base de datos: %v", err)
	}
	fmt.Printf("LOG recibido: %+v\n", logEntry)
}
