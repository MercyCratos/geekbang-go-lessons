package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 声明一个engine
	server := gin.Default()

	// ======================== 注册路由 ======================== //
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello, world")
	})
	server.GET("users/:name", func(ctx *gin.Context) {
		// 获取 PathVariable
		name := ctx.Param("name")
		// 获取查询参数
		age := ctx.Query("age")
		ctx.String(http.StatusOK, "这是你传过来的名字[%s]和年龄[%s]", name, age)
	})
	// 通配符路由匹配
	server.GET("views/*.html", func(ctx *gin.Context) {
		path := ctx.Param(".html")
		ctx.String(http.StatusOK, "匹配上的值是 %s", path)
	})

	// 启动engine并指定监听的端口
	// 如果不传入参数，实际上监听的是 8080 端口
	server.Run(":8080")
}
