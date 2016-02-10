package jelly

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"sync"
)

// ErrEmptyLogName error value if log name supplied is empty
var ErrEmptyLogName = errors.New("jelly: log name provided is empty")

// Logger represents a single log instance.
type Logger struct {
	Path    string
	Name    string
	logFile *log.Logger
	mutex   sync.Mutex
}

// NewLog constructs a new Logger instance.
func NewLog(logName string) (*Logger, error) {
	if logName == "" {
		return nil, ErrEmptyLogName
	}

	if !strings.Contains(logName, ".log") {
		logName = logName + ".log"
	}

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	base := fmt.Sprintf("%s%c%s", usr.HomeDir, os.PathSeparator, ".jelly")
	if _, err := os.Stat(base); os.IsNotExist(err) {
		os.Mkdir(base, os.ModePerm)
	}

	path := fmt.Sprintf("%s%c%s", base, os.PathSeparator, logName)

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Path:    path,
		Name:    logName,
		logFile: log.New(file, "", log.Ldate|log.Ltime|log.LUTC),
	}, nil
}

// Info writes v to the log file, at the INFO level. This method should
// be used to log information users need to know about.
func (l *Logger) Info(v ...interface{}) {
	l.write("INFO", stringify(v))
}

// Debug writes v to the log file, at the DEBUG level. This method should be
// used to log debug information a developer needs to know about.
func (l *Logger) Debug(v ...interface{}) {
	l.write("DEBUG", stringify(v))
}

// Die writes v to the log file and then causes the current program to exit with
// status code 1.
func (l *Logger) Die(v ...interface{}) {
	l.write("DIE", stringify(v))
	os.Exit(1)
}

// write writes level and msg strings to the log file. The write is protected
// by a mutex.
func (l *Logger) write(level, msg string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	_, file, line, _ := runtime.Caller(2)
	index := strings.LastIndex(file, "/")
	if index != -1 {
		file = file[index+1 : len(file)]
	}
	l.logFile.Println(fmt.Sprintf("%s:%d %s -", file, line, level), msg)
}

// stringify turns the contents of v into a string.
func stringify(v []interface{}) string {
	s := len(v)
	var buf bytes.Buffer
	for i := 0; i < s; i++ {
		buf.WriteString(fmt.Sprintf("%v", v[i]))
		buf.WriteString(" ")
	}
	return buf.String()
}
