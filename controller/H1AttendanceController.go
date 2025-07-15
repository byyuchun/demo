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

// CreateAttendance godoc
// @Summary      Create a new attendance record
// @Description  Create a new attendance record
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        attendance  body      dto.CreateAttendanceRequest  true  "Attendance info"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance [post]
func CreateAttendance(ctx *gin.Context) {
	var req dto.CreateAttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	attendance := model.H1Attendance{
		ScheduleID:       req.ScheduleID,
		StudentID:        req.StudentID,
		CheckedAt:        time.Now(),
		Status:           req.Status,
		MakeupScheduleID: req.MakeupScheduleID, // 现在是指针类型，可以正确处理 nil
	}

	if err := query.H1Attendance.WithContext(context.Background()).Create(&attendance); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, attendance, "创建成功")
}

// GetAttendance godoc
// @Summary      Get an attendance by ID
// @Description  Get an attendance by ID
// @Tags         H1Attendance
// @Produce      json
// @Param        id   path      int  true  "Attendance ID"
// @Success      200  {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance/{id} [get]
func GetAttendance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "ID格式错误")
		return
	}

	attendance, err := query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "考勤记录不存在")
		return
	}

	response.Success(ctx, attendance, "查询成功")
}

// DeleteAttendance godoc
// @Summary      Delete an attendance by ID
// @Description  Delete an attendance by ID
// @Tags         H1Attendance
// @Produce      json
// @Param        id   path      int  true  "Attendance ID"
// @Success      200  {object}  response.Response
// @Router       /h1/attendance/{id} [delete]
func DeleteAttendance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "ID格式错误")
		return
	}

	_, err = query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateAttendance godoc
// @Summary      Update an attendance record
// @Description  Update an attendance record
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        attendance  body      dto.UpdateAttendanceRequest  true  "Attendance info"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance [put]
func UpdateAttendance(ctx *gin.Context) {
	var req dto.UpdateAttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	attendance, err := query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "考勤记录不存在")
		return
	}

	if req.ScheduleID != 0 {
		attendance.ScheduleID = req.ScheduleID
	}
	if req.StudentID != 0 {
		attendance.StudentID = req.StudentID
	}
	if !req.CheckedAt.IsZero() {
		attendance.CheckedAt = req.CheckedAt
	}
	if req.Status != "" {
		attendance.Status = req.Status
	}
	if req.MakeupScheduleID != nil {
		attendance.MakeupScheduleID = req.MakeupScheduleID
	}

	if err := query.H1Attendance.WithContext(context.Background()).Save(attendance); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, attendance, "更新成功")
}

// GetAllAttendances godoc
// @Summary      Get all attendances with details
// @Description  Get all attendances with student and class information
// @Tags         H1Attendance
// @Produce      json
// @Param        status       query     string  false  "出勤状态过滤"
// @Param        start_date   query     string  false  "开始日期 (YYYY-MM-DD)"
// @Param        end_date     query     string  false  "结束日期 (YYYY-MM-DD)"
// @Success      200  {object}  response.Response{data=[]dto.AttendanceWithDetails}
// @Router       /h1/attendance [get]
func GetAllAttendances(ctx *gin.Context) {
	statusFilter := ctx.Query("status")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	// 构建查询
	attendanceQuery := query.H1Attendance.WithContext(context.Background())

	// 状态过滤
	if statusFilter != "" {
		attendanceQuery = attendanceQuery.Where(query.H1Attendance.Status.Eq(statusFilter))
	}

	// 日期范围过滤
	if startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			response.Fail(ctx, http.StatusBadRequest, "日期格式错误，请使用 YYYY-MM-DD")
			return
		}

		// 需要通过课表的日期来过滤
		schedules, err := query.H1Schedule.WithContext(context.Background()).
			Where(query.H1Schedule.Date.Between(startDate, endDate)).
			Find()
		if err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "查询课表失败")
			return
		}

		scheduleIDs := make([]int64, len(schedules))
		for i, schedule := range schedules {
			scheduleIDs[i] = schedule.ID
		}

		if len(scheduleIDs) > 0 {
			attendanceQuery = attendanceQuery.Where(query.H1Attendance.ScheduleID.In(scheduleIDs...))
		} else {
			// 如果没有找到符合日期范围的课表，返回空结果
			response.Success(ctx, []dto.AttendanceWithDetails{}, "查询成功")
			return
		}
	}

	attendances, err := attendanceQuery.Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	// 获取详细信息
	var result []dto.AttendanceWithDetails
	for _, attendance := range attendances {
		// 获取学生信息
		student, _ := query.H1Student.WithContext(context.Background()).
			Where(query.H1Student.ID.Eq(attendance.StudentID)).First()

		// 获取课表信息
		schedule, _ := query.H1Schedule.WithContext(context.Background()).
			Where(query.H1Schedule.ID.Eq(attendance.ScheduleID)).First()

		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(schedule.ClassCourseID)).First()

		// 获取班级信息
		class, _ := query.H1Class.WithContext(context.Background()).
			Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()

		// 获取课程信息
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

		detail := dto.AttendanceWithDetails{
			ID:               attendance.ID,
			ScheduleID:       attendance.ScheduleID,
			StudentID:        attendance.StudentID,
			CheckedAt:        attendance.CheckedAt,
			Status:           attendance.Status,
			MakeupScheduleID: attendance.MakeupScheduleID,
			ScheduleDate:     schedule.Date,
			StartTime:        schedule.StartTime,
			EndTime:          schedule.EndTime,
			ClassName:        "",
			CourseName:       "",
		}

		if student != nil {
			detail.StudentName = student.Name
		}
		if class != nil {
			detail.ClassName = class.Name
		}
		if course != nil {
			detail.CourseName = course.Name
		}

		result = append(result, detail)
	}

	response.Success(ctx, result, "查询成功")
}

