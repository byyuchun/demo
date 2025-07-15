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
)

// CreateEnrollment godoc
// @Summary      Create a new enrollment
// @Description  Create a new enrollment
// @Tags         H1Enrollment
// @Accept       json
// @Produce      json
// @Param        enrollment  body      dto.CreateEnrollmentRequest  true  "Enrollment info"
// @Success      200      {object}  response.Response{data=model.H1Enrollment}
// @Router       /h1/enrollment [post]
func CreateEnrollment(ctx *gin.Context) {
	var req dto.CreateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	enrollment := model.H1Enrollment{
		StudentID:     req.StudentID,
		ClassCourseID: req.ClassCourseID,
	}

	if err := query.H1Enrollment.WithContext(context.Background()).Create(&enrollment); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, enrollment, "创建成功")
}

// GetEnrollment godoc
// @Summary      Get an enrollment by ID
// @Description  Get an enrollment by ID
// @Tags         H1Enrollment
// @Produce      json
// @Param        id   path      int  true  "Enrollment ID"
// @Success      200  {object}  response.Response{data=model.H1Enrollment}
// @Router       /h1/enrollment/{id} [get]
func GetEnrollment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	enrollment, err := query.H1Enrollment.WithContext(context.Background()).Where(query.H1Enrollment.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "报名记录不存在")
		return
	}

	response.Success(ctx, enrollment, "查询成功")
}

// DeleteEnrollment godoc
// @Summary      Delete an enrollment by ID
// @Description  Delete an enrollment by ID
// @Tags         H1Enrollment
// @Produce      json
// @Param        id   path      int  true  "Enrollment ID"
// @Success      200  {object}  response.Response
// @Router       /h1/enrollment/{id} [delete]
func DeleteEnrollment(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Enrollment.WithContext(context.Background()).Where(query.H1Enrollment.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateEnrollment godoc
// @Summary      Update an enrollment
// @Description  Update an enrollment
// @Tags         H1Enrollment
// @Accept       json
// @Produce      json
// @Param        enrollment  body      dto.UpdateEnrollmentRequest  true  "Enrollment info"
// @Success      200      {object}  response.Response{data=model.H1Enrollment}
// @Router       /h1/enrollment [put]
func UpdateEnrollment(ctx *gin.Context) {
	var req dto.UpdateEnrollmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	enrollment, err := query.H1Enrollment.WithContext(context.Background()).Where(query.H1Enrollment.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "报名记录不存在")
		return
	}

	if req.StudentID != 0 {
		enrollment.StudentID = req.StudentID
	}
	if req.ClassCourseID != 0 {
		enrollment.ClassCourseID = req.ClassCourseID
	}
	if !req.EnrolledAt.IsZero() {
		enrollment.EnrolledAt = req.EnrolledAt
	}

	if err := query.H1Enrollment.WithContext(context.Background()).Save(enrollment); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, enrollment, "更新成功")
}

// GetAllEnrollments godoc
// @Summary      Get all enrollments
// @Description  Get all enrollments
// @Tags         H1Enrollment
// @Produce      json
// @Success      200  {object}  response.Response{data=[]model.H1Enrollment}
// @Router       /h1/enrollment [get]
func GetAllEnrollments(ctx *gin.Context) {
	enrollments, err := query.H1Enrollment.WithContext(context.Background()).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(ctx, enrollments, "查询成功")
}

