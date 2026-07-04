package database

import (
	"log"

	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	log.Println("[SEED] Starting seed data...")

	seedAdminUser(db)

	log.Println("[SEED] Seed complete.")
}

func seedAdminUser(db *gorm.DB) {
	// 1. Seed Permissions
	var pRead, pWrite, pRbacManage entity.Permission
	if err := db.Where("name = ?", "users.read").First(&pRead).Error; err != nil {
		pRead = entity.Permission{Name: "users.read", Description: "Read user data"}
		db.Create(&pRead)
	}
	if err := db.Where("name = ?", "users.write").First(&pWrite).Error; err != nil {
		pWrite = entity.Permission{Name: "users.write", Description: "Write/modify user data"}
		db.Create(&pWrite)
	}
	if err := db.Where("name = ?", "rbac.manage").First(&pRbacManage).Error; err != nil {
		pRbacManage = entity.Permission{Name: "rbac.manage", Description: "Manage RBAC (Roles and Permissions)"}
		db.Create(&pRbacManage)
	}

	// 2. Seed Role
	var roleAdmin entity.Role
	if err := db.Where("name = ?", "admin").First(&roleAdmin).Error; err != nil {
		roleAdmin = entity.Role{
			Name:        "admin",
			Description: "Administrator Role",
			Permissions: []entity.Permission{pRead, pWrite, pRbacManage},
		}
		db.Create(&roleAdmin)
	} else {
		db.Model(&roleAdmin).Association("Permissions").Append(&pRbacManage)
	}

	var user entity.User
	if db.Where("email = ?", "admin@example.com").First(&user).Error != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&entity.User{
			Username: "admin", Name: "System Administrator",
			Email:    "admin@example.com",
			Password: string(hash), Status: 1,
			Roles: []entity.Role{roleAdmin},
		})
		log.Println("[SEED] Admin user created: admin@example.com / password")
	} else {
		db.Model(&user).Association("Roles").Append(&roleAdmin)
	}
	if db.Where("email = ?", "operator@example.com").First(&user).Error != nil {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&entity.User{
			Username: "operator", Name: "System Operator",
			Email:    "operator@example.com",
			Password: string(hash), Status: 0,
		})
		log.Println("[SEED] operator user created: operator@example.com / password")
	}
}
