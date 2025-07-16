package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")

		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost:8082",     // 本地前端开发服务器
			"http://127.0.0.1:8082",     // 本地前端开发服务器
			"http://118.25.157.30:8082", // 服务器前端地址
			"http://118.25.157.30",      // 服务器根地址（如果前端部署在80端口）
		}

		// 检查请求的origin是否在允许列表中
		allowOrigin := ""
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				allowOrigin = origin
				break
			}
		}

		// 如果没有匹配的origin，默认允许localhost开发环境
		if allowOrigin == "" && (strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1")) {
			allowOrigin = origin
		}

		// 如果还是没有匹配，但是是服务器IP，也允许
		if allowOrigin == "" && strings.Contains(origin, "118.25.157.30") {
			allowOrigin = origin
		}

		if allowOrigin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}

		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}
