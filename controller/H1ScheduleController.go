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

// CreateSchedule godoc
// @Summary      Create a new schedule
// @Description  Create a new schedule
// @Tags         H1Schedule
// @Accept       json
// @Produce      json
// @Param        schedule  body      dto.CreateScheduleRequest  true  "Schedule info"
// @Success      200      {object}  response.Response{data=model.H1Schedule}
// @Router       /h1/schedule [post]
func CreateSchedule(ctx *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	schedule := model.H1Schedule{
		ClassCourseID: req.ClassCourseID,
		Date:          req.Date,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
	}

	if err := query.H1Schedule.WithContext(context.Background()).Create(&schedule); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败")
		return
	}

	response.Success(ctx, schedule, "创建成功")
}

// GetSchedule godoc
// @Summary      Get a schedule by ID
// @Description  Get a schedule by ID
// @Tags         H1Schedule
// @Produce      json
// @Param        id   path      int  true  "Schedule ID"
// @Success      200  {object}  response.Response{data=model.H1Schedule}
// @Router       /h1/schedule/{id} [get]
func GetSchedule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	schedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(id)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课程表不存在")
		return
	}

	response.Success(ctx, schedule, "查询成功")
}

// DeleteSchedule godoc
// @Summary      Delete a schedule by ID
// @Description  Delete a schedule by ID
// @Tags         H1Schedule
// @Produce      json
// @Param        id   path      int  true  "Schedule ID"
// @Success      200  {object}  response.Response
// @Router       /h1/schedule/{id} [delete]
func DeleteSchedule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	_, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(id)).Delete()
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(ctx, nil, "删除成功")
}

// UpdateSchedule godoc
// @Summary      Update a schedule
// @Description  Update a schedule
// @Tags         H1Schedule
// @Accept       json
// @Produce      json
// @Param        schedule  body      dto.UpdateScheduleRequest  true  "Schedule info"
// @Success      200      {object}  response.Response{data=model.H1Schedule}
// @Router       /h1/schedule [put]
func UpdateSchedule(ctx *gin.Context) {
	var req dto.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	schedule, err := query.H1Schedule.WithContext(context.Background()).Where(query.H1Schedule.ID.Eq(req.ID)).First()
	if err != nil {
		response.Fail(ctx, http.StatusNotFound, "课程表不存在")
		return
	}

	if req.ClassCourseID != 0 {
		schedule.ClassCourseID = req.ClassCourseID
	}
	if !req.Date.IsZero() {
		schedule.Date = req.Date
	}
	if !req.StartTime.IsZero() {
		schedule.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		schedule.EndTime = req.EndTime
	}

	if err := query.H1Schedule.WithContext(context.Background()).Save(schedule); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "更新失败")
		return
	}

	response.Success(ctx, schedule, "更新成功")
}