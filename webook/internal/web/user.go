package web

import (
	"geekbang-lessons/webook/internal/domain"
	"geekbang-lessons/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"unicode/utf8"
)

const (
	emailRegexPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
)

type UserHandler struct {
	emailRegExp    *regexp.Regexp
	passwordRegExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	userGroup := server.Group("/users")

	userGroup.POST("/signup", h.signup)
	userGroup.POST("/login", h.login)
	userGroup.PUT("/edit", h.edit)
	userGroup.GET("/profile", h.profile)
}

func (h *UserHandler) signup(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req SignupReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailRegExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	isPassword, err := h.passwordRegExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个邮箱")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		session := sessions.Default(ctx)
		session.Set("userId", u.Id)
		session.Options(sessions.Options{
			MaxAge:   900,
			HttpOnly: true,
		})
		err = session.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}

		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) profile(ctx *gin.Context) {
	userId := h.getUserIdFromSession(ctx)

	u, err := h.svc.GetProfile(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	type UserVO struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		AboutMe  string `json:"about_me"`
		Birthday string `json:"birthday"`
	}
	ctx.JSON(http.StatusOK, UserVO{
		Nickname: u.Nickname,
		Email:    u.Email,
		AboutMe:  u.AboutMe,
		Birthday: u.Birthday.Format(time.DateOnly),
	})
}

func (h *UserHandler) edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	userId := h.getUserIdFromSession(ctx)

	if utf8.RuneCountInString(req.Nickname) > 24 {
		ctx.String(http.StatusOK, "昵称不能超过24位")
		return
	}

	if utf8.RuneCountInString(req.AboutMe) > 128 {
		ctx.String(http.StatusOK, "关于我不能超过128个字符")
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}

	err = h.svc.UpdateUserInfo(ctx, domain.User{
		Id:       userId,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "更新成功")
}

func (h *UserHandler) getUserIdFromSession(ctx *gin.Context) int64 {
	session := sessions.Default(ctx)
	return session.Get("userId").(int64)
}
