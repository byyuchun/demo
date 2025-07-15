package controller

import (
	"context"
	"demo/dal/model"
	"demo/dal/query"
	"demo/dto"
	"demo/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetTeachingClassesWithDetail 获取教学班列表（带详细信息）
// @Summary      获取教学班列表
// @Description  获取所有教学班及其详细信息
// @Tags         TeachingClass
// @Produce      json
// @Param        semester_id  query     int  false  "学期ID"
// @Success      200  {object}  response.Response{data=[]dto.TeachingClassWithDetail}
// @Router       /h1/teaching-class [get]
func GetTeachingClassesWithDetail(ctx *gin.Context) {
	semesterIDStr := ctx.Query("semester_id")

	// 构建查询
	classCourseQuery := query.H1ClassCourse.WithContext(context.Background())

	if semesterIDStr != "" {
		semesterID, err := strconv.ParseInt(semesterIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "学期ID格式错误")
			return
		}
		classCourseQuery = classCourseQuery.Where(query.H1ClassCourse.SemesterID.Eq(semesterID))
	}

	classCourses, err := classCourseQuery.Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.TeachingClassWithDetail
	for _, classCourse := range classCourses {
		// 获取班级信息
		class, _ := query.H1Class.WithContext(context.Background()).
			Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()

		// 获取课程信息
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

		// 获取学期信息
		semester, _ := query.H1Semester.WithContext(context.Background()).
			Where(query.H1Semester.ID.Eq(classCourse.SemesterID)).First()

		// 获取学生数量
		studentCount, _ := query.H1Enrollment.WithContext(context.Background()).
			Where(query.H1Enrollment.ClassCourseID.Eq(classCourse.ID)).Count()

		// 获取课程表数量
		scheduleCount, _ := query.H1Schedule.WithContext(context.Background()).
			Where(query.H1Schedule.ClassCourseID.Eq(classCourse.ID)).Count()

		className := ""
		if class != nil {
			className = class.Name
		}

		courseName := ""
		if course != nil {
			courseName = course.Name
		}

		semesterName := ""
		if semester != nil {
			semesterName = semester.Name
		}

		displayName := fmt.Sprintf("%s-%s", className, courseName)

		result = append(result, dto.TeachingClassWithDetail{
			ID:            classCourse.ID,
			ClassID:       classCourse.ClassID,
			CourseID:      classCourse.CourseID,
			SemesterID:    classCourse.SemesterID,
			PricePerHour:  classCourse.PricePerHour,
			ClassName:     className,
			CourseName:    courseName,
			SemesterName:  semesterName,
			StudentCount:  int(studentCount),
			ScheduleCount: int(scheduleCount),
			DisplayName:   displayName,
		})
	}

	response.Success(ctx, result, "查询成功")
}

// CreateTeachingClass 创建教学班
// @Summary      创建教学班
// @Description  创建新的教学班
// @Tags         TeachingClass
// @Accept       json
// @Produce      json
// @Param        teachingClass  body      dto.TeachingClassRequest  true  "教学班信息"
// @Success      200      {object}  response.Response{data=model.H1ClassCourse}
// @Router       /h1/teaching-class [post]
func CreateTeachingClass(ctx *gin.Context) {
	var req dto.TeachingClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 检查是否已存在相同的教学班
	existing, _ := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ClassID.Eq(req.ClassID)).
		Where(query.H1ClassCourse.CourseID.Eq(req.CourseID)).
		Where(query.H1ClassCourse.SemesterID.Eq(req.SemesterID)).
		First()

	if existing != nil {
		response.Fail(ctx, http.StatusBadRequest, "该教学班已存在")
		return
	}

	classCourse := model.H1ClassCourse{
		ClassID:      req.ClassID,
		CourseID:     req.CourseID,
		SemesterID:   req.SemesterID,
		PricePerHour: req.PricePerHour,
	}

	if err := query.H1ClassCourse.WithContext(context.Background()).Create(&classCourse); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, classCourse, "创建成功")
}

