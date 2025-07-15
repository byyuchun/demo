package dto

import "time"

// CreateStudentSemesterBillRequest 创建学期账单的请求体
type CreateStudentSemesterBillRequest struct {
	StudentID  int64   `json:"student_id" binding:"required"`
	SemesterID int64   `json:"semester_id" binding:"required"`
	TotalFee   float64 `json:"total_fee" binding:"required"`
	Status     string  `json:"status"`
}

// UpdateStudentSemesterBillRequest 更新学期账单的请求体
type UpdateStudentSemesterBillRequest struct {
	ID          int64     `json:"id" binding:"required"`
	StudentID   int64     `json:"student_id"`
	SemesterID  int64     `json:"semester_id"`
	TotalFee    float64   `json:"total_fee"`
	FinalizedAt time.Time `json:"finalized_at"`
	Status      string    `json:"status"`
}

// GenerateBillRequest 生成账单请求体
type GenerateBillRequest struct {
	StudentID  int64 `json:"student_id" binding:"required"`
	SemesterID int64 `json:"semester_id" binding:"required"`
}

// PayBillRequest 支付账单请求体
type PayBillRequest struct {
	StudentID  int64 `json:"student_id" binding:"required"`
	SemesterID int64 `json:"semester_id" binding:"required"`
}

// ClassCourseDetail 教学班费用明细
type ClassCourseDetail struct {
	ClassCourseID  int64   `json:"class_course_id"`
	ClassName      string  `json:"class_name"`
	CourseName     string  `json:"course_name"`
	PricePerHour   float64 `json:"price_per_hour"`
	PricePerClass  float64 `json:"price_per_class"`
	TotalSchedules int64   `json:"total_schedules"`
	AttendedCount  int64   `json:"attended_count"`
	ClassFee       float64 `json:"class_fee"`
}

// BillDetail 账单详情
type BillDetail struct {
	StudentID              int64               `json:"student_id"`
	StudentName            string              `json:"student_name"`
	SemesterID             int64               `json:"semester_id"`
	SemesterName           string              `json:"semester_name"`
	DiscountRate           float64             `json:"discount_rate"`
	TotalFeeBeforeDiscount float64             `json:"total_fee_before_discount"`
	FinalFee               float64             `json:"final_fee"`
	ClassCourseDetails     []ClassCourseDetail `json:"class_course_details"`
	BillStatus             string              `json:"bill_status"`
	FinalizedAt            time.Time           `json:"finalized_at"`
}

// BillSummary 账单摘要
type BillSummary struct {
	BillID       int64     `json:"bill_id"`
	StudentID    int64     `json:"student_id"`
	StudentName  string    `json:"student_name"`
	SemesterID   int64     `json:"semester_id"`
	SemesterName string    `json:"semester_name"`
	TotalFee     float64   `json:"total_fee"`
	Status       string    `json:"status"`
	FinalizedAt  time.Time `json:"finalized_at"`
}

// BillStatistics 账单统计信息
type BillStatistics struct {
	TotalStudents    int64         `json:"total_students"`
	TotalAmount      float64       `json:"total_amount"`
	TotalAttendances int64         `json:"total_attendances"`
	BillDetails      []BillSummary `json:"bill_details"`
}
