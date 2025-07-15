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

// CreateSchedule godoc
// @Summary      Create a new schedule
// @Description  Create a new schedule
// @Tags         H1Schedule
// @Accept       json
// @Produce      json
// @Param        schedule  body      dto.CreateScheduleRequest  true  "Schedule info"
// @Success      200      {object}  response.Response{data=model.H1Schedule}
// @Router       /h1/schedule [post]
func CreateSchedule(ctx *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	schedule := model.H1Schedule{
		ClassCourseID: req.ClassCourseID,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	}

	if err := query.H1Schedule.WithContext(context.Background()).Create(&schedule); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, schedule, "创建成功")
}

// GetSchedule godoc
// @Summary      Get a schedule by ID
// @Description  Get a schedule by ID
// @Tags         H1Schedule
// @Produce      json
// @Param        id   path      int  true  "Schedule ID"
// @Success      200  {object}  response.Response{data=model.H1Schedule}
// @Router       /h1/schedule/{id} [get]
func GetSchedule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	schedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课程表不存在")
		return
	}

	response.Success(ctx, schedule, "查询成功")
}

// DeleteSchedule godoc
// @Summary      Delete a schedule by ID
// @Description  Delete a schedule by ID
// @Tags         H1Schedule
// @Produce      json
// @Param        id   path      int  true  "Schedule ID"
// @Success      200  {object}  response.Response
// @Router       /h1/schedule/{id} [delete]
func DeleteSchedule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateSchedule godoc
// @Summary      Update a schedule
// @Description  Update a schedule
// @Tags         H1Schedule
// @Accept       json
// @Produce      json
// @Param        schedule  body      dto.UpdateScheduleRequest  true  "Schedule info"
// @Success      200      {object}  response.Response{data=model.H1Schedule}
// @Router       /h1/schedule [put]
func UpdateSchedule(ctx *gin.Context) {
	var req dto.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	schedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课程表不存在")
		return
	}

	if req.ClassCourseID != 0 {
		schedule.ClassCourseID = req.ClassCourseID
	}
	if !req.Date.IsZero() {
		schedule.Date = req.Date
	}
	if !req.StartTime.IsZero() {
		schedule.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		schedule.EndTime = req.EndTime
	}

	if err := query.H1Schedule.WithContext(context.Background()).Save(schedule); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, schedule, "更新成功")
}

// GetAllSchedules godoc
// @Summary      Get all schedules
// @Description  Get all schedules
// @Tags         H1Schedule
// @Produce      json
// @Success      200  {object}  response.Response{data=[]model.H1Schedule}
// @Router       /h1/schedule [get]
func GetAllSchedules(ctx *gin.Context) {
	schedules, err := query.H1Schedule.WithContext(context.Background()).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(ctx, schedules, "查询成功")
}

// GetSchedulesWithDetail 获取课表详情列表
// @Summary      获取课表详情列表
// @Description  获取课表及其关联的班级、课程等详细信息
// @Tags         H1Schedule
// @Produce      json
// @Param        class_course_id  query     int     false  "教学班ID过滤"
// @Param        date            query     string  false  "日期过滤 (YYYY-MM-DD)"
// @Param        start_date      query     string  false  "开始日期 (YYYY-MM-DD)"
// @Param        end_date        query     string  false  "结束日期 (YYYY-MM-DD)"
// @Success      200  {object}  response.Response{data=[]dto.ScheduleWithDetail}
// @Router       /h1/schedule/with-detail [get]
func GetSchedulesWithDetail(ctx *gin.Context) {
	classCourseIDStr := ctx.Query("class_course_id")
	dateStr := ctx.Query("date")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	// 构建查询
	scheduleQuery := query.H1Schedule.WithContext(context.Background())

	if classCourseIDStr != "" {
		classCourseID, err := strconv.ParseInt(classCourseIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "教学班ID格式错误")
			return
		}
		scheduleQuery = scheduleQuery.Where(query.H1Schedule.ClassCourseID.Eq(classCourseID))
	}

	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "日期格式错误，请使用 YYYY-MM-DD")
			return
		}
		scheduleQuery = scheduleQuery.Where(query.H1Schedule.Date.Eq(date))
	}

	if startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			response.Fail(ctx, http.StatusBadRequest, "日期格式错误，请使用 YYYY-MM-DD")
			return
		}
		scheduleQuery = scheduleQuery.Where(query.H1Schedule.Date.Between(startDate, endDate))
	}

	schedules, err := scheduleQuery.Order(query.H1Schedule.Date, query.H1Schedule.StartTime).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.ScheduleWithDetail
	for _, schedule := range schedules {
		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(schedule.ClassCourseID)).First()

		if classCourse == nil {
			continue
		}

		// 获取班级信息
		class, _ := query.H1Class.WithContext(context.Background()).
			Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()
		// 获取课程信息
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()
		// 获取学期信息
		semester, _ := query.H1Semester.WithContext(context.Background()).
			Where(query.H1Semester.ID.Eq(classCourse.SemesterID)).First()

		detail := dto.ScheduleWithDetail{
			ID:            schedule.ID,
			ClassCourseID: schedule.ClassCourseID,
			Date:          schedule.Date,
			StartTime:     schedule.StartTime,
			EndTime:       schedule.EndTime,
			ClassName:     "",
			CourseName:    "",
			SemesterName:  "",
		}

		if class != nil {
			detail.ClassName = class.Name
		}
		if course != nil {
			detail.CourseName = course.Name
		}
		if semester != nil {
			detail.SemesterName = semester.Name
		}

		result = append(result, detail)
	}

	response.Success(ctx, result, "查询成功")
}

