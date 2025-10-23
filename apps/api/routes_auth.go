package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func jwtNew(uid uint, role string, exp time.Time) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UID: uid, Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
}

func registerAuthRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req struct{ Email, Password string }
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var u User
		if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		exp := time.Now().Add(7 * 24 * time.Hour)
		token := jwtToken(u.ID, u.Role, exp)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("token", token, int((7*24*time.Hour)/time.Second), "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email, "role": u.Role})
	})

	r.POST("/api/v1/auth/logout", func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Status(http.StatusNoContent)
	})

	r.GET("/api/v1/auth/me", func(c *gin.Context) {
		cl, err := currentClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var u User
		if err := db.First(&u, cl.UID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email, "role": u.Role})
	})
}

// 小工具：签 JWT
func jwtToken(uid uint, role string, exp time.Time) string {
	token := jwtNew(uid, role, exp)
	s, _ := token.SignedString(jwtSecret())
	return s
}
