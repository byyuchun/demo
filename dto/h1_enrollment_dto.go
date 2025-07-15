package dto

import "time"

// CreateEnrollmentRequest 创建报名记录的请求体
type CreateEnrollmentRequest struct {
	StudentID     int64 `json:"student_id" binding:"required"`
	ClassCourseID int64 `json:"class_course_id" binding:"required"`
}

// UpdateEnrollmentRequest 更新报名记录的请求体
type UpdateEnrollmentRequest struct {
	ID            int64     `json:"id" binding:"required"`
	StudentID     int64     `json:"student_id"`
	ClassCourseID int64     `json:"class_course_id"`
	EnrolledAt    time.Time `json:"enrolled_at"`
}

// EnrollStudentRequest 学生报名请求体
type EnrollStudentRequest struct {
	StudentID     int64 `json:"student_id" binding:"required"`
	ClassCourseID int64 `json:"class_course_id" binding:"required"`
}

// EnrollmentWithDetail 报名记录详情
type EnrollmentWithDetail struct {
	ID            int64     `json:"id"`
	StudentID     int64     `json:"student_id"`
	ClassCourseID int64     `json:"class_course_id"`
	EnrolledAt    time.Time `json:"enrolled_at"`
	StudentName   string    `json:"student_name"`
	ClassName     string    `json:"class_name"`
	CourseName    string    `json:"course_name"`
	SemesterName  string    `json:"semester_name"`
	PricePerHour  float64   `json:"price_per_hour"`
}
