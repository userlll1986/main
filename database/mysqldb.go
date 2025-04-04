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
	if config.Data.Init.Status {
		Db.AutoMigrate(&User{})
	}

	log.Printf("连接数据库【%s:%s/%s】成功", config.Data.Ip, config.Data.Part, config.Data.DataBase)
}
