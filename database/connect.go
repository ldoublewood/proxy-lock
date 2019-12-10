package database

import (
	"gin-gorm-example/config"
	"gin-gorm-example/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)


var  DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	conf := config.Get()
	db, err := gorm.Open(conf.Dialect, conf.DSN)

	if err == nil {
		db.LogMode(true)
		db.DB().SetMaxIdleConns(conf.MaxIdleConn)
		DB = db
		db.AutoMigrate(&models.Proxy{})
		return db, err
	}
	return nil, err
}