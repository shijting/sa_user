package inits

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

var Log *log.Logger

func InitLogger(path string) {
	Log = log.New()
	Log.SetFormatter(&log.JSONFormatter{})
	fmt.Println("get env", os.Getenv("debug"))
	if len(os.Getenv("debug")) > 0 {
		Log.Info("debug env")
		Log.SetLevel(log.DebugLevel)
	} else {
		Log.SetLevel(log.ErrorLevel)
		Log.Out = NewFileLogWriter(path, "log")
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
		if os.IsExist(err) {
			return true
		}
	}
	return false
}
