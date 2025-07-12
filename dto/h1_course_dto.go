package dto

// CreateCourseRequest 创建课程的请求体
type CreateCourseRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateCourseRequest 更新课程的请求体
type UpdateCourseRequest struct {
	ID   int64  `json:"id" binding:"required"`
	Name string `json:"name"`
}