package logging

import (
	"log"
	"os"
)

//LoadLogFile - load logging file
func LoadLogFile(logFilePath string) error {
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}
