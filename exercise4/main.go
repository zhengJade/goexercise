// 39.105.57.72
package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type PNum struct {
	id    int
	pnums string
}

// 定义一个全局对象db
type DB struct {
	db *sql.DB
}

var MDb DB

// 定义一个初始化数据库的函数

func initTabel() {
	s := "1234567890&123 456 7891&(123) 456 7892&(123) 456-7893&123-456-7894&123-456-7890&1234567892&(123)456-7892"
	data := strings.Split(s, "&")
	for _, num := range data {
		err := MDb.add(num)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func initDB() (err error) {
	// DSN:Data Source Name
	dsn := "root:jz1997...@tcp(39.105.57.72:3306)/phonenums?charset=utf8mb4&parseTime=True"
	// 不会校验账号密码是否正确
	// 注意！！！这里不要使用:=，我们是给全局变量赋值，然后在main函数中使用全局变量db
	MDb.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	MDb.db.SetConnMaxLifetime(time.Minute * 3)
	MDb.db.SetMaxOpenConns(10)
	MDb.db.SetMaxIdleConns(10)
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = MDb.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	initTabel()
	nums, err := MDb.DbChooseAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	filter := make(map[string]bool)
	for _, num := range nums {
		ret := Normalizer(num.pnums)
		if _, ok := filter[ret]; ok {
			MDb.remove(num.id)
		} else {
			filter[ret] = true
			MDb.update(num.id, ret)
		}
	}
	last, err := MDb.DbChooseAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(last)
}

func (mDb DB) DbChooseAll() ([]PNum, error) {
	rows, err := mDb.db.Query("select * from nums")
	ret := make([]PNum, 0)
	if err != nil {
		return ret, err
	}

	for rows.Next() {
		var p PNum
		rows.Scan(&p.id, &p.pnums)
		ret = append(ret, p)
	}
	return ret, nil
}

func Normalizer(num string) string {
	ret := ""
	for _, c := range num {
		if c >= '0' && c <= '9' {
			ret += string(c)
		}
	}
	return ret
}

func (mDb DB) update(id int, value string) error {
	_, err := mDb.db.Exec("update nums set pnums=? where id=?", value, id)
	if err != nil {
		return err
	}
	return nil
}

func (mDb DB) remove(id int) error {
	_, err := mDb.db.Exec("delete from nums where id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func (mDb DB) add(value string) error {
	_, err := mDb.db.Exec("insert into nums (pnums) value (?)", value)
	if err != nil {
		return err
	}
	return nil
}
