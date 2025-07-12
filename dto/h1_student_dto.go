package dto

// CreateStudentRequest 创建学生的请求体
type CreateStudentRequest struct {
	Name         string  `json:"name" binding:"required"`
	Contact      string  `json:"contact"`
	DiscountRate float64 `json:"discount_rate"`
}

// UpdateStudentRequest 更新学生的请求体
type UpdateStudentRequest struct {
	ID           int64   `json:"id" binding:"required"`
	Name         string  `json:"name"`
	Contact      string  `json:"contact"`
	DiscountRate float64 `json:"discount_rate"`
}