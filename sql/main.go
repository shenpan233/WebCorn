package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"regexp"
	"sync"
)

var SqliteDb *sql.DB
var rWLocker *sync.RWMutex

func init() {
	rWLocker = new(sync.RWMutex)
	var err error
	SqliteDb, err = sql.Open("sqlite3", "./db.db")
	if err != nil {
		fmt.Println(err)
		return
	} else {
		rWLocker.Lock()
		_, err = SqliteDb.Exec(
			`CREATE TABLE "urls" (
	"ID"	INTEGER NOT NULL UNIQUE,
	"Url"	INTEGER NOT NULL UNIQUE,
	"time" 	int,
	"OpenStatus" boolean,
	"Status"	varchar(255),
	PRIMARY KEY("ID" AUTOINCREMENT)
)`)
		rWLocker.Unlock()
		//不管是否成功都创建数据库
		if err != nil {
			fmt.Println(err)
		}
	}

}

type Url struct {
	ID, Time    int
	Url, Status string
	OpenStatus  bool
}

func GetUrl() (res []*Url) {
	//rWLocker.RLock()
	row, err := SqliteDb.Query("SELECT * FROM urls")
	//rWLocker.RUnlock()
	if err == nil {
		defer row.Close()
		for row.Next() {
			u := new(Url)
			_ = row.Scan(&u.ID, &u.Url, &u.Time, &u.OpenStatus, &u.Status)

			res = append(res, u)
		}
	} else {
		return nil
	}
	return
}

func GetAnUrl(url string) (u *Url, b bool) {
	//rWLocker.RLock()
	row, err := SqliteDb.Query("SELECT * FROM urls WHERE url ='" + url + "'")
	defer row.Close()
	//rWLocker.RUnlock()
	if err == nil {
		row.Next()
		u = new(Url)
		_ = row.Scan(&u.ID, &u.Url, &u.Time, &u.OpenStatus, &u.Status)
		b = true
	} else {
		return nil, false
	}

	return
}

func QueryUrl(url string) bool {
	//rWLocker.RLock()
	row, err := SqliteDb.Query("SELECT * FROM urls WHERE url ='" + url + "'")
	defer row.Close()
	if err == nil {
		return row.Next()
	}
	//rWLocker.Unlock()
	return false
}

func AddUrl(url Url) string {
	match, _ := regexp.MatchString(`[a-zA-z]+://[^\s]*`, url.Url)
	if !match {
		return "格式错误，正确应为:http(s)://blog.shenpan233.cn/api[这是个例子],需带有http://"
	}
	//rWLocker.Lock()
	_, err := SqliteDb.Exec(fmt.Sprintf(`INSERT INTO urls (Url,time,Status,OpenStatus) VALUES ('%s',%d,'never worked',%v)`, url.Url, url.Time, url.OpenStatus))
	//rWLocker.Unlock()
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: urls.url" {
			return "该域名已存在!无法重复添加"
		}

		return err.Error()
	}

	return "Success"
}

func UpdateUrl(url *Url) (b bool) {
	//理论上url都是对的所以不用做查询
	//rWLocker.Lock()
	_, err := SqliteDb.Exec(fmt.Sprintf("UPDATE urls SET status = '%s',time=%d,openStatus=%v WHERE url = '%s'", url.Status, url.Time, url.OpenStatus, url.Url))
	//rWLocker.Unlock()
	if err != nil {

		return false
	}
	return true
}

func DeleteId(id int) (string, bool) {
	//rWLocker.Lock()
	fmt.Println(fmt.Sprintf("DELETE FROM urls WHERE id=%d", id))
	_, err := SqliteDb.Exec(fmt.Sprintf("DELETE FROM urls WHERE id=%d", id))
	//rWLocker.Unlock()

	if err != nil {
		return err.Error(), false
	} else {
		return "Delete Success", true
	}
}