// CheckIn 学生打卡
// @Summary      学生上课打卡
// @Description  学生按照课表进行上课打卡
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        checkin  body      dto.CheckInRequest  true  "打卡信息"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance/check-in [post]
func CheckIn(ctx *gin.Context) {
	var req dto.CheckInRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证课表是否存在
	schedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ScheduleID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课表不存在")
		return
	}

	// 验证学生是否已报名该课程
	_, err = query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(req.StudentID)).
		Where(query.H1Enrollment.ClassCourseID.Eq(schedule.ClassCourseID)).
		First()
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "学生未报名该课程")
		return
	}

	// 检查是否已经打卡
	existingAttendance, _ := query.H1Attendance.WithContext(context.Background()).
		Where(query.H1Attendance.ScheduleID.Eq(req.ScheduleID)).
		Where(query.H1Attendance.StudentID.Eq(req.StudentID)).
		First()
	if existingAttendance != nil {
		response.Fail(ctx, http.StatusBadRequest, "已经打过卡了")
		return
	}

	// 创建出勤记录
	attendance := model.H1Attendance{
		ScheduleID: req.ScheduleID,
		StudentID:  req.StudentID,
		CheckedAt:  time.Now(),
		Status:     "出勤", // 正常出勤
	}

	if err := query.H1Attendance.WithContext(context.Background()).Create(&attendance); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "打卡失败")
		return
	}

	// 签到成功后，自动处理账单
	if err := updateStudentBill(req.StudentID, schedule.ClassCourseID); err != nil {
		// 账单更新失败，记录错误但不影响签到结果
		// 在实际应用中应该使用日志系统
		println("账单更新失败:", err.Error())
	}

	response.Success(ctx, attendance, "打卡成功")
}

