package models

type Role struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	Code        string `json:"code" form:"code"`
	Description string `json:"description" form:"description"`
}
