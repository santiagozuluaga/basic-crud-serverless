package logger

import (
	"io"
	"log"
	"os"
)

func InitMock(buffer io.Writer) {
	log.SetOutput(buffer)
}

func ClearMock() {
	log.SetOutput(os.Stderr)
}
