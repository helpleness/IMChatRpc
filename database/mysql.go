package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 表单的全局变量
var DB *gorm.DB

// 初始化数据库
func InitMysql() *gorm.DB {
	//得到yml中数据库配置信息，yml由viper管理
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	database := viper.GetString("mysql.database")
	//正则表达式，用于创建mysql连接
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		username, password, host, port, database)
	//db存储数据库信息，err存储报错
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql database,err: " + err.Error())
	}
	//自动创建表单，如果已经有表单，则不重复创建

	//db.AutoMigrate(&model.User{})
	//db.AutoMigrate(&model.File{})
	//db.AutoMigrate(&model.Friends{})
	//db.AutoMigrate(&model.Group{})
	//db.AutoMigrate(&model.GroupMember{})
	//db.AutoMigrate(&model.MyMessage{})
	//db.AutoMigrate(&model.FriendAdd{})
	//db.AutoMigrate(&model.GroupApplication{})
	DB = db
	return db
}

// 返回DB的函数
func GetDB() *gorm.DB {
	return DB
}
