package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerNewsRoutes(r *gin.Engine, db *gorm.DB) {
	// 读
	r.GET("/api/v1/news", func(c *gin.Context) {
		var list []News
		db.Order("pinned DESC").Order("published_at DESC").Order("id DESC").Find(&list)
		c.JSON(http.StatusOK, list)
	})

	// 写
	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	authz.POST("/news", func(c *gin.Context) {
		var req newsCreateReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ts, err := parsePublishedAt(req.PublishedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		n := News{
			Title:       req.Title,
			Summary:     req.Summary,
			Link:        req.Link,
			Pinned:      req.Pinned,
			PublishedAt: ts,
		}
		if err := db.Create(&n).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, n)
	})

	authz.DELETE("/news/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&News{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
