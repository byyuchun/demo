package dto

import "time"

// CreateAttendanceRequest 创建出勤记录的请求体
type CreateAttendanceRequest struct {
	ScheduleID       int64  `json:"schedule_id" binding:"required"`
	StudentID        int64  `json:"student_id" binding:"required"`
	Status           string `json:"status" binding:"required"`
	MakeupScheduleID *int64 `json:"makeup_schedule_id"`
}

// UpdateAttendanceRequest 更新出勤记录的请求体
type UpdateAttendanceRequest struct {
	ID               int64     `json:"id" binding:"required"`
	ScheduleID       int64     `json:"schedule_id"`
	StudentID        int64     `json:"student_id"`
	CheckedAt        time.Time `json:"checked_at"`
	Status           string    `json:"status"`
	MakeupScheduleID *int64    `json:"makeup_schedule_id"`
}

// CheckInRequest 学生打卡请求体
type CheckInRequest struct {
	ScheduleID int64 `json:"schedule_id" binding:"required"`
	StudentID  int64 `json:"student_id" binding:"required"`
}

// ApplyMakeupRequest 申请补签请求体
type ApplyMakeupRequest struct {
	OriginalScheduleID int64 `json:"original_schedule_id" binding:"required"`
	MakeupScheduleID   int64 `json:"makeup_schedule_id" binding:"required"`
	StudentID          int64 `json:"student_id" binding:"required"`
}

// AdminMakeupRequest 管理员补签请求体
type AdminMakeupRequest struct {
	ScheduleID       int64  `json:"schedule_id" binding:"required"`
	StudentID        int64  `json:"student_id" binding:"required"`
	MakeupScheduleID *int64 `json:"makeup_schedule_id"`
}

// MarkAbsentRequest 标记缺勤请求体
type MarkAbsentRequest struct {
	ScheduleID int64  `json:"schedule_id" binding:"required"`
	StudentID  int64  `json:"student_id" binding:"required"`
	Status     string `json:"status" binding:"required"` // "请假" 或 "旷课"
}

// AttendanceWithDetails 带详细信息的出勤记录
type AttendanceWithDetails struct {
	ID               int64     `json:"id"`
	ScheduleID       int64     `json:"schedule_id"`
	StudentID        int64     `json:"student_id"`
	StudentName      string    `json:"student_name"`
	CheckedAt        time.Time `json:"checked_at"`
	Status           string    `json:"status"`
	MakeupScheduleID *int64    `json:"makeup_schedule_id"`
	ScheduleDate     time.Time `json:"schedule_date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ClassName        string    `json:"class_name"`
	CourseName       string    `json:"course_name"`
}