// UpdateTeachingClass 更新教学班
// @Summary      更新教学班
// @Description  更新教学班信息
// @Tags         TeachingClass
// @Accept       json
// @Produce      json
// @Param        teachingClass  body      dto.TeachingClassRequest  true  "教学班信息"
// @Success      200      {object}  response.Response{data=model.H1ClassCourse}
// @Router       /h1/teaching-class [put]
func UpdateTeachingClass(ctx *gin.Context) {
	var req dto.TeachingClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	if req.ID == 0 {
		response.Fail(ctx, http.StatusBadRequest, "教学班ID不能为空")
		return
	}

	// 查找要更新的教学班
	classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "教学班不存在")
		return
	}

	// 检查是否与其他教学班冲突
	if req.ClassID != classCourse.ClassID || req.CourseID != classCourse.CourseID || req.SemesterID != classCourse.SemesterID {
		existing, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ClassID.Eq(req.ClassID)).
			Where(query.H1ClassCourse.CourseID.Eq(req.CourseID)).
			Where(query.H1ClassCourse.SemesterID.Eq(req.SemesterID)).
			Where(query.H1ClassCourse.ID.Neq(req.ID)).
			First()

		if existing != nil {
			response.Fail(ctx, http.StatusBadRequest, "该教学班已存在")
			return
		}
	}

	// 更新字段
	updates := map[string]interface{}{
		"class_id":       req.ClassID,
		"course_id":      req.CourseID,
		"semester_id":    req.SemesterID,
		"price_per_hour": req.PricePerHour,
	}

	_, err = query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(req.ID)).Updates(updates)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	// 获取更新后的记录
	updated, _ := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(req.ID)).First()

	response.Success(ctx, updated, "更新成功")
}

// DeleteTeachingClass 删除教学班
// @Summary      删除教学班
// @Description  删除教学班（会同时删除相关的报名记录和课程表）
// @Tags         TeachingClass
// @Produce      json
// @Param        id   path      int  true  "教学班ID"
// @Success      200  {object}  response.Response
// @Router       /h1/teaching-class/{id} [delete]
func DeleteTeachingClass(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "ID格式错误")
		return
	}

	// 检查教学班是否存在
	_, err = query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "教学班不存在")
		return
	}

	// 删除相关的报名记录
	query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.ClassCourseID.Eq(id)).Delete()

	// 删除相关的课程表
	query.H1Schedule.WithContext(context.Background()).
		Where(query.H1Schedule.ClassCourseID.Eq(id)).Delete()

	// 删除教学班
	_, err = query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// GetTeachingClassStudents 获取教学班的学生信息
// @Summary      获取教学班学生
// @Description  获取教学班的已报名学生和可报名学生列表
// @Tags         TeachingClass
// @Produce      json
// @Param        id   path      int  true  "教学班ID"
// @Success      200  {object}  response.Response{data=dto.TeachingClassStudent}
// @Router       /h1/teaching-class/{id}/students [get]
func GetTeachingClassStudents(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "ID格式错误")
		return
	}

	// 获取教学班信息
	classCourse, err := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "教学班不存在")
		return
	}

	// 获取班级和课程信息用于显示名称
	class, _ := query.H1Class.WithContext(context.Background()).
		Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()
	course, _ := query.H1Course.WithContext(context.Background()).
		Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

	displayName := fmt.Sprintf("%s-%s", class.Name, course.Name)

	// 获取已报名学生
	enrollments, _ := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.ClassCourseID.Eq(id)).Find()

	var enrolledStudents []dto.EnrolledStudent
	var enrolledStudentIDs []int64

	for _, enrollment := range enrollments {
		student, _ := query.H1Student.WithContext(context.Background()).
			Where(query.H1Student.ID.Eq(enrollment.StudentID)).First()

		if student != nil {
			enrolledStudents = append(enrolledStudents, dto.EnrolledStudent{
				ID:           enrollment.ID,
				StudentID:    student.ID,
				StudentName:  student.Name,
				Contact:      student.Contact,
				DiscountRate: student.DiscountRate,
				EnrolledAt:   enrollment.EnrolledAt,
			})
			enrolledStudentIDs = append(enrolledStudentIDs, student.ID)
		}
	}

	// 获取可报名学生（未报名此教学班的学生）
	allStudents, _ := query.H1Student.WithContext(context.Background()).Find()
	var availableStudents []dto.StudentInfo

	for _, student := range allStudents {
		// 检查是否已报名
		isEnrolled := false
		for _, enrolledID := range enrolledStudentIDs {
			if student.ID == enrolledID {
				isEnrolled = true
				break
			}
		}

		if !isEnrolled {
			availableStudents = append(availableStudents, dto.StudentInfo{
				ID:           student.ID,
				Name:         student.Name,
				Contact:      student.Contact,
				DiscountRate: student.DiscountRate,
			})
		}
	}

	result := dto.TeachingClassStudent{
		ClassCourseID:     id,
		DisplayName:       displayName,
		EnrolledStudents:  enrolledStudents,
		AvailableStudents: availableStudents,
	}

	response.Success(ctx, result, "查询成功")
}

