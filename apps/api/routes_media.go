package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerMediaRoutes(r *gin.Engine, db *gorm.DB) {
	// 读
	r.GET("/api/v1/news_images", func(c *gin.Context) {
		var list []NewsImage
		db.Order("created_at DESC").Find(&list)
		c.JSON(http.StatusOK, list)
	})

	// 写
	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	authz.POST("/news_images", func(c *gin.Context) {
		var ni NewsImage
		if err := c.BindJSON(&ni); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if ni.ImageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image_url required"})
			return
		}
		if err := db.Create(&ni).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, ni)
	})

	authz.DELETE("/news_images/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&NewsImage{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
