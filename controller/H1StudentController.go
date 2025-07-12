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

// CreateStudent godoc
// @Summary      Create a new student
// @Description  Create a new student
// @Tags         H1Student
// @Accept       json
// @Produce      json
// @Param        student  body      dto.CreateStudentRequest  true  "Student info"
// @Success      200      {object}  response.Response{data=model.H1Student}
// @Router       /h1/student [post]
func CreateStudent(ctx *gin.Context) {
	var req dto.CreateStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	student := model.H1Student{
		Name:         req.Name,
		Contact:      req.Contact,
		DiscountRate: req.DiscountRate,
	}

	if err := query.H1Student.WithContext(context.Background()).Create(&student); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, student, "创建成功")
}

// GetStudent godoc
// @Summary      Get a student by ID
// @Description  Get a student by ID
// @Tags         H1Student
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  response.Response{data=model.H1Student}
// @Router       /h1/student/{id} [get]
func GetStudent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	student, err := query.H1Student.WithContext(context.Background()).Where(query.H1Student.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生不存在")
		return
	}

	response.Success(ctx, student, "查询成功")
}

// DeleteStudent godoc
// @Summary      Delete a student by ID
// @Description  Delete a student by ID
// @Tags         H1Student
// @Produce      json
// @Param        id   path      int  true  "Student ID"
// @Success      200  {object}  response.Response
// @Router       /h1/student/{id} [delete]
func DeleteStudent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Student.WithContext(context.Background()).Where(query.H1Student.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateStudent godoc
// @Summary      Update a student
// @Description  Update a student
// @Tags         H1Student
// @Accept       json
// @Produce      json
// @Param        student  body      dto.UpdateStudentRequest  true  "Student info"
// @Success      200      {object}  response.Response{data=model.H1Student}
// @Router       /h1/student [put]
func UpdateStudent(ctx *gin.Context) {
	var req dto.UpdateStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	student, err := query.H1Student.WithContext(context.Background()).Where(query.H1Student.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生不存在")
		return
	}

	if req.Name != "" {
		student.Name = req.Name
	}
	if req.Contact != "" {
		student.Contact = req.Contact
	}
	if req.DiscountRate != 0 {
		student.DiscountRate = req.DiscountRate
	}

	if err := query.H1Student.WithContext(context.Background()).Save(student); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, student, "更新成功")
}