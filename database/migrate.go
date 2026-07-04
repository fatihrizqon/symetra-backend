package database

import (
	"github.com/fatihrizqon/gofiber-microservice/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		// ── Auth & Users ──────────────────────────────────────────────────────
		&entity.User{},
		&entity.Session{},
		&entity.Credential{},
		&entity.Role{},
		&entity.Permission{},
		&entity.File{},
		&entity.RedisJob{},
	)

	seedDefaultData(db)
}

func seedDefaultData(db *gorm.DB) {
	SeedData(db)
}
