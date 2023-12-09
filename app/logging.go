package usrcanettocan

import (
	"fmt"
	"os"
	"runtime"

	log "github.com/sirupsen/logrus"
)

const EmptySource = ""

func getDebugLvl() log.Level {
	logLevelStr, found := os.LookupEnv("LOG_LEVEL")
	if !found {
		return log.InfoLevel
	}

	if logLevelStr == "DEBUG" {
		return log.DebugLevel
	}

	if logLevelStr == "INFO" {
		return log.InfoLevel
	}

	if logLevelStr == "WARN" {
		return log.WarnLevel
	}

	if logLevelStr == "ERROR" {
		return log.ErrorLevel
	}

	return log.InfoLevel
}

func init() {
	log.SetOutput(os.Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "02-01-2006 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(getDebugLvl())
}

func getDetails() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return file, line
}

func LogDebug(source string, msg string) {
	file, line := getDetails()
	debugMessage := msg

	if source != EmptySource {
		debugMessage = fmt.Sprintf("%s: %s", source, msg)
	}

	log.Debugf("%s:%d %s", file, line, debugMessage)
}
