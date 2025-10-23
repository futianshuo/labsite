package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 基础数据模型：成员
type Member struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name"`
	Title         string `json:"title"`
	Role          string `json:"role"`
	AvatarURL     string `json:"avatar_url"`
	ResearchAreas string `json:"research_areas"`
	Email         string `json:"email"`
	Homepage      string `json:"homepage"`
	Destination   string `json:"destination"` // 校友去向
}

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	CreatedAt    time.Time
}

type Banner struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url"` // 轮播图地址（可 base64 或外链）
	Link      string    `json:"link"`      // 可选：点击跳转
	Sort      int       `json:"sort"`      // 排序序号（小的在前）
	CreatedAt time.Time `json:"created_at"`
}

// 独立的新闻图片库（与顶部 Banner 分开）
type NewsImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

// 新闻：改为引用新闻图库
type News struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	CoverURL    string    `json:"cover_url"` // 老数据兜底，可保留
	Link        string    `json:"link"`
	ImageID     *uint     `json:"image_id"` // ← 新字段，来自 NewsImage
	Pinned      bool      `json:"pinned"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// ====== JWT ======
type Claims struct {
	UID  uint   `json:"uid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev_secret_change_me"
	}
	return []byte(s)
}

// ====== CORS（开发环境）：允许携带 Cookie ======
func allowCORS() gin.HandlerFunc {
	// 开发固定允许本地前端
	origin := "http://localhost:5173"
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}
		c.Next()
	}
}

// 读取 Cookie 并校验 JWT
func currentClaims(c *gin.Context) (*Claims, error) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	t, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret(), nil
	})
	if err != nil || !t.Valid {
		return nil, errors.New("invalid token")
	}
	cl, ok := t.Claims.(*Claims)
	if !ok {
		return nil, errors.New("bad claims")
	}
	return cl, nil
}

// 角色守卫：仅允许 roles
func requireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		cl, err := currentClaims(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		if !allowed[cl.Role] {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Set("claims", cl)
		c.Next()
	}
}

// ====== DB 初始化与种子管理员 ======
func mustDB() *gorm.DB {
	_ = godotenv.Load()
	db, err := gorm.Open(sqlite.Open("labsite.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(&Member{}, &User{}, &News{}, &Banner{}, &NewsImage{}); err != nil {
		log.Fatal(err)
	}

	// 种子管理员（仅首次）
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminEmail != "" && adminPass != "" {
		var n int64
		db.Model(&User{}).Where("email = ?", adminEmail).Count(&n)
		if n == 0 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
			db.Create(&User{Email: adminEmail, PasswordHash: string(hash), Role: "admin"})
			log.Println("seeded admin:", adminEmail)
		}
	}
	return db
}

func main() {
	db := mustDB()
	r := gin.Default()
	r.Use(allowCORS())

	// 健康检查
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// ---- Auth ----
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req struct{ Email, Password string }
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		var u User
		if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		exp := time.Now().Add(7 * 24 * time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			UID: u.ID, Role: u.Role,
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(exp)},
		})
		s, _ := token.SignedString(jwtSecret())

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("token", s, int((7*24*time.Hour)/time.Second), "/", "", false, true) // HttpOnly
		c.JSON(200, gin.H{"id": u.ID, "email": u.Email, "role": u.Role})
	})

	r.POST("/api/v1/auth/logout", func(c *gin.Context) {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Status(204)
	})

	r.GET("/api/v1/auth/me", func(c *gin.Context) {
		cl, err := currentClaims(c)
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		var u User
		if err := db.First(&u, cl.UID).Error; err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.JSON(200, gin.H{"id": u.ID, "email": u.Email, "role": u.Role})
	})

	// ---- 写：仅 admin/editor ----
	authz := r.Group("/api/v1")
	authz.Use(requireRole("admin", "editor"))

	authz.POST("/members", func(c *gin.Context) {
		var m Member
		if err := c.BindJSON(&m); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&m).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, m)
	})

	authz.DELETE("/members/:id", func(c *gin.Context) {
		id := c.Param("id")
		var m Member
		if err := db.First(&m, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "member not found"})
			return
		}
		if err := db.Delete(&m).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	// 查询成员列表
	r.GET("/api/v1/members", func(c *gin.Context) {
		var list []Member
		db.Order("id asc").Find(&list)
		c.JSON(http.StatusOK, list)
	})

	// 公开读：新闻列表（置顶优先，其次按发布时间倒序）
	r.GET("/api/v1/news", func(c *gin.Context) {
		var list []News
		db.Order("pinned DESC").Order("published_at DESC").Order("id DESC").Find(&list)
		c.JSON(200, list)
	})

	// 受保护写：发布/删除新闻
	authz.POST("/news", func(c *gin.Context) {
		var n News
		if err := c.BindJSON(&n); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// 默认发布时间
		if n.PublishedAt.IsZero() {
			n.PublishedAt = time.Now()
		}
		if err := db.Create(&n).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, n)
	})

	authz.DELETE("/news/:id", func(c *gin.Context) {
		id := c.Param("id")
		var n News
		if err := db.First(&n, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "news not found"})
			return
		}
		if err := db.Delete(&n).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	// 公开读：轮播图（按 Sort 升序，其次按创建时间倒序）
	r.GET("/api/v1/banners", func(c *gin.Context) {
		var list []Banner
		db.Order("sort ASC").Order("created_at DESC").Find(&list)
		c.JSON(200, list)
	})

	// 受保护写：新增/删除轮播图
	authz.POST("/banners", func(c *gin.Context) {
		var b Banner
		if err := c.BindJSON(&b); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := db.Create(&b).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, b)
	})
	authz.DELETE("/banners/:id", func(c *gin.Context) {
		id := c.Param("id")
		var b Banner
		if err := db.First(&b, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "banner not found"})
			return
		}
		if err := db.Delete(&b).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	// 公开读：新闻图库
	r.GET("/api/v1/news_images", func(c *gin.Context) {
		var list []NewsImage
		db.Order("created_at DESC").Find(&list)
		c.JSON(200, list)
	})

	// 受保护写：上传/删除 新闻图库 图片
	authz.POST("/news_images", func(c *gin.Context) {
		var ni NewsImage
		if err := c.BindJSON(&ni); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if ni.ImageURL == "" {
			c.JSON(400, gin.H{"error": "image_url required"})
			return
		}
		if err := db.Create(&ni).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, ni)
	})
	authz.DELETE("/news_images/:id", func(c *gin.Context) {
		id := c.Param("id")
		var ni NewsImage
		if err := db.First(&ni, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "news image not found"})
			return
		}
		if err := db.Delete(&ni).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Status(204)
	})

	// 启动服务器
	_ = r.Run(":9000")
}
