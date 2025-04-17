package account

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	mymodals "main/modals"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func accouont_login(c *gin.Context, Type string, Tag string, Body map[string]interface{}) {
	// 这里进行实际的登录验证逻辑
	log.Printf("登录请求: %v,%v,%v", Type, Tag, Body)
	// 例如：检查用户名和密码是否匹配，验证码是否正确等
	var users []mymodals.User
	var db = c.MustGet("db").(*gorm.DB)
	db.Where("user_name = ?", Body["username"].(string)).Find(&users)

	var resp mymodals.AccountServiceResponse
	resp.Type = "AccountService"
	resp.Tag = "account_login"

	if users != nil {
		// 验证密码是否正确
		hashedPassword := hashPassword(Body["password"].(string))
		log.Printf("hashedPassword: %s", hashedPassword)
		if users[0].UserPwd == Body["password"].(string) {
			// 登录成功
			// c.JSON(http.StatusOK, gin.H{"result": "ok"})
			resp.Result = "ok"
			resp.Body.Length = 0
			// return
		} else {
			resp.Result = "error"
			resp.Body.Length = 1
			resp.Body.Detail = "用户名或密码错误"
		}
	} else {
		resp.Result = "error"
		resp.Body.Length = 1
		resp.Body.Detail = "用户名或密码错误"
	}

	// 将响应数据发送给客户端
	c.JSON(http.StatusOK, resp)
}

// 定义几个示例函数
func foo() {
	fmt.Println("Running foo")
}

func bar(a, b int) int {
	return a + b
}

func baz(s string) string {
	return "Hello, " + s
}
