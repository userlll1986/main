package mysqldb_test

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/userlll1986/main/config"
)

var Db *gorm.DB
var err error

// User 结构体声明
// type User struct {
// 	UserId    int64  `gorm:"primaryKey;autoIncrement"`
// 	UserName  string `gorm:"not null;type:varchar(32)"`
// 	UserPwd   string `gorm:"not null;type:varchar(128)"`
// 	UserPhone string `gorm:"unique;type:varchar(32)"`
// }

// 数据库配置
func InitDb(config *config.Config) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Data.Username,
		config.Data.Password,
		config.Data.Ip,
		config.Data.Part,
		config.Data.DataBase,
	)
	// 这个地方要注意，不要写称 :=  写成 = 才对
	Db, err = gorm.Open(config.Data.Category, url)

	// 设置表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		//fmt.Println("db.DefaultTableNameHandler", config.Data.Prefix, defaultTableName)
		return config.Data.Prefix + defaultTableName
	}

	if err != nil {
		log.Fatalf("连接数据库【%s:%s/%s】失败, 失败的原因：【%s】", config.Data.Ip, config.Data.Part, config.Data.DataBase, err)
	}

	// db配置输出SQL语句
	Db.LogMode(config.Data.Sql)
	// 使用表名不适用复数
	Db.SingularTable(true)
	// 连接池配置
	Db.DB().SetMaxOpenConns(20)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetConnMaxLifetime(10 * time.Second)
	// 判断是否需要用来映射结构体到到数据库
	// fmt.Println(config.Data.Init.Status)
	if config.Data.Init.Status {
		// 自动迁移数据库表结构
		// err := Db.AutoMigrate(&User{})
		// if err != nil {
		// 	fmt.Println("数据库表迁移失败!", err)
		// } else {
		// 	fmt.Println("数据库表迁移成功!")
		// }
		// 插入单条数据
		// var user = User{UserName: "wjj", UserPwd: "123", UserPhone: "111"}
		// Db.Create(&user)
		// var users = []User{
		// 	{UserName: "th", UserPwd: "123", UserPhone: "222"},
		// 	{UserName: "lhf", UserPwd: "123", UserPhone: "333"},
		// 	{UserName: "zcy", UserPwd: "123", UserPhone: "444"},
		// }
		// for _, user := range users {
		// 	Db.Create(&user)
		// }
		// 查询全部记录
		var users []mymodals.User
		Db.Find(&users)
		Db.Where("user_name = ?", "wjj").Find(&users)
		// Db.First(&users)
		// // 打印结果
		// fmt.Println(users)
		// 查询总数
		// var users []User
		// var totalSize int64
		// Db.Find(&users).Count(&totalSize)
		// fmt.Println("记录总数:", totalSize)
		// 查询user_id为1的记录
		// var stu User
		// Db.Where("user_id = ?", 1).Find(&stu)
		// // 修改stu姓名为wjj1
		// stu.UserName = "wjj1"
		// // 修改(按照主键修改)
		// Db.Save(&stu)
		// var stu User
		// Db.Model(&stu).Where("user_id = ?", 1).Update("user_name", "wjj2")
		// var fields = map[string]interface{}{"user_name": "WJJ", "user_pwd": "999"}
		// fmt.Println(fields)
		// Db.Model(&stu).Where("user_id = ?", 1).Updates(fields)
		// // 删除
		// var user = User{UserId: 1}
		// Db.Delete(&user)
		// 按照条件删除
		// Db.Where("user_id = ?", 10).Delete(&User{})
	}

	log.Printf("连接数据库【%s:%s/%s】成功", config.Data.Ip, config.Data.Part, config.Data.DataBase)
}
