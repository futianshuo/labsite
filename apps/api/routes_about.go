package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerAboutRoutes(r *gin.Engine, db *gorm.DB) {
	// 读：正文 + 图片
	r.GET("/api/v1/about", func(c *gin.Context) {
		var ap AboutPage
		if err := db.First(&ap).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var imgs []AboutImage
		db.Order("sort ASC").Order("created_at DESC").Find(&imgs)
		c.JSON(http.StatusOK, gin.H{"body": ap.Body, "images": imgs})
	})

	// 写正文
	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	authz.POST("/about", func(c *gin.Context) {
		var req struct {
			Body string `json:"body"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ap AboutPage
		if err := db.First(&ap).Error; err != nil {
			ap.Body = req.Body
			if err2 := db.Create(&ap).Error; err2 != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
				return
			}
		} else {
			ap.Body = req.Body
			if err := db.Save(&ap).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// 图片列表（公开）
	r.GET("/api/v1/about_images", func(c *gin.Context) {
		var imgs []AboutImage
		db.Order("sort ASC").Order("created_at DESC").Find(&imgs)
		c.JSON(http.StatusOK, imgs)
	})

	// 新增图片
	authz.POST("/about_images", func(c *gin.Context) {
		var req AboutImage
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.ImageURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image_url required"})
			return
		}
		if err := db.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})

	// 更新图片标题/排序
	authz.PATCH("/about_images/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Title *string `json:"title"`
			Sort  *int    `json:"sort"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var img AboutImage
		if err := db.First(&img, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}
		if req.Title != nil {
			img.Title = *req.Title
		}
		if req.Sort != nil {
			img.Sort = *req.Sort
		}
		if err := db.Save(&img).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, img)
	})

	// 删除图片
	authz.DELETE("/about_images/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&AboutImage{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
