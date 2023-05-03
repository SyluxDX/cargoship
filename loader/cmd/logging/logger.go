package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger stasd
type Logger struct {
	file        *os.File
	logger      *log.Logger
	numberlines int
}

func (l *Logger) Init(path string, toConsole bool) {
	var err error
	// check folder
	dir, _ := filepath.Split(path)
	_ = os.MkdirAll(dir, 0755)

	// open file
	l.file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Panicln(err)
	}
	// configure level loggers
	if toConsole {
		l.logger = log.New(io.MultiWriter(os.Stdout, l.file), "", log.LstdFlags|log.LUTC)
	} else {
		l.logger = log.New(l.file, "", log.LstdFlags|log.LUTC)
	}
	l.numberlines = 0
}

func (l *Logger) LogDebug(message string) {
	l.logger.Printf("DEBUG %s", message)
	l.numberlines += 1
}

func (l *Logger) LogInfo(message string) {
	l.logger.Printf("INFO %s", message)
	l.numberlines += 1
}

func (l *Logger) LogWarn(message string) {
	l.logger.Printf("WARN %s", message)
	l.numberlines += 1
}

func (l *Logger) LogError(message string) {
	l.logger.Printf("ERROR %s", message)
	l.numberlines += 1
}

func (l *Logger) Close() {
	// close file
	l.file.Close()

	// delete file if no lines are written
	if l.numberlines == 0 {
		os.Remove(l.file.Name())
	}
}
