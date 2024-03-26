package feidee

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// 刷新账本列表
func (cli *Client) SyncAccountBookList() error {
	resp, err := cli.Get(BaseUrl + "/report_index.do")
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return fmt.Errorf("读取出错: %s", err)
	}

	idNames := []IdName{}
	lists := doc.Find("ul.s-accountbook-all li")
	if len(lists.Nodes) == 0 {
		return fmt.Errorf("未找到账本列表")
	}

	for i := range lists.Nodes {
		list := lists.Eq(i)
		name, _ := list.Attr("title")
		idStr, _ := list.Attr("data-bookid")
		id, _ := strconv.Atoi(idStr)

		if id != 0 && name != "" {
			idNames = append(idNames, IdName{Id: id, Name: name})
		}
	}

	cli.AccountBookList = idNames
	return nil
}

// 切换账本
func (cli *Client) SwitchAccountBook(name string) error {
	var accountBookId int
	for _, accountBook := range cli.AccountBookList {
		if accountBook.Name == name {
			accountBookId = accountBook.Id
			break
		}
	}

	if accountBookId == 0 {
		return fmt.Errorf("未找到账本")
	}

	data := url.Values{}
	data.Set("opt", "switch")
	data.Set("switchId", strconv.Itoa(accountBookId))
	data.Set("return", "xxx") //该参数必须提供但值无所谓

	resp, err := cli.Get(BaseUrl + "/systemSet/book.do?" + data.Encode())
	if err != nil {
		return fmt.Errorf("请求错误: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("响应出错: %s", resp.Status)
	}

	err = cli.SyncMetaInfo()
	if err != nil {
		return fmt.Errorf("获取账基础信息错误: %s", err)
	}

	return nil
}
