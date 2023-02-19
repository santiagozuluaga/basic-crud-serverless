package logger

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type Logger struct {
	service string
}

type LoggerObject struct {
	ID         string        `json:"id"`
	Service    string        `json:"service"`
	Level      LogLevel      `json:"level"`
	Properties LogProperties `json:"properties"`
	Time       time.Time     `json:"time"`
}

type LogProperties map[string]interface{}

type LogLevel string

var (
	InfoLevel    LogLevel = "info"
	WarningLevel LogLevel = "warning"
	ErrorLevel   LogLevel = "error"
	FatalLevel   LogLevel = "fatal"
)

func NewLogger(service string) Logger {
	log.SetFlags(0)

	return Logger{
		service: service,
	}
}

func (logger Logger) Log(logLevel LogLevel, properties LogProperties) {
	bytes, _ := json.Marshal(LoggerObject{
		ID:         uuid.New().String(),
		Service:    logger.service,
		Properties: properties,
		Level:      logLevel,
		Time:       time.Now(),
	})

	log.Println(string(bytes))
}
