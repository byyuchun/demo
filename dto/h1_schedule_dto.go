package dto

import "time"

// CreateScheduleRequest 创建课表的请求体
type CreateScheduleRequest struct {
	ClassCourseID int64     `json:"class_course_id" binding:"required"`
	Date          time.Time `json:"date" binding:"required"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
}

// UpdateScheduleRequest 更新课表的请求体
type UpdateScheduleRequest struct {
	ID            int64     `json:"id" binding:"required"`
	ClassCourseID int64     `json:"class_course_id"`
	Date          time.Time `json:"date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

// CreateScheduleBatchRequest 批量创建课表请求体
type CreateScheduleBatchRequest struct {
	ClassCourseID int64     `json:"class_course_id" binding:"required"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	WeekDays      []int     `json:"week_days" binding:"required"` // 0=周日, 1=周一, ..., 6=周六
	StartHour     int       `json:"start_hour" binding:"required"`
	StartMinute   int       `json:"start_minute"`
	EndHour       int       `json:"end_hour" binding:"required"`
	EndMinute     int       `json:"end_minute"`
}

// ScheduleWithDetail 课表详情
type ScheduleWithDetail struct {
	ID            int64     `json:"id"`
	ClassCourseID int64     `json:"class_course_id"`
	Date          time.Time `json:"date"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	ClassName     string    `json:"class_name"`
	CourseName    string    `json:"course_name"`
	SemesterName  string    `json:"semester_name"`
}

// StudentSchedule 学生课表
type StudentSchedule struct {
	ScheduleID       int64     `json:"schedule_id"`
	ClassCourseID    int64     `json:"class_course_id"`
	Date             time.Time `json:"date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ClassName        string    `json:"class_name"`
	CourseName       string    `json:"course_name"`
	AttendanceStatus string    `json:"attendance_status"` // 出勤状态
}
