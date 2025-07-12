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

// CreateStudentSemesterBill godoc
// @Summary      Create a new student semester bill
// @Description  Create a new student semester bill
// @Tags         H1StudentSemesterBill
// @Accept       json
// @Produce      json
// @Param        bill  body      dto.CreateStudentSemesterBillRequest  true  "Bill info"
// @Success      200      {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill [post]
func CreateStudentSemesterBill(ctx *gin.Context) {
	var req dto.CreateStudentSemesterBillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	bill := model.H1StudentSemesterBill{
		StudentID:   req.StudentID,
		SemesterID:  req.SemesterID,
		TotalFee:    req.TotalFee,
		FinalizedAt: time.Now(),
		Status:      req.Status,
	}

	if err := query.H1StudentSemesterBill.WithContext(context.Background()).Create(&bill); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, bill, "创建成功")
}

// GetStudentSemesterBill godoc
// @Summary      Get a student semester bill by ID
// @Description  Get a student semester bill by ID
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        id   path      int  true  "StudentSemesterBill ID"
// @Success      200  {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill/{id} [get]
func GetStudentSemesterBill(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	bill, err := query.H1StudentSemesterBill.WithContext(context.Background()).Where(query.H1StudentSemesterBill.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "学生学期账单不存在")
		return
	}

	response.Success(ctx, bill, "查询成功")
}

// DeleteStudentSemesterBill godoc
// @Summary      Delete a student semester bill by ID
// @Description  Delete a student semester bill by ID
// @Tags         H1StudentSemesterBill
// @Produce      json
// @Param        id   path      int  true  "StudentSemesterBill ID"
// @Success      200  {object}  response.Response
// @Router       /h1/student-semester-bill/{id} [delete]
func DeleteStudentSemesterBill(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1StudentSemesterBill.WithContext(context.Background()).Where(query.H1StudentSemesterBill.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateStudentSemesterBill godoc
// @Summary      Update a student semester bill
// @Description  Update a student semester bill
// @Tags         H1StudentSemesterBill
// @Accept       json
// @Produce      json
// @Param        bill  body      dto.UpdateStudentSemesterBillRequest  true  "Bill info"
// @Success      200      {object}  response.Response{data=model.H1StudentSemesterBill}
// @Router       /h1/student-semester-bill [put]
func UpdateStudentSemesterBill(ctx *gin.Context) {
	var req dto.UpdateStudentSemesterBillRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	bill, err := query.H1StudentSemesterBill.WithContext(context.Background()).Where(query.H1StudentSemesterBill.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "账单不存在")
		return
	}

	if req.StudentID != 0 {
		bill.StudentID = req.StudentID
	}
	if req.SemesterID != 0 {
		bill.SemesterID = req.SemesterID
	}
	if req.TotalFee != 0 {
		bill.TotalFee = req.TotalFee
	}
	if !req.FinalizedAt.IsZero() {
		bill.FinalizedAt = req.FinalizedAt
	}
	if req.Status != "" {
		bill.Status = req.Status
	}

	if err := query.H1StudentSemesterBill.WithContext(context.Background()).Save(bill); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, bill, "更新成功")
}