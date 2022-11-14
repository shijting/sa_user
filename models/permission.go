package models

type Permission struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required,min=1,max=30"`
	Code        string `json:"code" form:"code" binding:"required,min=1,max=30"`
	Description string `json:"description" form:"description"`
	Action      string `json:"action" form:"action" binding:"required,min=1,max=50"`
	ParentID    int    `json:"parent_id" form:"parent_id"`
}
