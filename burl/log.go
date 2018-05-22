package burl

import (
	"fmt"
	"os"
	"time"
)

var logger []logMessage

type logLevel int

const (
	LOG_DEBUG logLevel = iota
	LOG_INFO
	LOG_WARN
	LOG_ERROR
)

func (l logLevel) String() string {
	switch l {
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_WARN:
		return "WARNING"
	case LOG_ERROR:
		return "ERROR"
	default:
		return "???"
	}
}

type logMessage struct {
	level   logLevel
	time    time.Time
	message string
}

func (l logMessage) String() string {
	return "[" + l.time.Format(time.Stamp) + "] " + l.level.String() + ": " + l.message
}

func init() {
	logger = make([]logMessage, 0, 1000)
	LogInfo("BURL Engine Online!")
}

func log(level logLevel, message ...interface{}) {
	logger = append(logger, logMessage{
		level:   level,
		time:    time.Now(),
		message: fmt.Sprint(message...),
	})

	//if we're in debug mode, add the new message to the debugger window
	if debug {
		debugger.logList.Append(logger[len(logger)-1].String())
		debugger.logList.ScrollToBottom()
	}
}

func outputLogToDisk() {
	f, err := os.Create("log.txt")
	if err != nil {
		return
	}
	defer f.Close()

	for _, m := range logger {
		f.WriteString(m.String() + "\n")
	}
}

func LogDebug(m ...interface{}) {
	log(LOG_DEBUG, m...)
}

func LogInfo(m ...interface{}) {
	log(LOG_INFO, m...)
}

func LogWarning(m ...interface{}) {
	log(LOG_WARN, m...)
}

func LogError(m ...interface{}) {
	log(LOG_ERROR, m...)
}
