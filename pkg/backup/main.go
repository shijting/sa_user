package backup

import (
	"bytes"
	"context"
	"github.com/robfig/cron/v3"
	"github.com/shjting0510/sa_user/inits"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"time"
)

const (
	port   = "5432"
	dbname = "example"
	path   = "/opt/test_backup2/"
)

func Backup() {
	quitCh := make(chan struct{}, 1)

	c := cron.New()
	defer c.Stop()
	c.AddFunc("@daily", func() {
		exit, _ := PathExists(path)
		if !exit {
			err := os.MkdirAll(path, 0777)
			if err != nil {
				quitCh <- struct{}{}
				inits.Log.WithFields(log.Fields{"action": "backup.Backup"}).Error(err)
				log.Println(err)
				return
			}
		}

		fileName := path + dbname + time.Now().Format("20060102150405") + ".bak"

		cmd := exec.CommandContext(context.Background(), "./backup.sh", dbname, fileName, port)
		output := bytes.Buffer{}
		stderr := bytes.Buffer{}
		cmd.Stdout = &output
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			quitCh <- struct{}{}
		}
		if err := cmd.Run(); err != nil {
			inits.Log.WithFields(log.Fields{"action": "backup.Backup"}).Error()
			quitCh <- struct{}{}
			return
		}
		// TODO
		log.Println(fileName)
	})

	c.Start()

	<-quitCh
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
