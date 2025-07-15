package controller

import (
	"context"
	"demo/dal/model"
	"demo/dal/query"
	"demo/dto"
	"demo/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// CreateStudentSemesterBill godoc
// @Summary      Create a new student semester bill
// @Description  Create a new student semester bill
// @Tags         H1StudentSemesterBill
// @Accept       json
// @Produce      json
// @Param        bill  body      dto.CreateStudentSemesterBillRequest  true  "Bill info"
// @Success      200      {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill [post]
func CreateStudentSemesterBill(ctx *gin.Context) {
	var req dto.CreateStudentSemesterBillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	bill := model.H1StudentSemesterBill{
		StudentID:   req.StudentID,
		SemesterID:  req.SemesterID,
		TotalFee:    req.TotalFee,
		FinalizedAt: time.Now(),
		Status:      req.Status,
	}

	if err := query.H1StudentSemesterBill.WithContext(context.Background()).Create(&bill); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, bill, "创建成功")
}

// GetStudentSemesterBill godoc
// @Summary      Get a student semester bill by ID
// @Description  Get a student semester bill by ID
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        id   path      int  true  "StudentSemesterBill ID"
// @Success      200  {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill/{id} [get]
func GetStudentSemesterBill(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	bill, err := query.H1StudentSemesterBill.WithContext(context.Background()).Where(query.H1StudentSemesterBill.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生学期账单不存在")
		return
	}

	response.Success(ctx, bill, "查询成功")
}

// DeleteStudentSemesterBill godoc
// @Summary      Delete a student semester bill by ID
// @Description  Delete a student semester bill by ID
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        id   path      int  true  "StudentSemesterBill ID"
// @Success      200  {object}  response.Response
// @Router       /h1/student-semester-bill/{id} [delete]
func DeleteStudentSemesterBill(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1StudentSemesterBill.WithContext(context.Background()).Where(query.H1StudentSemesterBill.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateStudentSemesterBill godoc
// @Summary      Update a student semester bill
// @Description  Update a student semester bill
// @Tags         H1StudentSemesterBill
// @Accept       json
// @Produce      json
// @Param        bill  body      dto.UpdateStudentSemesterBillRequest  true  "Bill info"
// @Success      200      {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill [put]
func UpdateStudentSemesterBill(ctx *gin.Context) {
	var req dto.UpdateStudentSemesterBillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	bill, err := query.H1StudentSemesterBill.WithContext(context.Background()).Where(query.H1StudentSemesterBill.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "账单不存在")
		return
	}

	if req.StudentID != 0 {
		bill.StudentID = req.StudentID
	}
	if req.SemesterID != 0 {
		bill.SemesterID = req.SemesterID
	}
	if req.TotalFee != 0 {
		bill.TotalFee = req.TotalFee
	}
	if !req.FinalizedAt.IsZero() {
		bill.FinalizedAt = req.FinalizedAt
	}
	if req.Status != "" {
		bill.Status = req.Status
	}

	if err := query.H1StudentSemesterBill.WithContext(context.Background()).Save(bill); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, bill, "更新成功")
}

