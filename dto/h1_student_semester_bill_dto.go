package dto

import "time"

// CreateStudentSemesterBillRequest 创建学生学期账单的请求体
type CreateStudentSemesterBillRequest struct {
	StudentID  int64   `json:"student_id" binding:"required"`
	SemesterID int64   `json:"semester_id" binding:"required"`
	TotalFee   float64 `json:"total_fee" binding:"required"`
	Status     string  `json:"status"`
}

// UpdateStudentSemesterBillRequest 更新学生学期账单的请求体
type UpdateStudentSemesterBillRequest struct {
	ID          int64     `json:"id" binding:"required"`
	StudentID   int64     `json:"student_id"`
	SemesterID  int64     `json:"semester_id"`
	TotalFee    float64   `json:"total_fee"`
	FinalizedAt time.Time `json:"finalized_at"`
	Status      string    `json:"status"`
}