package main

import (
	"fmt"
	"log"
	"os"
)

type telnetLogger struct {
	logger *log.Logger
}

func (logger *telnetLogger) Debug(v ...interface{}) {
	logger.logger.Print("[DEBUG]", v)
}

func (logger *telnetLogger) Debugf(format string, v ...interface{}) {
	logger.logger.Printf(format, fmt.Sprintf("[DEBUG] %v", v))
}

func (logger *telnetLogger) Error(v ...interface{}) {
	logger.logger.Print("[ERROR]", v)
}

func (logger *telnetLogger) Errorf(format string, v ...interface{}) {
	logger.logger.Printf(format, fmt.Sprintf("[ERROR] %v", v))
}

func (logger *telnetLogger) Trace(v ...interface{}) {
	logger.logger.Print("[TRACE]", v)
}

func (logger *telnetLogger) Tracef(format string, v ...interface{}) {
	logger.logger.Printf(format, fmt.Sprintf("[TRACE] %v", v))
}

func (logger *telnetLogger) Warn(v ...interface{}) {
	logger.logger.Print("[WARN]", v)
}

func (logger *telnetLogger) Warnf(format string, v ...interface{}) {
	logger.logger.Printf(format, fmt.Sprintf("[WARN] %v", v))
}

func createLogger(logFileName string) (*telnetLogger, *os.File, error) {
	logFile, err := getLogFile(logFileName)
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(logFile, "server", log.Ldate|log.Ltime|log.LUTC)
	return &telnetLogger{logger}, logFile, nil
}

func getLogFile(logFile string) (*os.File, error) {
	if file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); err == nil {
		return file, nil
	} else if file, err := os.Create(logFile); err == nil {
		return file, nil
	} else {
		return nil, err
	}
}