// GetAllStudentSemesterBills godoc
// @Summary      Get all student semester bills
// @Description  Get all student semester bills
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Success      200  {object}  response.Response{data=[]model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill [get]
func GetAllStudentSemesterBills(ctx *gin.Context) {
	bills, err := query.H1StudentSemesterBill.WithContext(context.Background()).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(ctx, bills, "查询成功")
}

// GenerateSemesterBill 生成学期账单
// @Summary      生成学生学期账单
// @Description  根据学生出勤情况自动计算并生成学期账单
// @Tags         H1StudentSemesterBill
// @Accept       json
// @Produce      json
// @Param        generate  body      dto.GenerateBillRequest  true  "生成账单信息"
// @Success      200      {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill/generate [post]
func GenerateSemesterBill(ctx *gin.Context) {
	var req dto.GenerateBillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证学生是否存在
	student, err := query.H1Student.WithContext(context.Background()).Where(query.H1Student.ID.Eq(req.StudentID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生不存在")
		return
	}

	// 验证学期是否存在
	_, err = query.H1Semester.WithContext(context.Background()).Where(query.H1Semester.ID.Eq(req.SemesterID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学期不存在")
		return
	}

	// 检查是否已经生成过账单
	existingBill, _ := query.H1StudentSemesterBill.WithContext(context.Background()).
		Where(query.H1StudentSemesterBill.StudentID.Eq(req.StudentID)).
		Where(query.H1StudentSemesterBill.SemesterID.Eq(req.SemesterID)).
		First()
	if existingBill != nil {
		response.Fail(ctx, http.StatusBadRequest, "该学期账单已存在")
		return
	}

	// 获取学生在该学期的所有报名记录
	enrollments, err := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(req.StudentID)).
		Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "获取报名记录失败")
		return
	}

	var totalFee float64 = 0

	// 遍历每个报名的教学班
	for _, enrollment := range enrollments {
		// 获取教学班信息
		classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(enrollment.ClassCourseID)).
			Where(query.H1ClassCourse.SemesterID.Eq(req.SemesterID)).
			First()
		if err != nil {
			continue // 跳过不属于该学期的报名
		}

		// 获取该教学班的所有课表
		schedules, err := query.H1Schedule.WithContext(context.Background()).
			Where(query.H1Schedule.ClassCourseID.Eq(classCourse.ID)).
			Find()
		if err != nil {
			continue
		}

		// 统计出勤次数
		var attendedCount int64 = 0
		for _, schedule := range schedules {
			// 查询该学生在这堂课的出勤记录
			attendance, err := query.H1Attendance.WithContext(context.Background()).
				Where(query.H1Attendance.ScheduleID.Eq(schedule.ID)).
				Where(query.H1Attendance.StudentID.Eq(req.StudentID)).
				First()

			if err == nil && (attendance.Status == "出勤" || attendance.Status == "补签") {
				attendedCount++
			}
		}

		// 计算该教学班的费用：出勤次数 × 每次课价格
		// 注意：PricePerHour 实际上是每次课的价格，不是每小时价格
		classFee := float64(attendedCount) * classCourse.PricePerHour
		totalFee += classFee
	}

	// 应用学生折扣
	finalFee := totalFee * student.DiscountRate

	// 创建账单
	bill := model.H1StudentSemesterBill{
		StudentID:   req.StudentID,
		SemesterID:  req.SemesterID,
		TotalFee:    finalFee,
		FinalizedAt: time.Now(),
		Status:      "未结算",
	}

	if err := query.H1StudentSemesterBill.WithContext(context.Background()).Create(&bill); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "生成账单失败")
		return
	}

	response.Success(ctx, bill, "账单生成成功")
}

// GetBillDetail 获取账单详情
// @Summary      获取账单详细信息
// @Description  获取学生学期账单的详细信息，包括各教学班的费用明细
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        student_id   query     int  true  "学生ID"
// @Param        semester_id  query     int  true  "学期ID"
// @Success      200  {object}  response.Response{data=dto.BillDetail}
// @Router       /h1/student-semester-bill/detail [get]
func GetBillDetail(ctx *gin.Context) {
	studentIDStr := ctx.Query("student_id")
	semesterIDStr := ctx.Query("semester_id")

	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "学生ID格式错误")
		return
	}

	semesterID, err := strconv.ParseInt(semesterIDStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "学期ID格式错误")
		return
	}

	// 获取学生信息
	student, err := query.H1Student.WithContext(context.Background()).Where(query.H1Student.ID.Eq(studentID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生不存在")
		return
	}

	// 获取学期信息
	semester, err := query.H1Semester.WithContext(context.Background()).Where(query.H1Semester.ID.Eq(semesterID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学期不存在")
		return
	}

	// 获取账单
	bill, _ := query.H1StudentSemesterBill.WithContext(context.Background()).
		Where(query.H1StudentSemesterBill.StudentID.Eq(studentID)).
		Where(query.H1StudentSemesterBill.SemesterID.Eq(semesterID)).
		First()

	// 获取学生在该学期的所有报名记录和费用明细
	enrollments, err := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(studentID)).
		Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "获取报名记录失败")
		return
	}

	var classCourseDetails []dto.ClassCourseDetail
	var totalFee float64 = 0

	for _, enrollment := range enrollments {
		// 获取教学班信息
		classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(enrollment.ClassCourseID)).
			Where(query.H1ClassCourse.SemesterID.Eq(semesterID)).
			First()
		if err != nil {
			continue
		}

		// 获取班级和课程信息
		class, _ := query.H1Class.WithContext(context.Background()).Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()
		course, _ := query.H1Course.WithContext(context.Background()).Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

		// 获取该教学班的所有课表
		schedules, err := query.H1Schedule.WithContext(context.Background()).
			Where(query.H1Schedule.ClassCourseID.Eq(classCourse.ID)).
			Find()
		if err != nil {
			continue
		}

		// 统计出勤次数
		var attendedCount int64 = 0
		var totalSchedules int64 = int64(len(schedules))

		for _, schedule := range schedules {
			attendance, err := query.H1Attendance.WithContext(context.Background()).
				Where(query.H1Attendance.ScheduleID.Eq(schedule.ID)).
				Where(query.H1Attendance.StudentID.Eq(studentID)).
				First()

			if err == nil && (attendance.Status == "出勤" || attendance.Status == "补签") {
				attendedCount++
			}
		}

		// 计算该教学班的费用：出勤次数 × 每次课价格
		// 注意：PricePerHour 实际上是每次课的价格，不是每小时价格
		pricePerClass := classCourse.PricePerHour
		classFee := float64(attendedCount) * pricePerClass
		totalFee += classFee

		detail := dto.ClassCourseDetail{
			ClassCourseID:  classCourse.ID,
			ClassName:      class.Name,
			CourseName:     course.Name,
			PricePerHour:   classCourse.PricePerHour,
			PricePerClass:  pricePerClass,
			TotalSchedules: totalSchedules,
			AttendedCount:  attendedCount,
			ClassFee:       classFee,
		}
		classCourseDetails = append(classCourseDetails, detail)
	}

	// 应用学生折扣
	finalFee := totalFee * student.DiscountRate

	billDetail := dto.BillDetail{
		StudentID:              studentID,
		StudentName:            student.Name,
		SemesterID:             semesterID,
		SemesterName:           semester.Name,
		DiscountRate:           student.DiscountRate,
		TotalFeeBeforeDiscount: totalFee,
		FinalFee:               finalFee,
		ClassCourseDetails:     classCourseDetails,
		BillStatus:             "",
		FinalizedAt:            time.Time{},
	}

	if bill != nil {
		billDetail.BillStatus = bill.Status
		billDetail.FinalizedAt = bill.FinalizedAt
	}

	response.Success(ctx, billDetail, "查询成功")
}

