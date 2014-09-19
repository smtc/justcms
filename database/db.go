package database

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guotie/config"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	dbs = map[string]*gorm.DB{}
)

// 增加一个name参数仅仅是为了以后的扩展
func GetDB(dbname string) *gorm.DB {
	//return &db
	if dbname == "" {
		dbname = config.GetStringDefault("dbname", "")
	}

	db := dbs[dbname]
	if db == nil {
		newDB, err := opendb(dbname, "", "")
		if err != nil {
			panic(err)
		}
		db = &newDB
		dbs[dbname] = db
	}
	return db
}

// 根据配置文件建立数据库连接
func OpenDefaultDB() {
	OpenDB(config.GetStringDefault("dbname", ""), "", "")
}

// 供外部调用的API
func OpenDB(dbname, dbuser, dbpass string) {
	db, err := opendb(dbname, dbuser, dbpass)
	if err != nil {
		panic(err)
	}
	dbs[dbname] = &db
}

func CloseDB(dbname string) {
	GetDB(dbname).DB().Close()
}

// 建立数据库连接
func opendb(dbname, dbuser, dbpass string) (gorm.DB, error) {
	var (
		dbtype, dsn string
		db          gorm.DB
		err         error
	)

	if dbuser == "" {
		dbuser = config.GetStringDefault("dbuser", "")
	}
	if dbpass == "" {
		dbpass = config.GetStringDefault("dbpass", "")
	}

	dbtype = strings.ToLower(config.GetStringDefault("dbtype", "mysql"))
	if dbtype == "mysql" {
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbuser,
			dbpass,
			config.GetStringDefault("dbproto", "tcp"),
			config.GetStringDefault("dbhost", "127.0.0.1"),
			config.GetIntDefault("dbport", 3306),
			dbname,
		)
		//dsn += "&loc=Asia%2FShanghai"
	} else if dbtype == "pg" || dbtype == "postgres" || dbtype == "postgresql" {
		dbtype = "postgres"
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			dbuser,
			dbpass,
			config.GetStringDefault("dbhost", "127.0.0.1"),
			config.GetIntDefault("dbport", 5432),
			dbname)
	}

	//println(dbtype, dsn)
	db, err = gorm.Open(dbtype, dsn)
	if err != nil {
		log.Println(err.Error())
		return db, err
	}

	err = db.DB().Ping()
	if err != nil {
		log.Println(err.Error())
		return db, err
	}

	return db, nil
}
