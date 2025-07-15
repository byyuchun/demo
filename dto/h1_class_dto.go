package dto

// CreateClassRequest 创建班级的请求体
type CreateClassRequest struct {
	Name       string `json:"name" binding:"required"`
	SemesterID int64  `json:"semester_id" binding:"required"`
}

// UpdateClassRequest 更新班级的请求体
type UpdateClassRequest struct {
	ID         int64  `json:"id" binding:"required"`
	Name       string `json:"name"`
	SemesterID int64  `json:"semester_id"`
}

// ClassWithSemester 班级及学期信息
type ClassWithSemester struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SemesterID   int64  `json:"semester_id"`
	SemesterName string `json:"semester_name"`
}

// ClassCourseWithInfo 教学班及关联信息
type ClassCourseWithInfo struct {
	ID           int64   `json:"id"`
	ClassID      int64   `json:"class_id"`
	CourseID     int64   `json:"course_id"`
	SemesterID   int64   `json:"semester_id"`
	PricePerHour float64 `json:"price_per_hour"`
	CourseName   string  `json:"course_name"`
	SemesterName string  `json:"semester_name"`
}