// PayBill 支付账单
// @Summary      支付学期账单
// @Description  标记学期账单为已支付
// @Tags         H1StudentSemesterBill
// @Accept       json
// @Produce      json
// @Param        payment  body      dto.PayBillRequest  true  "支付信息"
// @Success      200      {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill/pay [post]
func PayBill(ctx *gin.Context) {
	var req dto.PayBillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取账单
	bill, err := query.H1StudentSemesterBill.WithContext(context.Background()).
		Where(query.H1StudentSemesterBill.StudentID.Eq(req.StudentID)).
		Where(query.H1StudentSemesterBill.SemesterID.Eq(req.SemesterID)).
		First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "账单不存在")
		return
	}

	if bill.Status == "已结算" {
		response.Fail(ctx, http.StatusBadRequest, "账单已结算")
		return
	}

	// 更新账单状态
	bill.Status = "已结算"
	bill.FinalizedAt = time.Now()

	if err := query.H1StudentSemesterBill.WithContext(context.Background()).Save(bill); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "支付失败")
		return
	}

	response.Success(ctx, bill, "支付成功")
}

// GetAllBillsWithDetail 获取所有账单详情
// @Summary      获取所有账单的详细信息
// @Description  获取系统中所有学期账单的详细信息，用于管理员查看
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        semester_id  query     int  false  "学期ID过滤"
// @Success      200  {object}  response.Response{data=[]dto.BillSummary}
// @Router       /h1/student-semester-bill/all-with-detail [get]
func GetAllBillsWithDetail(ctx *gin.Context) {
	semesterIDStr := ctx.Query("semester_id")

	// 构建查询
	billQuery := query.H1StudentSemesterBill.WithContext(context.Background())
	if semesterIDStr != "" {
		semesterID, err := strconv.ParseInt(semesterIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "学期ID格式错误")
			return
		}
		billQuery = billQuery.Where(query.H1StudentSemesterBill.SemesterID.Eq(semesterID))
	}

	bills, err := billQuery.Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.BillSummary
	for _, bill := range bills {
		// 获取学生信息
		student, _ := query.H1Student.WithContext(context.Background()).Where(query.H1Student.ID.Eq(bill.StudentID)).First()
		// 获取学期信息
		semester, _ := query.H1Semester.WithContext(context.Background()).Where(query.H1Semester.ID.Eq(bill.SemesterID)).First()

		summary := dto.BillSummary{
			BillID:       bill.ID,
			StudentID:    bill.StudentID,
			StudentName:  student.Name,
			SemesterID:   bill.SemesterID,
			SemesterName: semester.Name,
			TotalFee:     bill.TotalFee,
			Status:       bill.Status,
			FinalizedAt:  bill.FinalizedAt,
		}
		result = append(result, summary)
	}

	response.Success(ctx, result, "查询成功")
}

