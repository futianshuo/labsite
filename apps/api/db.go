package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func mustDB() *gorm.DB {
	_ = godotenv.Load()
	db, err := gorm.Open(sqlite.Open("labsite.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	if err := db.AutoMigrate(
		&Member{}, &User{},
		&News{}, &Banner{}, &NewsImage{},
		&AboutPage{}, &AboutImage{},
		&Highlight{},
	); err != nil {
		log.Fatal(err)
	}

	// 种子管理员
	if email, pass := os.Getenv("ADMIN_EMAIL"), os.Getenv("ADMIN_PASSWORD"); email != "" && pass != "" {
		var n int64
		db.Model(&User{}).Where("email = ?", email).Count(&n)
		if n == 0 {
			hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
			db.Create(&User{Email: email, PasswordHash: string(hash), Role: "admin", CreatedAt: time.Now()})
			log.Println("seeded admin:", email)
		}
	}

	// AboutPage 兜底一条
	var cnt int64
	db.Model(&AboutPage{}).Count(&cnt)
	if cnt == 0 {
		db.Create(&AboutPage{Body: ""})
	}
	return db
}
