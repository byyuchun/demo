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

// CreateCourse godoc
// @Summary      Create a new course
// @Description  Create a new course
// @Tags         H1Course
// @Accept       json
// @Produce      json
// @Param        course  body      dto.CreateCourseRequest  true  "Course info"
// @Success      200      {object}  response.Response{data=model.H1Course}
// @Router       /h1/course [post]
func CreateCourse(ctx *gin.Context) {
	var req dto.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	course := model.H1Course{
		Name: req.Name,
	}

	if err := query.H1Course.WithContext(context.Background()).Create(&course); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, course, "创建成功")
}

// GetCourse godoc
// @Summary      Get a course by ID
// @Description  Get a course by ID
// @Tags         H1Course
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  response.Response{data=model.H1Course}
// @Router       /h1/course/{id} [get]
func GetCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	course, err := query.H1Course.WithContext(context.Background()).Where(query.H1Course.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课程不存在")
		return
	}

	response.Success(ctx, course, "查询成功")
}

// DeleteCourse godoc
// @Summary      Delete a course by ID
// @Description  Delete a course by ID
// @Tags         H1Course
// @Produce      json
// @Param        id   path      int  true  "Course ID"
// @Success      200  {object}  response.Response
// @Router       /h1/course/{id} [delete]
func DeleteCourse(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Course.WithContext(context.Background()).Where(query.H1Course.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateCourse godoc
// @Summary      Update a course
// @Description  Update a course
// @Tags         H1Course
// @Accept       json
// @Produce      json
// @Param        course  body      dto.UpdateCourseRequest  true  "Course info"
// @Success      200      {object}  response.Response{data=model.H1Course}
// @Router       /h1/course [put]
func UpdateCourse(ctx *gin.Context) {
	var req dto.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	course, err := query.H1Course.WithContext(context.Background()).Where(query.H1Course.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课程不存在")
		return
	}

	if req.Name != "" {
		course.Name = req.Name
	}

	if err := query.H1Course.WithContext(context.Background()).Save(course); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, course, "更新成功")
}