package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	User string
	PWD  string
	Addr string
	Name string
}

type DB struct {
	*gorm.DB
}

var db *DB = nil

func GetInstance() *DB {
	if db == nil {
		db = &DB{}
	}

	return db
}

func (db *DB) Init(c *Config) (err error) {

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.PWD, c.Addr, c.Name)
	db.DB, err = gorm.Open(
		mysql.New(
			mysql.Config{
				DSN: connStr,
			},
		),
		&gorm.Config{},
	)

	return err
}
