package dto

import "time"

// CreateSemesterRequest 创建学期的请求体
type CreateSemesterRequest struct {
	Name      string    `json:"name" binding:"required"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// UpdateSemesterRequest 更新学期的请求体
type UpdateSemesterRequest struct {
	ID        int64     `json:"id" binding:"required"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}