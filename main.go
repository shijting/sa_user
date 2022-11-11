package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shjting0510/sa_user/inits"
	"github.com/shjting0510/sa_user/pkg/auth"
)

const (
	DbUserName = "postgres"
	DbPassword = "123456"
	DbHost     = "192.168.1.7"
	DbPort     = 5432
	DbDbname   = "example"
)

func main() {
	r := gin.Default()
	inits.InitLogger("./logs")
	inits.InitDB(DbUserName, DbPassword, DbHost, DbPort, DbDbname)
	r.POST("/login", auth.Login)
	r.GET("/authorize", auth.CheckPermission)
	r.Run(":8080")
}