// GetEnrollmentsWithDetail 获取报名详情列表
// @Summary      获取报名详情列表
// @Description  获取所有报名记录及其关联的学生、班级、课程等详细信息
// @Tags         H1Enrollment
// @Produce      json
// @Param        student_id      query     int     false  "学生ID过滤"
// @Param        semester_id     query     int     false  "学期ID过滤"
// @Param        class_course_id query     int     false  "教学班ID过滤"
// @Param        student_name    query     string  false  "学生姓名模糊查询"
// @Success      200  {object}  response.Response{data=[]dto.EnrollmentWithDetail}
// @Router       /h1/enrollment/with-detail [get]
func GetEnrollmentsWithDetail(ctx *gin.Context) {
	studentIDStr := ctx.Query("student_id")
	semesterIDStr := ctx.Query("semester_id")
	classCourseIDStr := ctx.Query("class_course_id")
	studentName := ctx.Query("student_name")

	// 构建查询
	enrollmentQuery := query.H1Enrollment.WithContext(context.Background())

	// 如果有学生姓名搜索，先查找匹配的学生
	var targetStudentIDs []int64
	if studentName != "" {
		students, err := query.H1Student.WithContext(context.Background()).
			Where(query.H1Student.Name.Like("%" + studentName + "%")).
			Find()
		if err != nil {
			response.Fail(ctx, http.StatusInternalServerError, "查询学生失败")
			return
		}

		if len(students) == 0 {
			// 没有匹配的学生，返回空结果
			response.Success(ctx, []dto.EnrollmentWithDetail{}, "查询成功")
			return
		}

		for _, student := range students {
			targetStudentIDs = append(targetStudentIDs, student.ID)
		}
		enrollmentQuery = enrollmentQuery.Where(query.H1Enrollment.StudentID.In(targetStudentIDs...))
	}

	if studentIDStr != "" {
		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "学生ID格式错误")
			return
		}
		enrollmentQuery = enrollmentQuery.Where(query.H1Enrollment.StudentID.Eq(studentID))
	}

	if classCourseIDStr != "" {
		classCourseID, err := strconv.ParseInt(classCourseIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "教学班ID格式错误")
			return
		}
		enrollmentQuery = enrollmentQuery.Where(query.H1Enrollment.ClassCourseID.Eq(classCourseID))
	}

	enrollments, err := enrollmentQuery.Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.EnrollmentWithDetail
	for _, enrollment := range enrollments {
		// 获取学生信息
		student, _ := query.H1Student.WithContext(context.Background()).
			Where(query.H1Student.ID.Eq(enrollment.StudentID)).First()

		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(enrollment.ClassCourseID)).First()

		if classCourse == nil {
			continue
		}

		// 如果指定了学期过滤，检查教学班是否属于该学期
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

		// 获取班级信息
		class, _ := query.H1Class.WithContext(context.Background()).
			Where(query.H1Class.ID.Eq(classCourse.ClassID)).First()

		// 获取课程信息
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()

		// 获取学期信息
		semester, _ := query.H1Semester.WithContext(context.Background()).
			Where(query.H1Semester.ID.Eq(classCourse.SemesterID)).First()

		detail := dto.EnrollmentWithDetail{
			ID:            enrollment.ID,
			StudentID:     enrollment.StudentID,
			ClassCourseID: enrollment.ClassCourseID,
			EnrolledAt:    enrollment.EnrolledAt,
			StudentName:   "",
			ClassName:     "",
			CourseName:    "",
			SemesterName:  "",
			PricePerHour:  classCourse.PricePerHour,
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
		if semester != nil {
			detail.SemesterName = semester.Name
		}

		result = append(result, detail)
	}

	response.Success(ctx, result, "查询成功")
}

// EnrollStudent 学生报名
// @Summary      学生报名教学班
// @Description  学生报名指定的教学班，会检查重复报名等业务规则
// @Tags         H1Enrollment
// @Accept       json
// @Produce      json
// @Param        enroll  body      dto.EnrollStudentRequest  true  "报名信息"
// @Success      200      {object}  response.Response{data=model.H1Enrollment}
// @Router       /h1/enrollment/enroll [post]
func EnrollStudent(ctx *gin.Context) {
	var req dto.EnrollStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	// 验证学生是否存在
	_, err := query.H1Student.WithContext(context.Background()).
		Where(query.H1Student.ID.Eq(req.StudentID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生不存在")
		return
	}

	// 验证教学班是否存在
	_, err = query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ID.Eq(req.ClassCourseID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "教学班不存在")
		return
	}

	// 检查是否已经报名该教学班
	existingEnrollment, _ := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(req.StudentID)).
		Where(query.H1Enrollment.ClassCourseID.Eq(req.ClassCourseID)).
		First()
	if existingEnrollment != nil {
		response.Fail(ctx, http.StatusBadRequest, "学生已报名该教学班")
		return
	}

	// 创建报名记录
	enrollment := model.H1Enrollment{
		StudentID:     req.StudentID,
		ClassCourseID: req.ClassCourseID,
	}

	if err := query.H1Enrollment.WithContext(context.Background()).Create(&enrollment); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "报名失败")
		return
	}

	response.Success(ctx, enrollment, "报名成功")
}

// GetStudentEnrollments 获取学生的报名记录
// @Summary      获取学生的报名记录
// @Description  获取指定学生的所有报名记录及详细信息
// @Tags         H1Enrollment
// @Produce      json
// @Param        student_id  path      int  true  "学生ID"
// @Success      200  {object}  response.Response{data=[]dto.EnrollmentWithDetail}
// @Router       /h1/enrollment/student/{student_id} [get]
func GetStudentEnrollments(ctx *gin.Context) {
	studentIDStr := ctx.Param("student_id")
	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "学生ID格式错误")
		return
	}

	// 验证学生是否存在
	_, err = query.H1Student.WithContext(context.Background()).
		Where(query.H1Student.ID.Eq(studentID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生不存在")
		return
	}

	// 获取学生的所有报名记录
	enrollments, err := query.H1Enrollment.WithContext(context.Background()).
		Where(query.H1Enrollment.StudentID.Eq(studentID)).
		Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.EnrollmentWithDetail
	for _, enrollment := range enrollments {
		// 获取教学班信息
		classCourse, _ := query.H1ClassCourse.WithContext(context.Background()).
			Where(query.H1ClassCourse.ID.Eq(enrollment.ClassCourseID)).First()
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

		detail := dto.EnrollmentWithDetail{
			ID:            enrollment.ID,
			StudentID:     enrollment.StudentID,
			ClassCourseID: enrollment.ClassCourseID,
			EnrolledAt:    enrollment.EnrolledAt,
			StudentName:   "",
			ClassName:     "",
			CourseName:    "",
			SemesterName:  "",
			PricePerHour:  classCourse.PricePerHour,
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
