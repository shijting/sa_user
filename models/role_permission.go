package models

type RolePermission struct {
	Id             int    `json:"id" form:"id"`
	RoleCode       string `json:"role_code" form:"role_code" binding:"required"`
	PermissionCode string `json:"permission_code" form:"permission_code" binding:"required"`
}
