package dto

import "time"

// TeachingClassRequest 教学班创建/更新请求
type TeachingClassRequest struct {
	ID           int64   `json:"id,omitempty"`
	ClassID      int64   `json:"class_id" binding:"required"`
	CourseID     int64   `json:"course_id" binding:"required"`
	SemesterID   int64   `json:"semester_id" binding:"required"`
	PricePerHour float64 `json:"price_per_hour" binding:"required"`
}

// TeachingClassWithDetail 教学班详细信息
type TeachingClassWithDetail struct {
	ID            int64   `json:"id"`
	ClassID       int64   `json:"class_id"`
	CourseID      int64   `json:"course_id"`
	SemesterID    int64   `json:"semester_id"`
	PricePerHour  float64 `json:"price_per_hour"`
	ClassName     string  `json:"class_name"`
	CourseName    string  `json:"course_name"`
	SemesterName  string  `json:"semester_name"`
	StudentCount  int     `json:"student_count"`
	ScheduleCount int     `json:"schedule_count"`
	DisplayName   string  `json:"display_name"` // 格式："班级名-课程名"
}

// BatchEnrollRequest 批量报名请求
type BatchEnrollRequest struct {
	ClassCourseID int64   `json:"class_course_id" binding:"required"`
	StudentIDs    []int64 `json:"student_ids" binding:"required,min=1"`
}

// BatchUnenrollRequest 批量退课请求
type BatchUnenrollRequest struct {
	ClassCourseID int64   `json:"class_course_id" binding:"required"`
	StudentIDs    []int64 `json:"student_ids" binding:"required,min=1"`
}

// EnrolledStudent 已报名学生信息
type EnrolledStudent struct {
	ID           int64     `json:"id"`
	StudentID    int64     `json:"student_id"`
	StudentName  string    `json:"student_name"`
	Contact      string    `json:"contact"`
	DiscountRate float64   `json:"discount_rate"`
	EnrolledAt   time.Time `json:"enrolled_at"`
}

// TeachingClassStudent 教学班学生管理
type TeachingClassStudent struct {
	ClassCourseID     int64             `json:"class_course_id"`
	DisplayName       string            `json:"display_name"`
	EnrolledStudents  []EnrolledStudent `json:"enrolled_students"`
	AvailableStudents []StudentInfo     `json:"available_students"`
}

// StudentInfo 学生基本信息
type StudentInfo struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Contact      string  `json:"contact"`
	DiscountRate float64 `json:"discount_rate"`
}
