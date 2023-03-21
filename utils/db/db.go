package db

import (
	"database/sql"
	"financial-spider.go/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB
var err error

// 初始化数据库连接
func init() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	db, err = sql.Open("mysql", dataSource)

	db.SetMaxOpenConns(config.DbMaxOpenConns)
	db.SetMaxIdleConns(config.DbMaxIdleConns)
	db.SetConnMaxIdleTime(config.DbMaxIdleTime * time.Minute)
	db.SetConnMaxLifetime(config.DbMaxLifeTime * time.Minute)

	if err != nil {
		log.Fatalf("数据库连接出错 : %s", err)
	}
}

// GetDb 获取数据库连接
func GetDb() *sql.DB {
	return db
}

// ExecSQL 执行SQL
func ExecSQL(sqlStr string, args ...any) []map[string]interface{} {
	rows, err := db.Query(sqlStr, args...)
	if err != nil {
		log.Fatalf("SQL执行出错 : %s", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	data := make([]map[string]interface{}, 0)

	cols, _ := rows.Columns()
	colCount := len(cols)
	values, valuePoints := make([]interface{}, colCount), make([]interface{}, colCount)

	for rows.Next() {
		for i := 0; i < colCount; i++ {
			valuePoints[i] = &values[i]
		}
		_ = rows.Scan(valuePoints...)

		temp := make(map[string]interface{})
		for i, val := range values {
			bytes, ok := val.([]byte)
			if ok {
				temp[cols[i]] = string(bytes)
			} else {
				temp[cols[i]] = val
			}
		}
		data = append(data, temp)
	}

	return data
}
