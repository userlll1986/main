package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/userlll1986/main/config"
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

	// Define handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Listen and serve on defined port
	log.Printf("Listening on port %s", port)
	r.Run(":" + port)
}
