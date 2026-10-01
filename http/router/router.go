package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/controller/web"
	"net/http"
)

// XC: 管理后台静态资源缓存策略（部署健壮性）。
// 问题：浏览器对无 Cache-Control 的 index.html 做启发式缓存；部署新版本后旧 index.html
// 引用已删除的旧 chunk（文件名带 hash），SPA 导航报 "Failed to fetch dynamically imported module"。
// 策略：/_admin/static/ 构建产物（hash 文件名）永久缓存；其余（index.html 等）no-cache 每次协商。
func adminStaticCachePolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/_admin/static/") {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			c.Header("Cache-Control", "no-cache")
		}
		c.Next()
	}
}

func WebInit(g *gin.Engine) {
	i := &web.Index{}
	g.GET("/", i.Index)

	if global.Config.App.WebClient == 1 {
		g.GET("/webclient-config/index.js", i.ConfigJs)
	}

	if global.Config.App.WebClient == 1 {
		g.StaticFS("/webclient", http.Dir(global.Config.Gin.ResourcesPath+"/web"))
		g.StaticFS("/webclient2", http.Dir(global.Config.Gin.ResourcesPath+"/web2"))
	}
	// XC: 管理后台入口开关——关闭时 /_admin 静态资源不可访问
	if global.Config.Admin.IsEnabled() {
		// XC: 附带缓存策略中间件（防部署后旧 index.html 引用失效 chunk）
		adminGroup := g.Group("/_admin", adminStaticCachePolicy())
		adminGroup.StaticFS("/", http.Dir(global.Config.Gin.ResourcesPath+"/admin"))
	} else {
		g.Any("/_admin/*any", func(c *gin.Context) {
			c.String(http.StatusForbidden, "Admin panel is disabled")
		})
	}
}
