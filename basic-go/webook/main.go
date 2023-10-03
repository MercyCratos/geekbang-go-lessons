package main

import (
	"geekbang-lessons/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	hdl := web.NewUserHandler()

	server := gin.Default()
	hdl.RegisterRoutes(server)

	server.Run(":8080")
}
