package middleware

import (
	"go_learning/webook_project/webook/internal/web"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwt的登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//用JWT进行校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//没登录 返回401
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//按照空格切割
		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			//没登录 有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("qOYZLAuWmwkxAKG6bijwru9ghNNS9rHc"), nil
		})
		if err != nil {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//err为nil，token不为nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//每10秒钟刷新一次 其实相当于生成了一个新的token
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("qOYZLAuWmwkxAKG6bijwru9ghNNS9rHc"))
			if err != nil {
				//记录日志
				log.Println("jwt续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
		ctx.Set("userId", claims.Uid)
	}
}
