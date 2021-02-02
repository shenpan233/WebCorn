package main

import (
	"WebCorn/corn"
	"WebCorn/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

func main() {
	corn.Corn()
	r := gin.Default()
	//r.LoadHTMLGlob("E:\\go\\WebCorn\\html\\*.html")
	r.LoadHTMLGlob("./html/*.html")
	r.GET("", func(c *gin.Context) {
		urls := sql.GetUrl()
		var load string
		for _, url := range urls {
			isOpen := "false"
			if url.OpenStatus {
				isOpen = "true"
			}

			load += fmt.Sprintf(`
                    <tr>
						<td>%d</td>
                        <td>%s</td>
                        <td>%d</td>
                        <td>%s</td>
                        <td>%s</td>
                        <td>
                            <a href="edit?url=%s" class="btn btn-info btn-xs">编辑</a>&nbsp;
                            <a data-id="%d" id="del" class="btn btn-xs btn-danger" onclick="return confirm('你确实要删除此记录吗？');">删除</a>
                        </td>
                    </tr>`, url.ID, url.Url, url.Time, isOpen, url.Status, url.Url, url.ID)

		}

		c.HTML(200, "index.html", gin.H{
			"loadData": template.HTML(load),
		})
	})

	r.GET("edit", func(c *gin.Context) {
		url := c.Query("url")
		r, b := sql.GetAnUrl(url)
		if !b {
			c.String(200, "URL找不到指定的信息")
			return
		} else {
			c.HTML(200, "edit.html", gin.H{
				"url":        url,
				"time":       r.Time,
				"openStatus": r.OpenStatus,
			})
			return
		}

	})

	r.GET("add", func(c *gin.Context) {
		c.HTML(200, "edit.html", gin.H{
			"url":        "",
			"time":       "",
			"openStatus": "true",
		})
	})

	r.POST("api/add/url", func(c *gin.Context) {
		url := c.PostForm("url")
		time, _ := strconv.Atoi(c.PostForm("time"))
		var status bool
		if c.PostForm("status") == "true" {
			status = true
		}
		u := sql.Url{
			Url:        url,
			Time:       time,
			OpenStatus: status,
		}
		if sql.QueryUrl(url) {
			u2, b := sql.GetAnUrl(url)
			u.Status = u2.Status
			if b == true {
				c.JSON(200, gin.H{
					"code": sql.UpdateUrl(&u),
					"msg":  "←👈看图标知状态",
				})
				return
			}
		}

		msg := sql.AddUrl(u)
		code := false
		if msg == "Success" {
			code = true
			go corn.CornUrl(&u)
		}
		c.JSON(200, gin.H{
			"code": code,
			"msg":  msg,
		})
	})

	r.POST("api/delete", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.PostForm("id"))
		fmt.Println(id)
		s, b := sql.DeleteId(id) //我这个真的不是骂人，相信我简化变量命名的艺术
		c.JSON(200, gin.H{
			"code": b,
			"msg":  s,
		})
	})

	r.Run(":8686")
}
