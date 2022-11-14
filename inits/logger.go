package inits

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

var Log *log.Logger

func InitLogger(path string) {
	Log = log.New()

	if len(os.Getenv("debug")) > 0 {
		Log.SetLevel(log.DebugLevel)
		Log.SetFormatter(&log.TextFormatter{})
		Log.Info("dev mode")
	} else {
		Log.SetLevel(log.ErrorLevel)
		Log.Out = NewFileLogWriter(path, "log")
		Log.SetFormatter(&log.JSONFormatter{})
	}
}

func NewFileLogWriter(filePath string, fileName string) io.Writer {
	var name string
	if !Exist(filePath) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return nil
		}
	}
	name = filePath + "/" + fileName
	if !strings.HasSuffix(fileName, ".log") {
		name = filePath + "/" + fileName + ".log"
	}
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func Exist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return false
}
