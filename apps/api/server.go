package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func initRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(allowCORS())

	// 健康检查
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// 各子路由
	registerAuthRoutes(r, db)
	registerMemberRoutes(r, db)
	registerNewsRoutes(r, db)
	registerBannerRoutes(r, db)
	registerMediaRoutes(r, db)
	registerAboutRoutes(r, db)
	registerHighlightRoutes(r, db)

	return r
}
