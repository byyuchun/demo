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

	h1Group := r.Group("/api/h1")
	h1Group.Use(middleware.AdminMiddleware()) // 添加管理员认证中间件
	{
		studentGroup := h1Group.Group("/student")
		{
			studentGroup.POST("", controller.CreateStudent)
			studentGroup.PUT("", controller.UpdateStudent)
			studentGroup.GET("", controller.GetAllStudents)
			studentGroup.GET("/:id", controller.GetStudent)
			studentGroup.DELETE("/:id", controller.DeleteStudent)
		}

		semesterGroup := h1Group.Group("/semester")
		{
			semesterGroup.POST("", controller.CreateSemester)
			semesterGroup.PUT("", controller.UpdateSemester)
			semesterGroup.GET("", controller.GetAllSemesters)
			semesterGroup.GET("/:id", controller.GetSemester)
			semesterGroup.DELETE("/:id", controller.DeleteSemester)
		}

		courseGroup := h1Group.Group("/course")
		{
			courseGroup.POST("", controller.CreateCourse)
			courseGroup.PUT("", controller.UpdateCourse)
			courseGroup.GET("", controller.GetAllCourses)
			courseGroup.GET("/:id", controller.GetCourse)
			courseGroup.DELETE("/:id", controller.DeleteCourse)
		}

		classGroup := h1Group.Group("/class")
		{
			classGroup.POST("", controller.CreateH1Class)
			classGroup.PUT("", controller.UpdateH1Class)
			classGroup.GET("", controller.GetAllH1Classes)
			classGroup.GET("/:id", controller.GetH1Class)
			classGroup.DELETE("/:id", controller.DeleteH1Class)
			// 新增业务接口
			classGroup.GET("/with-semester", controller.GetClassesWithSemester)
			classGroup.GET("/:id/courses", controller.GetClassCourses)
		}

		classCourseGroup := h1Group.Group("/class-course")
		{
			classCourseGroup.POST("", controller.CreateClassCourse)
			classCourseGroup.PUT("", controller.UpdateClassCourse)
			classCourseGroup.GET("", controller.GetAllClassCourses)
			classCourseGroup.GET("/:id", controller.GetClassCourse)
			classCourseGroup.DELETE("/:id", controller.DeleteClassCourse)
		}

		// 教学班管理路由
		teachingClassGroup := h1Group.Group("/teaching-class")
		{
			teachingClassGroup.GET("", controller.GetTeachingClassesWithDetail)
			teachingClassGroup.POST("", controller.CreateTeachingClass)
			teachingClassGroup.PUT("", controller.UpdateTeachingClass)
			teachingClassGroup.DELETE("/:id", controller.DeleteTeachingClass)
			teachingClassGroup.GET("/:id/students", controller.GetTeachingClassStudents)
			teachingClassGroup.POST("/batch-enroll", controller.BatchEnrollStudents)
			teachingClassGroup.POST("/batch-unenroll", controller.BatchUnenrollStudents)
		}

		enrollmentGroup := h1Group.Group("/enrollment")
		{
			enrollmentGroup.POST("", controller.CreateEnrollment)
			enrollmentGroup.PUT("", controller.UpdateEnrollment)
			enrollmentGroup.GET("", controller.GetAllEnrollments)
			enrollmentGroup.GET("/:id", controller.GetEnrollment)
			enrollmentGroup.DELETE("/:id", controller.DeleteEnrollment)
			// 新增业务接口
			enrollmentGroup.GET("/with-detail", controller.GetEnrollmentsWithDetail)
			enrollmentGroup.POST("/enroll", controller.EnrollStudent)
			enrollmentGroup.GET("/student/:student_id", controller.GetStudentEnrollments)
		}

		scheduleGroup := h1Group.Group("/schedule")
		{
			scheduleGroup.POST("", controller.CreateSchedule)
			scheduleGroup.PUT("", controller.UpdateSchedule)
			scheduleGroup.GET("", controller.GetAllSchedules)
			scheduleGroup.GET("/:id", controller.GetSchedule)
			scheduleGroup.DELETE("/:id", controller.DeleteSchedule)
			// 新增业务接口
			scheduleGroup.GET("/with-detail", controller.GetSchedulesWithDetail)
			scheduleGroup.POST("/batch", controller.CreateScheduleBatch)
			scheduleGroup.GET("/today", controller.GetTodaySchedules)
			scheduleGroup.GET("/student", controller.GetStudentSchedules)
		}

		attendanceGroup := h1Group.Group("/attendance")
		{
			attendanceGroup.POST("", controller.CreateAttendance)
			attendanceGroup.PUT("", controller.UpdateAttendance)
			attendanceGroup.GET("", controller.GetAllAttendances)
			attendanceGroup.GET("/:id", controller.GetAttendance)
			attendanceGroup.DELETE("/:id", controller.DeleteAttendance)
			// 新增业务接口
			attendanceGroup.POST("/check-in", controller.CheckIn)
			attendanceGroup.POST("/apply-makeup", controller.ApplyMakeup)
			attendanceGroup.POST("/admin-makeup", controller.AdminMakeup)
			attendanceGroup.GET("/student", controller.GetStudentAttendance)
			attendanceGroup.POST("/mark-absent", controller.MarkAbsent)
		}

		billGroup := h1Group.Group("/student-semester-bill")
		{
			billGroup.POST("", controller.CreateStudentSemesterBill)
			billGroup.PUT("", controller.UpdateStudentSemesterBill)
			billGroup.GET("", controller.GetAllStudentSemesterBills)
			billGroup.GET("/:id", controller.GetStudentSemesterBill)
			billGroup.DELETE("/:id", controller.DeleteStudentSemesterBill)
			// 新增业务接口
			billGroup.POST("/generate", controller.GenerateSemesterBill)
			billGroup.GET("/detail", controller.GetBillDetail)
			billGroup.GET("/statistics", controller.GetBillStatistics)
			billGroup.POST("/pay", controller.PayBill)
			billGroup.GET("/all-with-detail", controller.GetAllBillsWithDetail)
		}
	}

	return r
}