// CreateScheduleBatch 批量创建课表
// @Summary      批量创建课表
// @Description  为教学班批量创建多天的课表
// @Tags         H1Schedule
// @Accept       json
// @Produce      json
// @Param        batch  body      dto.CreateScheduleBatchRequest  true  "批量课表信息"
// @Success      200      {object}  response.Response{data=[]model.H1Schedule}
// @Router       /h1/schedule/batch [post]
func CreateScheduleBatch(ctx *gin.Context) {
	var req dto.CreateScheduleBatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证教学班是否存在
	_, err := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(req.ClassCourseID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "教学班不存在")
		return
	}

	var schedules []model.H1Schedule
	current := req.StartDate

	for !current.After(req.EndDate) {
		// 根据星期规则判断是否需要排课
		weekday := current.Weekday()
		shouldCreateSchedule := false

		for _, day := range req.WeekDays {
			if int(weekday) == day {
				shouldCreateSchedule = true
				break
			}
		}

		if shouldCreateSchedule {
			// 检查是否已经存在该日期的课表
			existingSchedule, _ := query.H1Schedule.WithContext(context.Background()).
				Where(query.H1Schedule.ClassCourseID.Eq(req.ClassCourseID)).
				Where(query.H1Schedule.Date.Eq(current)).
				First()

			if existingSchedule == nil {
				// 创建时间戳 - 使用日期加上时间
				startTime := time.Date(current.Year(), current.Month(), current.Day(),
					req.StartHour, req.StartMinute, 0, 0, current.Location())
				endTime := time.Date(current.Year(), current.Month(), current.Day(),
					req.EndHour, req.EndMinute, 0, 0, current.Location())

				schedule := model.H1Schedule{
					ClassCourseID: req.ClassCourseID,
					Date:          current,
					StartTime:     startTime,
					EndTime:       endTime,
				}
				schedules = append(schedules, schedule)
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	// 批量创建
	if len(schedules) > 0 {
		// 转换为指针切片
		schedulePointers := make([]*model.H1Schedule, len(schedules))
		for i := range schedules {
			schedulePointers[i] = &schedules[i]
		}

		if err := query.H1Schedule.WithContext(context.Background()).CreateInBatches(schedulePointers, 100); err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "批量创建失败")
			return
		}
	}

	response.Success(ctx, schedules, "批量创建成功")
}

