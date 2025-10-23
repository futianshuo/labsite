package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerMemberRoutes(r *gin.Engine, db *gorm.DB) {
	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	// 写
	authz.POST("/members", func(c *gin.Context) {
		var m Member
		if err := c.BindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&m).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, m)
	})

	authz.DELETE("/members/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := db.Delete(&Member{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// 读
	r.GET("/api/v1/members", func(c *gin.Context) {
		var list []Member
		db.Order("id ASC").Find(&list)
		c.JSON(http.StatusOK, list)
	})
}