// BatchEnrollStudents 批量报名学生
// @Summary      批量报名学生
// @Description  将多个学生批量报名到教学班
// @Tags         TeachingClass
// @Accept       json
// @Produce      json
// @Param        request  body      dto.BatchEnrollRequest  true  "批量报名请求"
// @Success      200      {object}  response.Response
// @Router       /h1/teaching-class/batch-enroll [post]
func BatchEnrollStudents(ctx *gin.Context) {
	var req dto.BatchEnrollRequest
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

	var successCount int
	var failureReasons []string

	for _, studentID := range req.StudentIDs {
		// 验证学生是否存在
		_, err := query.H1Student.WithContext(context.Background()).
			Where(query.H1Student.ID.Eq(studentID)).First()
		if err != nil {
			failureReasons = append(failureReasons, fmt.Sprintf("学生ID %d 不存在", studentID))
			continue
		}

		// 检查是否已经报名
		existing, _ := query.H1Enrollment.WithContext(context.Background()).
			Where(query.H1Enrollment.StudentID.Eq(studentID)).
			Where(query.H1Enrollment.ClassCourseID.Eq(req.ClassCourseID)).
			First()

		if existing != nil {
			failureReasons = append(failureReasons, fmt.Sprintf("学生ID %d 已报名此教学班", studentID))
			continue
		}

		// 创建报名记录
		enrollment := model.H1Enrollment{
			StudentID:     studentID,
			ClassCourseID: req.ClassCourseID,
		}

		if err := query.H1Enrollment.WithContext(context.Background()).Create(&enrollment); err != nil {
			failureReasons = append(failureReasons, fmt.Sprintf("学生ID %d 报名失败: %v", studentID, err))
			continue
		}

		successCount++
	}

	message := fmt.Sprintf("批量报名完成，成功: %d，失败: %d", successCount, len(failureReasons))
	if len(failureReasons) > 0 {
		message += "，失败原因: " + fmt.Sprintf("%v", failureReasons)
	}

	response.Success(ctx, gin.H{
		"success_count":   successCount,
		"failure_count":   len(failureReasons),
		"failure_reasons": failureReasons,
	}, message)
}

// BatchUnenrollStudents 批量退课
// @Summary      批量退课
// @Description  将多个学生从教学班中退课
// @Tags         TeachingClass
// @Accept       json
// @Produce      json
// @Param        request  body      dto.BatchUnenrollRequest  true  "批量退课请求"
// @Success      200      {object}  response.Response
// @Router       /h1/teaching-class/batch-unenroll [post]
func BatchUnenrollStudents(ctx *gin.Context) {
	var req dto.BatchUnenrollRequest
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

	var successCount int
	var failureReasons []string

	for _, studentID := range req.StudentIDs {
		// 删除报名记录
		result, err := query.H1Enrollment.WithContext(context.Background()).
			Where(query.H1Enrollment.StudentID.Eq(studentID)).
			Where(query.H1Enrollment.ClassCourseID.Eq(req.ClassCourseID)).
			Delete()

		if err != nil {
			failureReasons = append(failureReasons, fmt.Sprintf("学生ID %d 退课失败: %v", studentID, err))
			continue
		}

		if result.RowsAffected == 0 {
			failureReasons = append(failureReasons, fmt.Sprintf("学生ID %d 未报名此教学班", studentID))
			continue
		}

		successCount++
	}

	message := fmt.Sprintf("批量退课完成，成功: %d，失败: %d", successCount, len(failureReasons))
	if len(failureReasons) > 0 {
		message += "，失败原因: " + fmt.Sprintf("%v", failureReasons)
	}

	response.Success(ctx, gin.H{
		"success_count":   successCount,
		"failure_count":   len(failureReasons),
		"failure_reasons": failureReasons,
	}, message)
}
