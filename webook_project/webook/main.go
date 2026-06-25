package main

import (
	"go_learning/webook_project/webook/internal/repository"
	"go_learning/webook_project/webook/internal/repository/dao"
	"go_learning/webook_project/webook/internal/service"
	"go_learning/webook_project/webook/internal/web/middleware"
	"time"

	"go_learning/webook_project/webook/internal/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化DB
	db := initDB()
	//初始化Server
	server := initWebServer()

	//初始化User
	u := initUser(db)
	//注册路由
	u.RegisterRoutes(server)

	//启动
	server.Run(":9090")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
	// CORS 配置，允许前端跨域请求
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	//步骤一
	//store := cookie.NewStore([]byte("your-secret-key"))
	store, err := redis.NewStore(16, "tcp",
		"localhost:6379", "root", "",
		[]byte("qOYZLAuWmwkxAKG6bijwru9ghNNS9rHc"),
		[]byte("8Mv11Olt6x3DX97rUE1exp9XISEMSZJl"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))
	//步骤三 链式调用
	server.Use(middleware.NewLoginMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:040725ge@tcp(localhost:13316)/webook"))
	if err != nil {
		//只在初始化过程中用panic
		//panic相当于整个goroutine结束
		//一旦初始化出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