// GetTodaySchedules 获取今日课表
// @Summary      获取今日课表
// @Description  获取今天的所有课表安排
// @Tags         H1Schedule
// @Produce      json
// @Success      200  {object}  response.Response{data=[]dto.ScheduleWithDetail}
// @Router       /h1/schedule/today [get]
func GetTodaySchedules(ctx *gin.Context) {
	today := time.Now().Format("2006-01-02")
	todayTime, _ := time.Parse("2006-01-02", today)

	schedules, err := query.H1Schedule.WithContext(context.Background()).
		Where(query.H1Schedule.Date.Eq(todayTime)).
		Order(query.H1Schedule.StartTime).
		Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.ScheduleWithDetail
	for _, schedule := range schedules {
		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(schedule.ClassCourseID)).First()

		if classCourse == nil {
			continue
		}

		// 获取班级和课程信息
		class, _ := query.H1Class.WithContext(context.Background()).
			Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()
		semester, _ := query.H1Semester.WithContext(context.Background()).
			Where(query.H1Semester.ID.Eq(classCourse.SemesterID)).First()

		detail := dto.ScheduleWithDetail{
			ID:            schedule.ID,
			ClassCourseID: schedule.ClassCourseID,
			Date:          schedule.Date,
			StartTime:     schedule.StartTime,
			EndTime:       schedule.EndTime,
			ClassName:     "",
			CourseName:    "",
			SemesterName:  "",
		}

		if class != nil {
			detail.ClassName = class.Name
		}
		if course != nil {
			detail.CourseName = course.Name
		}
		if semester != nil {
			detail.SemesterName = semester.Name
		}

		result = append(result, detail)
	}

	response.Success(ctx, result, "查询成功")
}

// GetStudentSchedules 获取学生的课表
// @Summary      获取学生的课表
// @Description  获取指定学生的课表安排
// @Tags         H1Schedule
// @Produce      json
// @Param        student_id   query     int     true   "学生ID"
// @Param        start_date   query     string  false  "开始日期 (YYYY-MM-DD)"
// @Param        end_date     query     string  false  "结束日期 (YYYY-MM-DD)"
// @Success      200  {object}  response.Response{data=[]dto.StudentSchedule}
// @Router       /h1/schedule/student [get]
func GetStudentSchedules(ctx *gin.Context) {
	studentIDStr := ctx.Query("student_id")
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "学生ID格式错误")
		return
	}

	// 获取学生的所有报名记录
	enrollments, err := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(studentID)).
		Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "获取报名记录失败")
		return
	}

	var classCourseIDs []int64
	for _, enrollment := range enrollments {
		classCourseIDs = append(classCourseIDs, enrollment.ClassCourseID)
	}

	if len(classCourseIDs) == 0 {
		response.Success(ctx, []dto.StudentSchedule{}, "查询成功")
		return
	}

	// 构建课表查询
	scheduleQuery := query.H1Schedule.WithContext(context.Background()).
		Where(query.H1Schedule.ClassCourseID.In(classCourseIDs...))

	if startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			response.Fail(ctx, http.StatusBadRequest, "日期格式错误，请使用 YYYY-MM-DD")
			return
		}
		scheduleQuery = scheduleQuery.Where(query.H1Schedule.Date.Between(startDate, endDate))
	}

	schedules, err := scheduleQuery.Order(query.H1Schedule.Date, query.H1Schedule.StartTime).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询课表失败")
		return
	}

	var result []dto.StudentSchedule
	for _, schedule := range schedules {
		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(schedule.ClassCourseID)).First()
		if classCourse == nil {
			continue
		}

		// 获取班级和课程信息
		class, _ := query.H1Class.WithContext(context.Background()).
			Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

		// 检查出勤记录
		attendance, _ := query.H1Attendance.WithContext(context.Background()).
			Where(query.H1Attendance.ScheduleID.Eq(schedule.ID)).
			Where(query.H1Attendance.StudentID.Eq(studentID)).
			First()

		studentSchedule := dto.StudentSchedule{
			ScheduleID:       schedule.ID,
			ClassCourseID:    schedule.ClassCourseID,
			Date:             schedule.Date,
			StartTime:        schedule.StartTime,
			EndTime:          schedule.EndTime,
			ClassName:        "",
			CourseName:       "",
			AttendanceStatus: "未打卡",
		}

		if class != nil {
			studentSchedule.ClassName = class.Name
		}
		if course != nil {
			studentSchedule.CourseName = course.Name
		}
		if attendance != nil {
			studentSchedule.AttendanceStatus = attendance.Status
		}

		result = append(result, studentSchedule)
	}

	response.Success(ctx, result, "查询成功")
}
