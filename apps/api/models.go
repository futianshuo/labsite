package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ===== 模型 =====
type Member struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name"`
	Title         string `json:"title"`
	Role          string `json:"role"`
	AvatarURL     string `json:"avatar_url"`
	ResearchAreas string `json:"research_areas"`
	Email         string `json:"email"`
	Homepage      string `json:"homepage"`
	Destination   string `json:"destination"`
}

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Banner struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url"`
	Link      string    `json:"link"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

type NewsImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type News struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	CoverURL    string    `json:"cover_url"`
	Link        string    `json:"link"`
	ImageID     *uint     `json:"image_id"`
	Pinned      bool      `json:"pinned"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type AboutPage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Body      string    `json:"body"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AboutImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url"`
	Title     string    `json:"title"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
}

// ===== JWT Claims / 请求体 =====
type Claims struct {
	UID  uint   `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type newsCreateReq struct {
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Link        string `json:"link"`
	Pinned      bool   `json:"pinned"`
	PublishedAt string `json:"published_at"` // 支持 YYYY-MM-DD / RFC3339
}