// GetBillStatistics 获取账单统计信息
// @Summary      获取账单统计信息
// @Description  根据签到记录统计生成账单信息
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        semester_id  query     int  false  "学期ID过滤"
// @Success      200  {object}  response.Response{data=dto.BillStatistics}
// @Router       /h1/student-semester-bill/statistics [get]
func GetBillStatistics(ctx *gin.Context) {
	semesterIDStr := ctx.Query("semester_id")

	// 获取所有学生
	students, err := query.H1Student.WithContext(context.Background()).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询学生失败")
		return
	}

	var totalAmount float64 = 0
	var totalStudents int64 = 0
	var totalAttendances int64 = 0
	var billDetails []dto.BillSummary

	for _, student := range students {
		// 获取学生的报名记录
		enrollments, err := query.H1Enrollment.WithContext(context.Background()).
			Where(query.H1Enrollment.StudentID.Eq(student.ID)).
			Find()
		if err != nil {
			continue
		}

		var studentTotalFee float64 = 0
		var studentAttendanceCount int64 = 0

		for _, enrollment := range enrollments {
			// 获取教学班信息
			classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
				Where(query.H1ClassCourse.ID.Eq(enrollment.ClassCourseID)).
				First()
			if err != nil {
				continue
			}

			// 如果指定了学期过滤
			if semesterIDStr != "" {
				semesterID, err := strconv.ParseInt(semesterIDStr, 10, 64)
				if err != nil {
					response.Fail(ctx, http.StatusBadRequest, "学期ID格式错误")
					return
				}
				if classCourse.SemesterID != semesterID {
					continue
				}
			}

			// 获取该教学班的所有课表
			schedules, err := query.H1Schedule.WithContext(context.Background()).
				Where(query.H1Schedule.ClassCourseID.Eq(classCourse.ID)).
				Find()
			if err != nil {
				continue
			}

			// 统计出勤次数
			var attendedCount int64 = 0
			for _, schedule := range schedules {
				attendance, err := query.H1Attendance.WithContext(context.Background()).
					Where(query.H1Attendance.ScheduleID.Eq(schedule.ID)).
					Where(query.H1Attendance.StudentID.Eq(student.ID)).
					First()

				if err == nil && (attendance.Status == "出勤" || attendance.Status == "补签") {
					attendedCount++
				}
			}

			// 计算课时费用：出勤次数 × 每次课价格
			// 注意：PricePerHour 实际上是每次课的价格，不是每小时价格
			pricePerClass := classCourse.PricePerHour
			classFee := float64(attendedCount) * pricePerClass
			studentTotalFee += classFee
			studentAttendanceCount += attendedCount
		}

		// 应用学生折扣
		finalFee := studentTotalFee * student.DiscountRate

		if studentTotalFee > 0 {
			totalStudents++
			totalAmount += finalFee
			totalAttendances += studentAttendanceCount

			// 查找或创建账单记录
			if semesterIDStr != "" {
				semesterID, _ := strconv.ParseInt(semesterIDStr, 10, 64)

				// 查找现有账单
				existingBill, err := query.H1StudentSemesterBill.WithContext(context.Background()).
					Where(query.H1StudentSemesterBill.StudentID.Eq(student.ID)).
					Where(query.H1StudentSemesterBill.SemesterID.Eq(semesterID)).
					First()

				var bill *model.H1StudentSemesterBill
				if err != nil {
					// 创建新账单
					newBill := model.H1StudentSemesterBill{
						StudentID:   student.ID,
						SemesterID:  semesterID,
						TotalFee:    finalFee,
						FinalizedAt: time.Now(),
						Status:      "待支付",
					}
					if createErr := query.H1StudentSemesterBill.WithContext(context.Background()).Create(&newBill); createErr == nil {
						bill = &newBill
					}
				} else {
					// 更新现有账单
					existingBill.TotalFee = finalFee
					existingBill.FinalizedAt = time.Now()
					if updateErr := query.H1StudentSemesterBill.WithContext(context.Background()).Save(existingBill); updateErr == nil {
						bill = existingBill
					}
				}

				if bill != nil {
					// 获取学期信息
					semester, _ := query.H1Semester.WithContext(context.Background()).
						Where(query.H1Semester.ID.Eq(semesterID)).First()

					semesterName := ""
					if semester != nil {
						semesterName = semester.Name
					}

					billDetails = append(billDetails, dto.BillSummary{
						BillID:       bill.ID,
						StudentID:    student.ID,
						StudentName:  student.Name,
						SemesterID:   semesterID,
						SemesterName: semesterName,
						TotalFee:     finalFee,
						Status:       bill.Status,
						FinalizedAt:  bill.FinalizedAt,
					})
				}
			}
		}
	}

	statistics := dto.BillStatistics{
		TotalStudents:    totalStudents,
		TotalAmount:      totalAmount,
		TotalAttendances: totalAttendances,
		BillDetails:      billDetails,
	}

	response.Success(ctx, statistics, "统计成功")
}
