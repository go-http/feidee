// 随手记导出工具
package main

import (
	"github.com/go-http/feidee/v2"

	"flag"
	"log"
	"time"
)

func main() {
	var username, password, accountBook, filename string

	flag.StringVar(&username, "u", "", "用户名/邮箱地址")
	flag.StringVar(&password, "p", "", "密码")
	flag.StringVar(&accountBook, "b", "默认账本", "账本名")
	flag.StringVar(&filename, "f", time.Now().Format("2006-01-02.xls"), "保存的文件名")
	flag.Parse()

	if username == "" || password == "" {
		flag.Usage()
		return
	}

	client, err := feidee.New(username, password)
	if err != nil {
		log.Fatalf("登录失败:%s", err)
	}

	err = client.SyncAccountBookList()
	if err != nil {
		log.Fatalf("刷新账本列表失败: %s", err)
	}

	err = client.SwitchAccountBook(accountBook)
	if err != nil {
		log.Fatalf("切换账本失败: %s", err)
	}

	err = client.ExportToFile(filename)
	if err != nil {
		log.Fatalf("下载失败%s", err)
	}
}
