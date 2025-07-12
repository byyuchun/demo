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

// CreateAttendance godoc
// @Summary      Create a new attendance record
// @Description  Create a new attendance record
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        attendance  body      dto.CreateAttendanceRequest  true  "Attendance info"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance [post]
func CreateAttendance(ctx *gin.Context) {
	var req dto.CreateAttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	attendance := model.H1Attendance{
		ScheduleID:       req.ScheduleID,
		StudentID:        req.StudentID,
		CheckedAt:        time.Now(),
		Status:           req.Status,
		MakeupScheduleID: req.MakeupScheduleID,
	}

	if err := query.H1Attendance.WithContext(context.Background()).Create(&attendance); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, attendance, "创建成功")
}

// GetAttendance godoc
// @Summary      Get an attendance by ID
// @Description  Get an attendance by ID
// @Tags         H1Attendance
// @Produce      json
// @Param        id   path      int  true  "Attendance ID"
// @Success      200  {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance/{id} [get]
func GetAttendance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "ID格式错误")
		return
	}

	attendance, err := query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "考勤记录不存在")
		return
	}

	response.Success(ctx, attendance, "查询成功")
}

// DeleteAttendance godoc
// @Summary      Delete an attendance by ID
// @Description  Delete an attendance by ID
// @Tags         H1Attendance
// @Produce      json
// @Param        id   path      int  true  "Attendance ID"
// @Success      200  {object}  response.Response
// @Router       /h1/attendance/{id} [delete]
func DeleteAttendance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "ID格式错误")
		return
	}

	_, err = query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateAttendance godoc
// @Summary      Update an attendance record
// @Description  Update an attendance record
// @Tags         H1Attendance
// @Accept       json
// @Produce      json
// @Param        attendance  body      dto.UpdateAttendanceRequest  true  "Attendance info"
// @Success      200      {object}  response.Response{data=model.H1Attendance}
// @Router       /h1/attendance [put]
func UpdateAttendance(ctx *gin.Context) {
	var req dto.UpdateAttendanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	attendance, err := query.H1Attendance.WithContext(context.Background()).Where(query.H1Attendance.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "考勤记录不存在")
		return
	}

	if req.ScheduleID != 0 {
		attendance.ScheduleID = req.ScheduleID
	}
	if req.StudentID != 0 {
		attendance.StudentID = req.StudentID
	}
	if !req.CheckedAt.IsZero() {
		attendance.CheckedAt = req.CheckedAt
	}
	if req.Status != "" {
		attendance.Status = req.Status
	}
	if req.MakeupScheduleID != 0 {
		attendance.MakeupScheduleID = req.MakeupScheduleID
	}

	if err := query.H1Attendance.WithContext(context.Background()).Save(attendance); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, attendance, "更新成功")
}