// updateStudentBill 更新学生账单
func updateStudentBill(studentID int64, classCourseID int64) error {
	// 获取教学班信息
	classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(classCourseID)).
		First()
	if err != nil {
		println("获取教学班信息失败:", err.Error())
		return err
	}

	// 获取学生信息（需要折扣率）
	student, err := query.H1Student.WithContext(context.Background()).
		Where(query.H1Student.ID.Eq(studentID)).
		First()
	if err != nil {
		println("获取学生信息失败:", err.Error())
		return err
	}

	// 计算该学生在当前学期的总费用
	totalFee := calculateStudentSemesterFee(studentID, classCourse.SemesterID, student.DiscountRate)
	println("计算出的总费用:", totalFee, "学生ID:", studentID, "学期ID:", classCourse.SemesterID)

	// 检查该学生在当前学期是否已有账单记录
	existingBill, err := query.H1StudentSemesterBill.WithContext(context.Background()).
		Where(query.H1StudentSemesterBill.StudentID.Eq(studentID)).
		Where(query.H1StudentSemesterBill.SemesterID.Eq(classCourse.SemesterID)).
		First()

	if err != nil {
		// 账单不存在，创建新账单
		newBill := model.H1StudentSemesterBill{
			StudentID:   studentID,
			SemesterID:  classCourse.SemesterID,
			TotalFee:    totalFee,
			FinalizedAt: time.Now(),
			Status:      "待支付",
		}
		println("创建新账单，学生ID:", studentID, "学期ID:", classCourse.SemesterID, "费用:", totalFee)
		err = query.H1StudentSemesterBill.WithContext(context.Background()).Create(&newBill)
		if err != nil {
			println("创建账单失败:", err.Error())
		} else {
			println("账单创建成功，ID:", newBill.ID)
		}
		return err
	} else {
		// 账单已存在，更新费用
		println("更新现有账单，账单ID:", existingBill.ID, "原费用:", existingBill.TotalFee, "新费用:", totalFee)
		existingBill.TotalFee = totalFee
		existingBill.FinalizedAt = time.Now()
		err = query.H1StudentSemesterBill.WithContext(context.Background()).Save(existingBill)
		if err != nil {
			println("更新账单失败:", err.Error())
		} else {
			println("账单更新成功")
		}
		return err
	}
}

// calculateStudentSemesterFee 计算学生在指定学期的总费用
func calculateStudentSemesterFee(studentID int64, semesterID int64, discountRate float64) float64 {
	// 获取学生在该学期的所有报名记录
	enrollments, err := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(studentID)).
		Find()
	if err != nil {
		println("获取报名记录失败:", err.Error())
		return 0
	}

	println("找到报名记录数量:", len(enrollments))
	var totalFee float64 = 0

	for _, enrollment := range enrollments {
		// 获取教学班信息，确保属于指定学期
		classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(enrollment.ClassCourseID)).
			Where(query.H1ClassCourse.SemesterID.Eq(semesterID)).
			First()
		if err != nil {
			println("教学班不属于指定学期，跳过. ClassCourseID:", enrollment.ClassCourseID)
			continue // 跳过不属于该学期的报名
		}

		println("处理教学班ID:", classCourse.ID, "每次课价格:", classCourse.PricePerHour)

		// 获取该教学班的所有课表
		schedules, err := query.H1Schedule.WithContext(context.Background()).
			Where(query.H1Schedule.ClassCourseID.Eq(classCourse.ID)).
			Find()
		if err != nil {
			println("获取课表失败:", err.Error())
			continue
		}

		println("找到课表数量:", len(schedules))

		// 统计有效出勤次数（出勤和补签）
		var attendedCount int64 = 0
		for _, schedule := range schedules {
			attendance, err := query.H1Attendance.WithContext(context.Background()).
				Where(query.H1Attendance.ScheduleID.Eq(schedule.ID)).
				Where(query.H1Attendance.StudentID.Eq(studentID)).
				First()

			if err == nil && (attendance.Status == "出勤" || attendance.Status == "补签") {
				attendedCount++
				println("找到有效出勤记录，课表ID:", schedule.ID, "状态:", attendance.Status)
			}
		}

		// 计算该教学班的费用：有效出勤次数 × 每次课价格
		// 注意：price_per_hour 实际上是每次课的价格，不是每小时价格
		classFee := float64(attendedCount) * classCourse.PricePerHour
		println("教学班费用计算: 出勤次数", attendedCount, "× 每次课价格", classCourse.PricePerHour, "= 费用", classFee)
		totalFee += classFee
	}

	// 应用学生折扣
	finalFee := totalFee * discountRate
	println("费用计算完成: 总费用", totalFee, "× 折扣率", discountRate, "= 最终费用", finalFee)
	return finalFee
}

