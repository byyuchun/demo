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