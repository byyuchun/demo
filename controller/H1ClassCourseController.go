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

// CreateClassCourse godoc
// @Summary      Create a new class course
// @Description  Create a new class course
// @Tags         H1ClassCourse
// @Accept       json
// @Produce      json
// @Param        class_course  body      dto.CreateClassCourseRequest  true  "ClassCourse info"
// @Success      200      {object}  response.Response{data=model.H1ClassCourse}
// @Router       /h1/class-course [post]
func CreateClassCourse(ctx *gin.Context) {
	var req dto.CreateClassCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
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

// GetClassCourse godoc
// @Summary      Get a class course by ID
// @Description  Get a class course by ID
// @Tags         H1ClassCourse
// @Produce      json
// @Param        id   path      int  true  "ClassCourse ID"
// @Success      200  {object}  response.Response{data=model.H1ClassCourse}
// @Router       /h1/class-course/{id} [get]
func GetClassCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	classCourse, err := query.H1ClassCourse.WithContext(context.Background()).Where(query.H1ClassCourse.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "班级课程不存在")
		return
	}

	response.Success(ctx, classCourse, "查询成功")
}

// DeleteClassCourse godoc
// @Summary      Delete a class course by ID
// @Description  Delete a class course by ID
// @Tags         H1ClassCourse
// @Produce      json
// @Param        id   path      int  true  "ClassCourse ID"
// @Success      200  {object}  response.Response
// @Router       /h1/class-course/{id} [delete]
func DeleteClassCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1ClassCourse.WithContext(context.Background()).Where(query.H1ClassCourse.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateClassCourse godoc
// @Summary      Update a class course
// @Description  Update a class course
// @Tags         H1ClassCourse
// @Accept       json
// @Produce      json
// @Param        class_course  body      dto.UpdateClassCourseRequest  true  "ClassCourse info"
// @Success      200      {object}  response.Response{data=model.H1ClassCourse}
// @Router       /h1/class-course [put]
func UpdateClassCourse(ctx *gin.Context) {
	var req dto.UpdateClassCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	classCourse, err := query.H1ClassCourse.WithContext(context.Background()).Where(query.H1ClassCourse.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "班级课程不存在")
		return
	}

	if req.ClassID != 0 {
		classCourse.ClassID = req.ClassID
	}
	if req.CourseID != 0 {
		classCourse.CourseID = req.CourseID
	}
	if req.SemesterID != 0 {
		classCourse.SemesterID = req.SemesterID
	}
	if req.PricePerHour != 0 {
		classCourse.PricePerHour = req.PricePerHour
	}

	if err := query.H1ClassCourse.WithContext(context.Background()).Save(classCourse); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, classCourse, "更新成功")
}

// GetAllClassCourses godoc
// @Summary      Get all class courses
// @Description  Get all class courses
// @Tags         H1ClassCourse
// @Produce      json
// @Success      200  {object}  response.Response{data=[]model.H1ClassCourse}
// @Router       /h1/class-course [get]
func GetAllClassCourses(ctx *gin.Context) {
	classCourses, err := query.H1ClassCourse.WithContext(context.Background()).Find()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(ctx, classCourses, "查询成功")
}
