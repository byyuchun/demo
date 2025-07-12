package dto

// CreateClassCourseRequest 创建班级课程的请求体
type CreateClassCourseRequest struct {
	ClassID      int64   `json:"class_id" binding:"required"`
	CourseID     int64   `json:"course_id" binding:"required"`
	SemesterID   int64   `json:"semester_id" binding:"required"`
	PricePerHour float64 `json:"price_per_hour" binding:"required"`
}

// UpdateClassCourseRequest 更新班级课程的请求体
type UpdateClassCourseRequest struct {
	ID           int64   `json:"id" binding:"required"`
	ClassID      int64   `json:"class_id"`
	CourseID     int64   `json:"course_id"`
	SemesterID   int64   `json:"semester_id"`
	PricePerHour float64 `json:"price_per_hour"`
}