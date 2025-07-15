package common

import (
	"demo/model"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Admin{})
	db.AutoMigrate(&model.Info{})
	db.AutoMigrate(&model.Msg{})
	db.AutoMigrate(&model.School{})
	db.AutoMigrate(&model.Class{})
	db.AutoMigrate(&model.Major{})

	// 自动迁移H1表
	//db.AutoMigrate(&dalmodel.H1Student{})
	//db.AutoMigrate(&dalmodel.H1Semester{})
	//db.AutoMigrate(&dalmodel.H1Course{})
	//db.AutoMigrate(&dalmodel.H1Class{})
	//db.AutoMigrate(&dalmodel.H1ClassCourse{})
	//db.AutoMigrate(&dalmodel.H1Enrollment{})
	//db.AutoMigrate(&dalmodel.H1Schedule{})
	//db.AutoMigrate(&dalmodel.H1Attendance{})
	//db.AutoMigrate(&dalmodel.H1StudentSemesterBill{})

	DB = db
	return db
}

// GetDB 定义一个方法获取DB实例
func GetDB() *gorm.DB {
	return DB
}
