package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func clamp(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func registerHighlightRoutes(r *gin.Engine, db *gorm.DB) {
	// 读：原样
	r.GET("/api/v1/highlights", func(c *gin.Context) {
		var list []Highlight
		db.Order("sort ASC").Order("published_at DESC").Order("id DESC").Find(&list)
		c.JSON(http.StatusOK, list)
	})

	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	// 创建
	authz.POST("/highlights", func(c *gin.Context) {
		var req struct {
			Title       string `json:"title"`
			Summary     string `json:"summary"`
			PaperLink   string `json:"paper_link"`
			ImageURL    string `json:"image_url"`
			Sort        int    `json:"sort"`
			PublishedAt string `json:"published_at"`
			Split       int    `json:"split"` // 新增
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ts, err := parsePublishedAt(req.PublishedAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h := Highlight{
			Title:       req.Title,
			Summary:     req.Summary,
			PaperLink:   req.PaperLink,
			ImageURL:    req.ImageURL,
			Sort:        req.Sort,
			PublishedAt: ts,
			Split:       clamp(req.Split, 15, 60), // 新增
		}
		if err := db.Create(&h).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, h)
	})

	// 修改（排序/比例/其它）
	authz.PATCH("/highlights/:id", func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Title       *string `json:"title"`
			Summary     *string `json:"summary"`
			PaperLink   *string `json:"paper_link"`
			ImageURL    *string `json:"image_url"`
			Sort        *int    `json:"sort"`
			PublishedAt *string `json:"published_at"`
			Split       *int    `json:"split"` // 新增
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var h Highlight
		if err := db.First(&h, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "highlight not found"})
			return
		}
		if req.Title != nil {
			h.Title = *req.Title
		}
		if req.Summary != nil {
			h.Summary = *req.Summary
		}
		if req.PaperLink != nil {
			h.PaperLink = *req.PaperLink
		}
		if req.ImageURL != nil {
			h.ImageURL = *req.ImageURL
		}
		if req.Sort != nil {
			h.Sort = *req.Sort
		}
		if req.PublishedAt != nil {
			if ts, err := parsePublishedAt(*req.PublishedAt); err == nil {
				h.PublishedAt = ts
			}
		}
		if req.Split != nil {
			h.Split = clamp(*req.Split, 15, 60)
		} // 新增

		if err := db.Save(&h).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, h)
	})

	// 删除：不变
	authz.DELETE("/highlights/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Highlight{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
