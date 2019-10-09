package cache

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const userName = "root"
const passWord = "********"
const host = "39.105.20.61"
const db = "Reader"
const port = 3306

var mysqlUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", userName, passWord, host, port, db)

var con *gorm.DB

func NewMysql() (*gorm.DB, error) {
	var err error
	con, err = gorm.Open("mysql", mysqlUrl)
	if err != nil {
		return nil, err
	}
	return con, nil
}
