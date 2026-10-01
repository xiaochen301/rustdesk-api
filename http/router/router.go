package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/controller/web"
	"net/http"
)

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
		g.StaticFS("/_admin", http.Dir(global.Config.Gin.ResourcesPath+"/admin"))
	} else {
		g.Any("/_admin/*any", func(c *gin.Context) {
			c.String(http.StatusForbidden, "Admin panel is disabled")
		})
	}
}
