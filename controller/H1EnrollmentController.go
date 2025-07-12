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