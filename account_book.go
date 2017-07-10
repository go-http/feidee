package feidee

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

//刷新账本列表
func (cli *Client) SyncAccountBookList() error {
	resp, err := cli.httpClient.Get(BaseUrl + "/money/report_index.do")
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return fmt.Errorf("读取出错: %s", err)
	}

	idNameMap := map[int]IdName{}
	lists := doc.Find("ul.s-accountbook-all li")
	for i := range lists.Nodes {
		list := lists.Eq(i)
		name, _ := list.Attr("title")
		idStr, _ := list.Attr("data-bookid")
		id, _ := strconv.Atoi(idStr)

		idNameMap[id] = IdName{Id: id, Name: name}
	}

	cli.AccountBookMap = idNameMap
	return nil
}

//切换当前操作的账本
func (cli *Client) SwitchBook(name string) error {
	var accountBookId int
	for _, info := range cli.AccountBookMap {
		if info.Name == name {
			accountBookId = info.Id
		}
	}

	if accountBookId == 0 {
		return fmt.Errorf("未找到账本")
	}

	data := url.Values{}
	data.Set("opt", "switch")
	data.Set("switchId", strconv.Itoa(accountBookId))
	data.Set("return", "xxx") //该参数必须提供但值无所谓

	resp, err := cli.httpClient.Get(BaseUrl + "/money/systemSet/book.do?" + data.Encode())
	if err != nil {
		return fmt.Errorf("请求错误: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("响应出错: %s", resp.Status)
	}

	return nil
}
