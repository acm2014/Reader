package tools

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/axgle/mahonia"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
)

func Convert(src string, srcCode string, targetCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(targetCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}

func ChineseGBKEncode(name string) string {
	enc := mahonia.NewEncoder("GBK")
	s := enc.ConvertString(name)
	res := ""
	for i := 0; i < len(s); i++ {
		res = res + fmt.Sprintf("%%%X", s[i])
	}
	return res
}

const mysqlUserName = "root"
const mysqlPassWord = "********"
const mysqlHost = "39.105.20.61"
const mysqlDb = "Reader"
const mysqlPort = 3306

var once sync.Once
var db *sql.DB

func NewMysqlConnection() (*sql.DB, error) {
	once.Do(func() {
		db, _ = manager.New(mysqlDb, mysqlUserName, mysqlPassWord, mysqlHost).Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
			manager.SetTimeout(1*time.Second),
			manager.SetReadTimeout(1*time.Second),
		).Port(mysqlPort).Open(true)
		db.SetMaxIdleConns(2000)
		db.SetMaxOpenConns(2000)
	})
	return db, nil
}
