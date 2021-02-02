package corn

import (
	"WebCorn/sql"
	"fmt"
	"net/http"
	"time"
)

func Corn() {
	res := sql.GetUrl()
	fmt.Println(res)
	if res != nil {
		for _, rs := range res {
			fmt.Println(rs)
			go CornUrl(rs)
		}
	}
}

func CornUrl(url *sql.Url) {
	for {
		for url.OpenStatus {
			url, _ = sql.GetAnUrl(url.Url)
			if url == nil {
				return
			}
			if !getOnce(url) {
				return
			}
			time.Sleep(time.Duration((*url).Time) * time.Second)
		}
	}
}

func getOnce(url *sql.Url) bool {
	fmt.Println(*url)
	r, err := http.Get((*url).Url)
	if err != nil {
		(*url).Status = "GET访问错误,errorMsg=" + err.Error()
		return sql.UpdateUrl(url)
	} else {
		if r.Status == "200 OK" {
			(*url).Status = "GetSuccess,Last Get Time=" + time.Now().Format("01-02 15:04:05")
			return sql.UpdateUrl(url)

		} else {
			(*url).Status = "Get Failed,error code=" + r.Status
			return sql.UpdateUrl(url)

		}

	}
}
