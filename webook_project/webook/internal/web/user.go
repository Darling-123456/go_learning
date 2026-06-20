package web

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

const (
	emailPattern    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	passwordPattern = `^[a-zA-Z0-9!@#$%^&*]{8,}$`
)

// 在userhandler中定义和用户有关的路由
type UserHandler struct {
	emailCompiled    *regexp.Regexp
	passwordCompiled *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		emailCompiled:    regexp.MustCompile(emailPattern),
		passwordCompiled: regexp.MustCompile(passwordPattern),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {

	//不分组的写法
	//注册
	server.POST("/users/signup", u.SignUp)
	//登录
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

	ctx.String(http.StatusOK, "注册成功")

	//数据库操作

}

// 登录
func (u *UserHandler) Login(ctx *gin.Context) {

}

// 编辑
func (u *UserHandler) Edit(ctx *gin.Context) {

}

// 用户信息
func (u *UserHandler) Profile(ctx *gin.Context) {

}
