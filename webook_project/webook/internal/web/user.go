package web

import (
	"fmt"
	"go_learning/webook_project/webook/internal/domain"
	"go_learning/webook_project/webook/internal/service"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	emailPattern    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	passwordPattern = `^[a-zA-Z0-9!@#$%^&*]{8,}$`
)

// 在userhandler中定义和用户有关的路由
type UserHandler struct {
	svc              *service.UserService
	emailCompiled    *regexp.Regexp
	passwordCompiled *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		emailCompiled:    regexp.MustCompile(emailPattern),
		passwordCompiled: regexp.MustCompile(passwordPattern),
	}
}

func (u *UserHandler) RegisterRoutesV1(server *gin.Engine) {

	//不分组的写法
	//注册
	server.POST("/users/signup", u.SignUp)
	//登录
	//server.POST("/users/login", u.Login)
	server.POST("/users/login", u.Login)
	//编辑
	server.POST("/users/edit", u.Edit)
	//用户信息
	server.GET("/users/profile", u.Profile)
	/*//分组
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)*/
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {

	//不分组的写法
	//注册
	server.POST("/users/signup", u.SignUp)
	//登录
	//server.POST("/users/login", u.Login)
	server.POST("/users/login", u.LoginJWT)
	//编辑
	server.POST("/users/edit", u.Edit)
	//用户信息
	server.GET("/users/profile", u.ProfileJWT)
	/*//分组
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)*/
}

// 注册
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		fmt.Println("绑定失败:", err)
		ctx.String(http.StatusBadRequest, "请求格式错误，请检查输入")
		return
	}
	fmt.Printf("收到注册请求: %+v\n", req)

	// 校验邮箱
	if !u.emailCompiled.MatchString(req.Email) {
		ctx.String(http.StatusBadRequest, "邮箱格式不对")
		return
	}

	// 校验密码
	if !u.passwordCompiled.MatchString(req.Password) {
		ctx.String(http.StatusBadRequest, "密码必须不少于8位")
		return
	}

	// 校验确认密码
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusBadRequest, "两次输入密码不一致")
		return
	}

	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")

	//数据库操作

}

// JWT实现的登录
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	//jwt登录
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("qOYZLAuWmwkxAKG6bijwru9ghNNS9rHc"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	fmt.Println(user)
	ctx.String(http.StatusOK, "登录成功")
}

// Gin Session插件实现的登录
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//在这里登陆成功了,设置session
	sess := sessions.Default(ctx)
	//可以随便设置值了
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		//半小时过期
		MaxAge:   30 * 60,
		Secure:   true,
		HttpOnly: true,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
}

// 退出登录
func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	//可以随便设置值了
	sess.Options(sessions.Options{MaxAge: -1})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录")
}

// 编辑
func (u *UserHandler) Edit(ctx *gin.Context) {

}

// 用户信息
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	//可以断定必然有claims
	if !ok {
		//考虑监控住这里
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	//类型断言 ok代表是不是*UserClaims
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	println(claims.Uid)
	//后面就是补充profile的其它代码
}

// 用户信息
func (u *UserHandler) Profile(ctx *gin.Context) {
	ctx.String(http.StatusOK, "这是你的用户信息")
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
