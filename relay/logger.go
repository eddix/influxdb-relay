package relay

import (
	"log"
	"os"
	"time"
)

type Logger struct {
	prefix string
	low *log.Logger
	high *log.Logger

	lowFile *os.File
	highFile *os.File
}

func NewLogger(prefix string) (l *Logger, err error) {
	l = &Logger{prefix: prefix}
	if err = l.openLogs(); err != nil {
		return nil, err
	}

	go l.rotate()

	return l, nil
}

func (l Logger) Trace(format string, args ...interface{}) {
	l.low.Printf("[TRACE] " + format, args)
}

func (l Logger) Debug(format string, args ...interface{}) {
	l.low.Printf("[DEBUG] " + format, args)
}

func (l Logger) Notice(format string, args ...interface{}) {
	l.low.Printf("[NOTICE] " + format, args)
}

func (l Logger) Warning(format string, args ...interface{}) {
	l.high.Printf("[WARNING] " + format, args)
}

func (l Logger) Fatal(format string, args ...interface{}) {
	l.high.Printf("[FATAL] " + format, args)
}

func (l Logger) openLogs() (err error) {
	lowFile, err := os.OpenFile(l.prefix + ".log." + time.Now().Format("2006010215"), os.O_APPEND | os.O_CREATE,
		0644)
	if err != nil {
		return err
	}
	l.lowFile = lowFile
	l.low = log.New(lowFile, "", log.LstdFlags)

	highFile, err := os.OpenFile(l.prefix + ".log.wf." + time.Now().Format("2006010215"), os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	l.highFile = highFile
	l.high = log.New(highFile, "", log.LstdFlags)

	return nil
}

func (l Logger) rotate() {
	for {
		diff := 3600 - (time.Now().Unix() % 3600)
		time.Sleep(time.Duration(diff) * time.Second)

		oldLowFile := l.lowFile
		oldHighFile := l.highFile

		if err := l.openLogs(); err != nil {
			l.Fatal("Create new log files failed.")
		}
		if err := oldLowFile.Close(); err != nil {
			l.Fatal("Close old log file failed.")
		}
		if err := oldHighFile.Close(); err != nil {
			l.Fatal("Close old high file failed.")
		}
	}
}