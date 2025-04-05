package main

import (
	"fmt"
	"log"
	mysqldb_test "main/database"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/userlll1986/main/config"
)

var (
	limiter = NewLimiter(10, 1*time.Minute) // 设置限流器，允许每分钟最多请求10次
)

// NewLimiter 创建限流器
func NewLimiter(limit int, duration time.Duration) *Limiter {
	return &Limiter{
		limit:      limit,
		duration:   duration,
		timestamps: make(map[string][]int64),
	}
}

// Limiter 限流器
type Limiter struct {
	limit      int                // 限制的请求数量
	duration   time.Duration      // 时间窗口
	timestamps map[string][]int64 // 请求的时间戳
}

// Middleware 限流中间件
func (l *Limiter) Middleware(c *gin.Context) {
	ip := c.ClientIP() // 获取客户端IP地址

	// 检查请求时间戳切片是否存在
	if _, ok := l.timestamps[ip]; !ok {
		l.timestamps[ip] = make([]int64, 0)
	}

	now := time.Now().Unix() // 当前时间戳

	// 移除过期的请求时间戳
	for i := 0; i < len(l.timestamps[ip]); i++ {
		if l.timestamps[ip][i] < now-int64(l.duration.Seconds()) {
			l.timestamps[ip] = append(l.timestamps[ip][:i], l.timestamps[ip][i+1:]...)
			i--
		}
	}

	// 检查请求数量是否超过限制
	if len(l.timestamps[ip]) >= l.limit {
		c.JSON(429, gin.H{
			"message": "Too Many Requests",
		})
		c.Abort()
		return
	}

	// 添加当前请求时间戳到切片
	l.timestamps[ip] = append(l.timestamps[ip], now)

	// 继续处理请求
	c.Next()
}

// corsMiddleware 返回CORS中间件处理函数
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有的跨域请求
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400") // 预检请求缓存时间，单位为秒

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// 继续处理其他请求
		c.Next()
	}
}

// main 函数是应用程序的入口点。它启动一个 HTTP 服务器，使用 Gin 框架处理请求。
// 服务器监听环境变量 PORT 指定的端口(默认为 8080)，并提供以下 API 端点:
//   - GET /: 返回 "Hello World!"
//   - GET /ping: 返回 "pong"
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// 读取配置文件
	myconfig := config.NewConfig()
	myconfig.ReadConfig()
	// 初始化数据库连接
	mysqldb_test.InitDb(myconfig)

	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// 使用 Recovery 中间件
	r.Use(gin.Recovery())

	// 使用 Logger 中间件
	r.Use(gin.Logger())

	// 使用 CORSMiddleware 中间件
	r.Use(corsMiddleware())

	// 使用限流中间件
	r.Use(limiter.Middleware)

	//使用数据库中间件
	// 将db作为中间件传递给路由处理函数
	r.Use(func(c *gin.Context) {
		c.Set("db", mysqldb_test.Db)
		c.Next()
	})
	// 在路由处理函数中可以通过c.MustGet("db").(*gorm.DB)获取到db对象，然后进行数据库操作

	// Define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		//截取/
		action = strings.Trim(action, "/")
		c.String(http.StatusOK, name+" is "+action)
	})
	//处理表单数据
	//<form action="http://localhost:8080/form" method="post" action="application/x-www-form-urlencoded">
	// 	用户名：<input type="text" name="username" placeholder="请输入你的用户名">  <br>
	// 	密&nbsp;&nbsp;&nbsp;码：<input type="password" name="userpassword" placeholder="请输入你的密码">  <br>
	// 	<input type="submit" value="提交">
	// </form>
	r.POST("/form", func(c *gin.Context) {
		types := c.DefaultPostForm("type", "post")
		username := c.PostForm("username")
		password := c.PostForm("userpassword")
		// c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
		c.String(http.StatusOK, fmt.Sprintf("username:%s,password:%s,type:%s", username, password, types))
	})
	//处理文件上传
	// <form action="http://localhost:8080/upload" method="post" enctype="multipart/form-data">
	//       上传文件:<input type="file" name="file" >
	//       <input type="submit" value="提交">
	// </form>
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.String(500, "上传图片出错")
		}
		// c.JSON(200, gin.H{"message": file.Header.Context})
		c.SaveUploadedFile(file, file.Filename)
		c.String(http.StatusOK, file.Filename)
	})

	//上传多个文件
	//<form action="http://localhost:8000/uploads" method="post" enctype="multipart/form-data">
	// 	上传文件:<input type="file" name="files" multiple>
	// 	<input type="submit" value="提交">
	// </form>
	// 限制表单上传大小 8MB，默认为32MB
	r.MaxMultipartMemory = 8 << 20
	r.POST("/uploads", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		// 获取所有图片
		files := form.File["files"]
		// 遍历所有图片
		for _, file := range files {
			// 逐个存
			if err := c.SaveUploadedFile(file, file.Filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}
		}
		c.String(200, fmt.Sprintf("upload ok %d files", len(files)))
	})

	// 路由组1 ，处理GET请求
	v1 := r.Group("/v1")
	// {} 是书写规范
	{
		v1.GET("/login", login)
		v1.GET("submit", submit)
	}
	// 路由组2 ，处理POST请求
	v2 := r.Group("/v2")
	{
		v2.POST("/login", login)
		v2.POST("/submit", submit)
	}

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port)
	r.Run(":" + port)
}

func login(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}

func submit(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n", name))
}
