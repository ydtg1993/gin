package core

import (
	_ "database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var Mysql *gorm.DB

func init() {
	db1Dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local&timeout=%s",
		Config.GetString("mysql.username"),
		Config.GetString("mysql.password"),
		Config.GetString("mysql.host"),
		Config.GetString("mysql.port"),
		Config.GetString("mysql.database"),
		Config.GetString("mysql.charset"),
		Config.GetString("mysql.timeout"),
	)
	var err error
	Mysql, err = gorm.Open(mysql.Open(db1Dsn), &gorm.Config{
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(fmt.Sprintf("unable connect to database : %s\n", err))
	}

	sqlDB, err := Mysql.DB()
	if err != nil {
		panic(fmt.Sprintf("unable to get underlying DB: %s\n", err))
	}
	// Set maximum open connections
	sqlDB.SetMaxOpenConns(Config.GetInt("mysql.max_open_num"))
	// Set maximum idle connections
	sqlDB.SetMaxIdleConns(Config.GetInt("mysql.max_idle_num"))
	// Set connection's maximum reusable time
	sqlDB.SetConnMaxLifetime(Config.GetDuration("mysql.connect_lifetime") * time.Minute)
}
