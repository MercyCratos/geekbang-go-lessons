package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}

		session := sessions.Default(ctx)
		userId := session.Get("userId")
		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		// 注册一下该类型，gob 才能知道如何序列化 Time
		gob.Register(now)

		const updateTimeKey = "update_time"
		val := session.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Minute {
			session.Set(updateTimeKey, now)
			session.Set("userId", userId)
			session.Options(sessions.Options{
				MaxAge:   900,
				HttpOnly: true,
			})
			err := session.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
