package dto

// CreateClassRequest 创建班级的请求体
type CreateClassRequest struct {
	Name       string `json:"name" binding:"required"`
	SemesterID int64  `json:"semester_id" binding:"required"`
}

// UpdateClassRequest 更新班级的请求体
type UpdateClassRequest struct {
	ID         int64  `json:"id" binding:"required"`
	Name       string `json:"name"`
	SemesterID int64  `json:"semester_id"`
}