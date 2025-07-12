package routes

import (
	"demo/controller"
	"demo/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/admin/register", controller.ARegister)

	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/admin/login", controller.ALogin)

	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.GET("/api/admin/info", middleware.AdminMiddleware(), controller.AdminInfo)

	r.PUT("/api/admin/:id", middleware.AdminMiddleware(), controller.AUpdate)
	r.DELETE("/api/admin/:id", middleware.AdminMiddleware(), controller.ADelete)
	r.GET("/api/admin/show", middleware.AdminMiddleware(), controller.AShow)
	r.GET("/api/admin/search", middleware.AdminMiddleware(), controller.ASearch)

	r.PUT("/api/auth/:id", middleware.AuthMiddleware(), controller.PassUpdate)
	r.POST("/api/auth/:id", controller.InfoUpdate)
	r.DELETE("/api/auth/:id", middleware.AdminMiddleware(), controller.Delete)

	r.POST("/api/message", controller.ForgetPass)
	r.GET("/api/message/show", middleware.AdminMiddleware(), controller.ShowMsg)
	r.DELETE("/api/message/:sid", middleware.AdminMiddleware(), controller.DeleteMsg)

	r.POST("/api/class/add", middleware.AdminMiddleware(), controller.AddClass)
	r.GET("/api/class/show", middleware.AdminMiddleware(), controller.ShowClass)
	r.DELETE("/api/class/:name", middleware.AdminMiddleware(), controller.DeleteClass)
	r.POST("/api/class/:name", middleware.AdminMiddleware(), controller.UpdateClass)

	h1Group := r.Group("/h1")
	{
		studentGroup := h1Group.Group("/student")
		{
			studentGroup.POST("", controller.CreateStudent)
			studentGroup.PUT("", controller.UpdateStudent)
			studentGroup.GET("/:id", controller.GetStudent)
			studentGroup.DELETE("/:id", controller.DeleteStudent)
		}

		semesterGroup := h1Group.Group("/semester")
		{
			semesterGroup.POST("", controller.CreateSemester)
			semesterGroup.PUT("", controller.UpdateSemester)
			semesterGroup.GET("/:id", controller.GetSemester)
			semesterGroup.DELETE("/:id", controller.DeleteSemester)
		}

		courseGroup := h1Group.Group("/course")
		{
			courseGroup.POST("", controller.CreateCourse)
			courseGroup.PUT("", controller.UpdateCourse)
			courseGroup.GET("/:id", controller.GetCourse)
			courseGroup.DELETE("/:id", controller.DeleteCourse)
		}

		classGroup := h1Group.Group("/class")
		{
			classGroup.POST("", controller.CreateH1Class)
			classGroup.PUT("", controller.UpdateH1Class)
			classGroup.GET("/:id", controller.GetH1Class)
			classGroup.DELETE("/:id", controller.DeleteH1Class)
		}

		classCourseGroup := h1Group.Group("/class-course")
		{
			classCourseGroup.POST("", controller.CreateClassCourse)
			classCourseGroup.PUT("", controller.UpdateClassCourse)
			classCourseGroup.GET("/:id", controller.GetClassCourse)
			classCourseGroup.DELETE("/:id", controller.DeleteClassCourse)
		}

		enrollmentGroup := h1Group.Group("/enrollment")
		{
			enrollmentGroup.POST("", controller.CreateEnrollment)
			enrollmentGroup.PUT("", controller.UpdateEnrollment)
			enrollmentGroup.GET("/:id", controller.GetEnrollment)
			enrollmentGroup.DELETE("/:id", controller.DeleteEnrollment)
		}

		scheduleGroup := h1Group.Group("/schedule")
		{
			scheduleGroup.POST("", controller.CreateSchedule)
			scheduleGroup.PUT("", controller.UpdateSchedule)
			scheduleGroup.GET("/:id", controller.GetSchedule)
			scheduleGroup.DELETE("/:id", controller.DeleteSchedule)
		}

		attendanceGroup := h1Group.Group("/attendance")
		{
			attendanceGroup.POST("", controller.CreateAttendance)
			attendanceGroup.PUT("", controller.UpdateAttendance)
			attendanceGroup.GET("/:id", controller.GetAttendance)
			attendanceGroup.DELETE("/:id", controller.DeleteAttendance)
		}

		billGroup := h1Group.Group("/student-semester-bill")
		{
			billGroup.POST("", controller.CreateStudentSemesterBill)
			billGroup.PUT("", controller.UpdateStudentSemesterBill)
			billGroup.GET("/:id", controller.GetStudentSemesterBill)
			billGroup.DELETE("/:id", controller.DeleteStudentSemesterBill)
		}
	}

	return r
}
