package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerBannerRoutes(r *gin.Engine, db *gorm.DB) {
	// 读
	r.GET("/api/v1/banners", func(c *gin.Context) {
		var list []Banner
		db.Order("sort ASC").Order("created_at DESC").Find(&list)
		c.JSON(http.StatusOK, list)
	})

	// 写
	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	authz.POST("/banners", func(c *gin.Context) {
		var b Banner
		if err := c.BindJSON(&b); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&b).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, b)
	})

	authz.DELETE("/banners/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Banner{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
