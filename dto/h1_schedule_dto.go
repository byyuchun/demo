package dto

import "time"

// CreateScheduleRequest 创建课程表的请求体
type CreateScheduleRequest struct {
	ClassCourseID int64     `json:"class_course_id" binding:"required"`
	Date          time.Time `json:"date" binding:"required"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

// UpdateScheduleRequest 更新课程表的请求体
type UpdateScheduleRequest struct {
	ID            int64     `json:"id" binding:"required"`
	ClassCourseID int64     `json:"class_course_id"`
	Date          time.Time `json:"date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}