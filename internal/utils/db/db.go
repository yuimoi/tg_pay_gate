package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"tg_pay_gate/internal/utils/config"
)

var DB *gorm.DB

func InitDB() {
	var db *gorm.DB
	var err error
	//db, err = gorm.Open(sqlite.Open(".env/db.db"))

	dsn := filepath.Join(config.GetRootDir(), ".env", "db.db")
	db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

}
