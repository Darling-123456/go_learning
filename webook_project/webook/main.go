package main

import (
	"time"

	"go_learning/webook_project/webook/internal/web"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	// CORS 配置，允许前端跨域请求
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