// ApplyMakeup 申请补签
// @Summary      申请课程补签
// @Description  学生申请缺勤课程的补签
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        makeup  body      dto.ApplyMakeupRequest  true  "补签申请信息"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance/apply-makeup [post]
func ApplyMakeup(ctx *gin.Context) {
	var req dto.ApplyMakeupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证原课表是否存在
	schedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.OriginalScheduleID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "原课表不存在")
		return
	}

	// 验证补课课表是否存在
	makeupSchedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.MakeupScheduleID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "补课课表不存在")
		return
	}

	// 验证两个课表是否属于同一个教学班
	if schedule.ClassCourseID != makeupSchedule.ClassCourseID {
		response.Fail(ctx, http.StatusBadRequest, "补课课表与原课表不属于同一教学班")
		return
	}

	// 检查是否已有出勤记录
	existingAttendance, _ := query.H1Attendance.WithContext(context.Background()).
		Where(query.H1Attendance.ScheduleID.Eq(req.OriginalScheduleID)).
		Where(query.H1Attendance.StudentID.Eq(req.StudentID)).
		First()

	if existingAttendance != nil {
		// 更新现有记录为补签
		existingAttendance.Status = "补签"
		existingAttendance.MakeupScheduleID = &req.MakeupScheduleID
		existingAttendance.CheckedAt = time.Now()

		if err := query.H1Attendance.WithContext(context.Background()).Save(existingAttendance); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "补签失败")
			return
		}

		// 补签成功后，更新账单
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.OriginalScheduleID)).First()
		if schedule != nil {
			updateStudentBill(req.StudentID, schedule.ClassCourseID)
		}

		response.Success(ctx, existingAttendance, "补签成功")
	} else {
		// 创建新的补签记录
		attendance := model.H1Attendance{
			ScheduleID:       req.OriginalScheduleID,
			StudentID:        req.StudentID,
			CheckedAt:        time.Now(),
			Status:           "补签",
			MakeupScheduleID: &req.MakeupScheduleID,
		}

		if err := query.H1Attendance.WithContext(context.Background()).Create(&attendance); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "补签失败")
			return
		}

		// 补签成功后，更新账单
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.OriginalScheduleID)).First()
		if schedule != nil {
			updateStudentBill(req.StudentID, schedule.ClassCourseID)
		}

		response.Success(ctx, attendance, "补签成功")
	}
}

// GetStudentAttendance 获取学生出勤记录
// @Summary      获取学生的出勤记录
// @Description  获取指定学生在指定学期的出勤记录
// @Tags         H1Attendance
// @Produce      json
// @Param        student_id   query     int  true  "学生ID"
// @Param        semester_id  query     int  false "学期ID"
// @Success      200  {object}  response.Response{data=[]dto.AttendanceWithDetails}
// @Router       /h1/attendance/student [get]
func GetStudentAttendance(ctx *gin.Context) {
	studentIDStr := ctx.Query("student_id")
	semesterIDStr := ctx.Query("semester_id")

	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "学生ID格式错误")
		return
	}

	// 构建查询
	attendanceQuery := query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.StudentID.Eq(studentID))

	// 如果指定了学期，需要通过课表关联过滤
	if semesterIDStr != "" {
		semesterID, err := strconv.ParseInt(semesterIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "学期ID格式错误")
			return
		}

		// 先获取该学期的所有课表ID
		schedules, err := query.H1Schedule.WithContext(context.Background()).
			Select(query.H1Schedule.ID).
			LeftJoin(query.H1ClassCourse, query.H1Schedule.ClassCourseID.EqCol(query.H1ClassCourse.ID)).
			Where(query.H1ClassCourse.SemesterID.Eq(semesterID)).
			Find()
		if err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "查询课表失败")
			return
		}

		scheduleIDs := make([]int64, len(schedules))
		for i, schedule := range schedules {
			scheduleIDs[i] = schedule.ID
		}

		attendanceQuery = attendanceQuery.Where(query.H1Attendance.ScheduleID.In(scheduleIDs...))
	}

	attendances, err := attendanceQuery.Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	// 获取详细信息
	var result []dto.AttendanceWithDetails
	for _, attendance := range attendances {
		// 获取课表信息
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(attendance.ScheduleID)).First()
		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).Where(query.H1ClassCourse.ID.Eq(schedule.ClassCourseID)).First()
		// 获取班级信息
		class, _ := query.H1Class.WithContext(context.Background()).Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()
		// 获取课程信息
		course, _ := query.H1Course.WithContext(context.Background()).Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

		detail := dto.AttendanceWithDetails{
			ID:               attendance.ID,
			ScheduleID:       attendance.ScheduleID,
			StudentID:        attendance.StudentID,
			CheckedAt:        attendance.CheckedAt,
			Status:           attendance.Status,
			MakeupScheduleID: attendance.MakeupScheduleID,
			ScheduleDate:     schedule.Date,
			StartTime:        schedule.StartTime,
			EndTime:          schedule.EndTime,
			ClassName:        class.Name,
			CourseName:       course.Name,
		}
		result = append(result, detail)
	}

	response.Success(ctx, result, "查询成功")
}

