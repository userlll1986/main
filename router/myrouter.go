package myrouter

import (
	"github.com/gin-gonic/gin"
)

func ArticleRouter(group *gin.RouterGroup) {
	// 文章相关的接口
	article := group.Group("/article")
	{
		// 获取文章列表
		article.GET("/:id", v1.GetOneUser)
	}
}

// 初始化路由
func InitRouter(model string) *gin.Engine {
	gin.SetMode(model)
	r := gin.Default()
	groupV1 := r.Group("api/v1")
	// 文章相关的接口组
	ArticleRouter(groupV1)
	return r
}
