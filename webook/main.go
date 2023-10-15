package main

import (
	"geekbang-lessons/webook/internal/repository"
	"geekbang-lessons/webook/internal/repository/dao"
	"geekbang-lessons/webook/internal/service"
	"geekbang-lessons/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()

	server := initWebServer()

	initUserHandler(db, server)

	server.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}))
	return server
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	userHandler.RegisterRoutes(server)
}
