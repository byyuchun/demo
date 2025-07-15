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

// CreateH1Class godoc
// @Summary      Create a new class
// @Description  Create a new class
// @Tags         H1Class
// @Accept       json
// @Produce      json
// @Param        class  body      dto.CreateClassRequest  true  "Class info"
// @Success      200      {object}  response.Response{data=model.H1Class}
// @Router       /h1/class [post]
func CreateH1Class(ctx *gin.Context) {
	var req dto.CreateClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	class := model.H1Class{
		Name:       req.Name,
		SemesterID: req.SemesterID,
	}

	if err := query.H1Class.WithContext(context.Background()).Create(&class); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, class, "创建成功")
}

// GetH1Class godoc
// @Summary      Get a class by ID
// @Description  Get a class by ID
// @Tags         H1Class
// @Produce      json
// @Param        id   path      int  true  "Class ID"
// @Success      200  {object}  response.Response{data=model.H1Class}
// @Router       /h1/class/{id} [get]
func GetH1Class(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	class, err := query.H1Class.WithContext(context.Background()).Where(query.H1Class.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "班级不存在")
		return
	}

	response.Success(ctx, class, "查询成功")
}

// DeleteH1Class godoc
// @Summary      Delete a class by ID
// @Description  Delete a class by ID
// @Tags         H1Class
// @Produce      json
// @Param        id   path      int  true  "Class ID"
// @Success      200  {object}  response.Response
// @Router       /h1/class/{id} [delete]
func DeleteH1Class(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Class.WithContext(context.Background()).Where(query.H1Class.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateH1Class godoc
// @Summary      Update a class
// @Description  Update a class
// @Tags         H1Class
// @Accept       json
// @Produce      json
// @Param        class  body      dto.UpdateClassRequest  true  "Class info"
// @Success      200      {object}  response.Response{data=model.H1Class}
// @Router       /h1/class [put]
func UpdateH1Class(ctx *gin.Context) {
	var req dto.UpdateClassRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	class, err := query.H1Class.WithContext(context.Background()).Where(query.H1Class.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "班级不存在")
		return
	}

	if req.Name != "" {
		class.Name = req.Name
	}
	if req.SemesterID != 0 {
		class.SemesterID = req.SemesterID
	}

	if err := query.H1Class.WithContext(context.Background()).Save(class); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, class, "更新成功")
}

// GetAllH1Classes godoc
// @Summary      Get all classes
// @Description  Get all classes
// @Tags         H1Class
// @Produce      json
// @Success      200  {object}  response.Response{data=[]model.H1Class}
// @Router       /h1/class [get]
func GetAllH1Classes(ctx *gin.Context) {
	classes, err := query.H1Class.WithContext(context.Background()).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(ctx, classes, "查询成功")
}

// GetClassesWithSemester 获取班级及学期信息
// @Summary      获取班级及关联的学期信息
// @Description  获取所有班级及其关联的学期信息
// @Tags         H1Class
// @Produce      json
// @Param        semester_id  query     int  false  "学期ID过滤"
// @Success      200  {object}  response.Response{data=[]dto.ClassWithSemester}
// @Router       /h1/class/with-semester [get]
func GetClassesWithSemester(ctx *gin.Context) {
	semesterIDStr := ctx.Query("semester_id")

	// 构建查询
	classQuery := query.H1Class.WithContext(context.Background())
	if semesterIDStr != "" {
		semesterID, err := strconv.ParseInt(semesterIDStr, 10, 64)
		if err != nil {
			response.Fail(ctx, http.StatusBadRequest, "学期ID格式错误")
			return
		}
		classQuery = classQuery.Where(query.H1Class.SemesterID.Eq(semesterID))
	}

	classes, err := classQuery.Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.ClassWithSemester
	for _, class := range classes {
		// 获取学期信息
		semester, _ := query.H1Semester.WithContext(context.Background()).
			Where(query.H1Semester.ID.Eq(class.SemesterID)).First()

		classWithSemester := dto.ClassWithSemester{
			ID:           class.ID,
			Name:         class.Name,
			SemesterID:   class.SemesterID,
			SemesterName: "",
		}

		if semester != nil {
			classWithSemester.SemesterName = semester.Name
		}

		result = append(result, classWithSemester)
	}

	response.Success(ctx, result, "查询成功")
}

// GetClassCourses 获取班级的课程列表
// @Summary      获取班级开设的课程列表
// @Description  获取指定班级开设的所有课程及教学班信息
// @Tags         H1Class
// @Produce      json
// @Param        class_id  path      int  true  "班级ID"
// @Success      200  {object}  response.Response{data=[]dto.ClassCourseWithInfo}
// @Router       /h1/class/{id}/courses [get]
func GetClassCourses(ctx *gin.Context) {
	classIDStr := ctx.Param("id")
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "班级ID格式错误")
		return
	}

	// 获取该班级的所有教学班
	classCourses, err := query.H1ClassCourse.WithContext(context.Background()).
		Where(query.H1ClassCourse.ClassID.Eq(classID)).
		Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	var result []dto.ClassCourseWithInfo
	for _, classCourse := range classCourses {
		// 获取课程信息
		course, _ := query.H1Course.WithContext(context.Background()).
			Where(query.H1Course.ID.Eq(classCourse.CourseID)).First()
		// 获取学期信息
		semester, _ := query.H1Semester.WithContext(context.Background()).
			Where(query.H1Semester.ID.Eq(classCourse.SemesterID)).First()

		info := dto.ClassCourseWithInfo{
			ID:           classCourse.ID,
			ClassID:      classCourse.ClassID,
			CourseID:     classCourse.CourseID,
			SemesterID:   classCourse.SemesterID,
			PricePerHour: classCourse.PricePerHour,
			CourseName:   "",
			SemesterName: "",
		}

		if course != nil {
			info.CourseName = course.Name
		}
		if semester != nil {
			info.SemesterName = semester.Name
		}

		result = append(result, info)
	}

	response.Success(ctx, result, "查询成功")
}
