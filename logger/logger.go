package logger

import (
	"log"
	"os"
)

const(
	POSTGRES = "POSTGRES"
	MEMCACHE = "MEMCACHE"
	CONTROLLER = "CONTROLLER"
	BIDDER = "BIDDER"
	SERVER = "SERVER"
)

type Logger struct {
	lg *log.Logger
}	

var loggerInstance *Logger

// Returns logger instance
func InitLogger(pref string) *Logger {
	lg := log.New(os.Stdout, pref+": ", log.LstdFlags)
	// loggerInstance = &Logger{lg:lg}
	// "mini-bidder: "+
	return &Logger{lg:lg}
}

func GetLoggerInstance(pref string) *Logger {
	if loggerInstance == nil {
		loggerInstance = InitLogger(pref)
	} else {
		loggerInstance.lg.SetPrefix(pref+":\t")
	}
	
	return loggerInstance
}

func (l *Logger) Print(v ...interface{}) {
	l.lg.Print(v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.lg.Fatal(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.lg.Fatalf(format, v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.lg.Printf(format, v...)
}