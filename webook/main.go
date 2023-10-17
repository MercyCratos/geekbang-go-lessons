package main

import (
	"geekbang-lessons/webook/internal/repository"
	"geekbang-lessons/webook/internal/repository/dao"
	"geekbang-lessons/webook/internal/service"
	"geekbang-lessons/webook/internal/web"
	"geekbang-lessons/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
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

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
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
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 允许前端访问你的后端响应中带的哪些头部
		ExposeHeaders: []string{"X-Auth-Token"},
		MaxAge:        12 * time.Hour,
	}))

	//useSession(server)
	useJWT(server)

	return server
}

func useJWT(server *gin.Engine) {
	loginMiddleware := &middleware.LoginJWTMiddlewareBuilder{}
	server.Use(loginMiddleware.CheckLogin())
}

func useSession(server *gin.Engine) {
	loginMiddleware := &middleware.LoginMiddlewareBuilder{}
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK"),
		[]byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgA"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("ssid", store), loginMiddleware.CheckLogin())
}

func initUserHandler(db *gorm.DB, server *gin.Engine) {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	userHandler.RegisterRoutes(server)
}