// AdminMakeup 管理员补签
// @Summary      管理员为学生补签
// @Description  管理员为缺勤学生进行补签操作
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        makeup  body      dto.AdminMakeupRequest  true  "补签信息"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance/admin-makeup [post]
func AdminMakeup(ctx *gin.Context) {
	var req dto.AdminMakeupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 检查是否已存在考勤记录
	existingAttendance, err := query.H1Attendance.WithContext(context.Background()).
		Where(query.H1Attendance.ScheduleID.Eq(req.ScheduleID)).
		Where(query.H1Attendance.StudentID.Eq(req.StudentID)).
		First()

	if err != nil {
		// 不存在考勤记录，创建新的补签记录
		attendance := model.H1Attendance{
			ScheduleID:       req.ScheduleID,
			StudentID:        req.StudentID,
			CheckedAt:        time.Now(),
			Status:           "补签",
			MakeupScheduleID: req.MakeupScheduleID,
		}

		if err := query.H1Attendance.WithContext(context.Background()).Create(&attendance); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "补签失败")
			return
		}

		// 补签成功后，更新账单
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ScheduleID)).First()
		if schedule != nil {
			updateStudentBill(req.StudentID, schedule.ClassCourseID)
		}

		response.Success(ctx, attendance, "补签成功")
	} else {
		// 已存在考勤记录，更新为补签状态
		existingAttendance.Status = "补签"
		existingAttendance.MakeupScheduleID = req.MakeupScheduleID
		existingAttendance.CheckedAt = time.Now()

		if err := query.H1Attendance.WithContext(context.Background()).Save(existingAttendance); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "补签失败")
			return
		}

		// 补签成功后，更新账单
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ScheduleID)).First()
		if schedule != nil {
			updateStudentBill(req.StudentID, schedule.ClassCourseID)
		}

		response.Success(ctx, existingAttendance, "补签成功")
	}
}

// MarkAbsent 标记缺勤
// @Summary      标记学生缺勤
// @Description  管理员标记学生缺勤（旷课或请假）
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        absent  body      dto.MarkAbsentRequest  true  "缺勤信息"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance/mark-absent [post]
func MarkAbsent(ctx *gin.Context) {
	var req dto.MarkAbsentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证状态是否有效
	if req.Status != "请假" && req.Status != "旷课" {
		response.Fail(ctx, http.StatusBadRequest, "状态只能是'请假'或'旷课'")
		return
	}

	// 检查是否已有出勤记录
	existingAttendance, _ := query.H1Attendance.WithContext(context.Background()).
		Where(query.H1Attendance.ScheduleID.Eq(req.ScheduleID)).
		Where(query.H1Attendance.StudentID.Eq(req.StudentID)).
		First()

	if existingAttendance != nil {
		// 更新现有记录
		existingAttendance.Status = req.Status
		if err := query.H1Attendance.WithContext(context.Background()).Save(existingAttendance); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "更新失败")
			return
		}

		// 状态变更后，更新账单
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ScheduleID)).First()
		if schedule != nil {
			updateStudentBill(req.StudentID, schedule.ClassCourseID)
		}

		response.Success(ctx, existingAttendance, "标记成功")
	} else {
		// 创建新记录
		attendance := model.H1Attendance{
			ScheduleID: req.ScheduleID,
			StudentID:  req.StudentID,
			CheckedAt:  time.Now(),
			Status:     req.Status,
		}

		if err := query.H1Attendance.WithContext(context.Background()).Create(&attendance); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "标记失败")
			return
		}

		// 标记完成后，更新账单
		schedule, _ := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ScheduleID)).First()
		if schedule != nil {
			updateStudentBill(req.StudentID, schedule.ClassCourseID)
		}

		response.Success(ctx, attendance, "标记成功")
	}
}
