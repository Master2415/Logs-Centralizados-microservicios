package communication

// import (
// 	"encoding/json"
// 	"log"
// 	"os"
// 	"sync"
// 	"user_server/models"

// 	"github.com/nats-io/nats.go"
// )

// var (
// 	once        sync.Once
// 	instance    *nats.Conn
// 	logger      *NatsLogger
// 	logsChannel chan *models.Log
// 	mutex       sync.Mutex
// )

// type NatsLogger struct {
// 	conn *nats.Conn
// }

// func Communicate() *NatsLogger {

// 	NATS_HOSTS := os.Getenv("US_NATS_SERVER")
// 	if NATS_HOSTS == "" {
// 		NATS_HOSTS = "localhost"
// 	}
// 	NATS_URL := "nats://" + NATS_HOSTS + ":4222"

// 	once.Do(func() {
// 		nc, err := nats.Connect(NATS_URL)
// 		if err != nil {
// 			log.Println("Error al conectar con NATS: %v", err)
// 		}
// 		instance = nc
// 	})

// 	return &NatsLogger{conn: instance}
// }

// func init() {
// 	logger = Communicate()
// 	logsChannel = make(chan *models.Log, 1000)
// 	go processLogs()
// }

// func processLogs() {
// 	for {
// 		logToSend := <-logsChannel
// 		logger.SendLog(logToSend)
// 	}
// }

// func (nl *NatsLogger) SendLog(newLog *models.Log) {
// 	if nl.CheckCommunicationLive() {
// 		jsonData, err := json.Marshal(newLog)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		if err := nl.conn.Publish("UsersServer", jsonData); err != nil {
// 			log.Println(err)
// 			addToLogsQueue(newLog)
// 			return
// 		}

// 		log.Println("LOG Sent.")
// 	} else {
// 		if logger.conn == nil || logger.conn.IsClosed() {
// 			logger.conn = Communicate().conn
// 		} else {
// 			addToLogsQueue(newLog)
// 		}
// 	}
// }

// func addToLogsQueue(logToAdd *models.Log) {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	logsChannel <- logToAdd
// }

// func (nl *NatsLogger) CheckCommunicationLive() bool {
// 	return nl.conn != nil && nl.conn.Status() == nats.CONNECTED
// }

// func (nl *NatsLogger) CheckCommunicationReady() bool {
// 	var subject = "test"
// 	var message = "Sample message"

// 	if nl.conn == nil {
// 		log.Println("Error: NATS connection is not established.")
// 		return false
// 	}

// 	if err := nl.conn.Publish(subject, []byte(message)); err != nil {
// 		log.Printf("Error al enviar mensaje a NATS: %v", err)
// 		return false
// 	}

// 	return true
// }
