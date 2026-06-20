package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	//get方法
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//post方法
	server.POST("/post", func(c *gin.Context) {
		c.String(http.StatusOK, "hello post方法")
	})

	//参数路由
	server.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello 参数路由"+name)
	})

	//查询参数
	server.GET("/order", func(c *gin.Context) {

		oid := c.Query("id")
		c.String(http.StatusOK, "hello 查询参数"+oid)
	})

	//通配符路由
	server.GET("/uviews/*.html", func(c *gin.Context) {
		page := c.Param(".html")
		c.String(http.StatusOK, "hello 通配符路由"+page)
	})

	server.Run() // listens on 0.0.0.0:8080 by default
}
