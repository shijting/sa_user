package models

type Permission struct {
	Id          int    `json:"id" form:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Code        string `json:"code" form:"code" binding:"required"`
	Description string `json:"description" form:"description"`
	Action      string `json:"action" form:"action" binding:"required"`
	ParentID    int    `json:"parent_id" form:"parent_id"`
}
