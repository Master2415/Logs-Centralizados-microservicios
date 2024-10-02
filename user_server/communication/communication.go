package communication

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"user_server/models"

	"github.com/nats-io/nats.go"
)

var (
	once        sync.Once
	instance    *nats.Conn
	logger      *NatsLogger
	logsChannel chan *models.Log
	mutex       sync.Mutex
	natsURL     string
)

type NatsLogger struct {
	conn *nats.Conn
}

func handleError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %v\n", msg, err)
	}
}

func Communicate() *NatsLogger {
	once.Do(func() {
		nc, err := nats.Connect(natsURL)
		handleError(err, "Error al conectar con NATS")
		instance = nc
	})
	return &NatsLogger{conn: instance}
}

func init() {
	NATS_HOSTS := os.Getenv("US_NATS_SERVER")
	if NATS_HOSTS == "" {
		NATS_HOSTS = "localhost"
	}
	natsURL = "nats://" + NATS_HOSTS + ":4222"
	logger = Communicate()
	logsChannel = make(chan *models.Log, 1000)
	go processLogs()
}

func processLogs() {
	var wg sync.WaitGroup
	for logToSend := range logsChannel {
		wg.Add(1)
		go func(logEntry *models.Log) {
			defer wg.Done()
			logger.SendLog(logEntry)
		}(logToSend)
	}
	wg.Wait()
}

func (nl *NatsLogger) SendLog(newLog *models.Log) {
	if !nl.CheckCommunicationLive() {
		logger.conn = Communicate().conn
	}

	jsonData, err := json.Marshal(newLog)
	if err != nil {
		log.Println(err)
		return
	}

	if err := nl.conn.Publish("UsersServer", jsonData); err != nil {
		log.Println(err)
		addToLogsQueue(newLog)
		return
	}

	log.Println("LOG Sent.")
}

func addToLogsQueue(logToAdd *models.Log) {
	mutex.Lock()
	defer mutex.Unlock()
	select {
	case logsChannel <- logToAdd:
	default:
		log.Println("Logs queue full, dropping log entry.")
	}
}

func (nl *NatsLogger) CheckCommunicationLive() bool {
	return nl.conn != nil && nl.conn.Status() == nats.CONNECTED
}

func (nl *NatsLogger) CheckCommunicationReady() bool {
	const subject = "test"
	const message = "Sample message"

	if nl.conn == nil || nl.conn.Status() != nats.CONNECTED {
		return false
	}

	if err := nl.conn.Publish(subject, []byte(message)); err != nil {
		return false
	}

	return true
}
