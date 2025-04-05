package mymodals

// User 结构体声明
type User struct {
	UserId    int64  `gorm:"primaryKey;autoIncrement"`
	UserName  string `gorm:"not null;type:varchar(32)"`
	UserPwd   string `gorm:"not null;type:varchar(128)"`
	UserPhone string `gorm:"unique;type:varchar(32)"`
}
