package models

//	type BaseInfo[T any] struct {
//		Result    T      `json:"result"`
//		Success   bool   `json:"success"`
//		ErrorInfo string `json:"errorInfo"`
//	}
type UsersModel struct {
	UserID   int64  `form:"user_id" json:"user_id"`
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
