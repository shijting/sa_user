package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shjting0510/sa_user/inits"
	"github.com/shjting0510/sa_user/pkg/auth"
	"github.com/shjting0510/sa_user/pkg/permission"
	"github.com/shjting0510/sa_user/pkg/roles_permissions"
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
	r.GET("/check_auth", auth.CheckPermission)

	r.GET("/permission", permission.GetPermissions)
	r.POST("/permission", permission.AddPermission)
	r.PUT("/permission", permission.EditPermission)
	r.DELETE("/permission/:id", permission.DelPermission)

	r.POST("/role_permission", roles_permissions.AddRolePermission)
	r.DELETE("/role_permission/:id", roles_permissions.DelRolePermission)
	r.GET("/permissions_with_roles", roles_permissions.GetRoleWithPermissionsByRoles)
	r.Run(":8080")
}
