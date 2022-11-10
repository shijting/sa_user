package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shjting0510/sa_user/inits"
	"github.com/shjting0510/sa_user/pkg/auth"
)

const (
	userName = "postgres"
	password = "123456"
	host     = "192.168.1.7"
	port     = 5432
	dbname   = "example"
)

func main() {
	r := gin.Default()
	inits.InitLogger("./logs")
	inits.InitDB(userName, password, host, port, dbname)
	r.POST("/login", auth.Login)
	r.Run(":8080")
}
