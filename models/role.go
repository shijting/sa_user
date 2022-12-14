package models

type Role struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Code        string `json:"code" form:"code" binding:"required"`
	Description string `json:"description" form:"description"`
}
