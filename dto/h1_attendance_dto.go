package dto

import "time"

// CreateAttendanceRequest 创建考勤记录的请求体
type CreateAttendanceRequest struct {
	ScheduleID       int64     `json:"schedule_id" binding:"required"`
	StudentID        int64     `json:"student_id" binding:"required"`
	Status           string    `json:"status" binding:"required"`
	MakeupScheduleID int64     `json:"makeup_schedule_id"`
}

// UpdateAttendanceRequest 更新考勤记录的请求体
type UpdateAttendanceRequest struct {
	ID               int64     `json:"id" binding:"required"`
	ScheduleID       int64     `json:"schedule_id"`
	StudentID        int64     `json:"student_id"`
	CheckedAt        time.Time `json:"checked_at"`
	Status           string    `json:"status"`
	MakeupScheduleID int64     `json:"makeup_schedule_id"`
}