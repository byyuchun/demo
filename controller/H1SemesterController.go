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

// CreateSemester godoc
// @Summary      Create a new semester
// @Description  Create a new semester
// @Tags         H1Semester
// @Accept       json
// @Produce      json
// @Param        semester  body      dto.CreateSemesterRequest  true  "Semester info"
// @Success      200      {object}  response.Response{data=model.H1Semester}
// @Router       /h1/semester [post]
func CreateSemester(ctx *gin.Context) {
	var req dto.CreateSemesterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	semester := model.H1Semester{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := query.H1Semester.WithContext(context.Background()).Create(&semester); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, semester, "创建成功")
}

// GetSemester godoc
// @Summary      Get a semester by ID
// @Description  Get a semester by ID
// @Tags         H1Semester
// @Produce      json
// @Param        id   path      int  true  "Semester ID"
// @Success      200  {object}  response.Response{data=model.H1Semester}
// @Router       /h1/semester/{id} [get]
func GetSemester(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	semester, err := query.H1Semester.WithContext(context.Background()).Where(query.H1Semester.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学期不存在")
		return
	}

	response.Success(ctx, semester, "查询成功")
}

// DeleteSemester godoc
// @Summary      Delete a semester by ID
// @Description  Delete a semester by ID
// @Tags         H1Semester
// @Produce      json
// @Param        id   path      int  true  "Semester ID"
// @Success      200  {object}  response.Response
// @Router       /h1/semester/{id} [delete]
func DeleteSemester(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Semester.WithContext(context.Background()).Where(query.H1Semester.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateSemester godoc
// @Summary      Update a semester
// @Description  Update a semester
// @Tags         H1Semester
// @Accept       json
// @Produce      json
// @Param        semester  body      dto.UpdateSemesterRequest  true  "Semester info"
// @Success      200      {object}  response.Response{data=model.H1Semester}
// @Router       /h1/semester [put]
func UpdateSemester(ctx *gin.Context) {
	var req dto.UpdateSemesterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	semester, err := query.H1Semester.WithContext(context.Background()).Where(query.H1Semester.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学期不存在")
		return
	}

	if req.Name != "" {
		semester.Name = req.Name
	}
	if !req.StartDate.IsZero() {
		semester.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		semester.EndDate = req.EndDate
	}

	if err := query.H1Semester.WithContext(context.Background()).Save(semester); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, semester, "更新成功")
}