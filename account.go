package feidee

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type AccountInfo struct {
	Id   int
	Name string

	// 账户当前余额
	Money    float64
	Currency string
}

//刷新账户余额
func (cli *Client) SyncAccountInfoList() error {
	resp, err := cli.Get(BaseUrl + "/account/account.do")
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("读取出错: %s", err)
	}

	for _, account := range cli.Accounts {
		idStr := fmt.Sprintf("#acc-money-%d", account.Id)
		selection := doc.Find(idStr)
		moneySelection := selection.Find(".child-r1-money")
		money, err := strconv.ParseFloat(strings.ReplaceAll(moneySelection.Text(), ",", ""), 64)
		if err != nil {
			return fmt.Errorf("读取账户余额出错: %s", err)
		}
		currencySelection := selection.Find(".child-r1-currency")
		cli.AccountInfoList = append(cli.AccountInfoList, AccountInfo{
			Id:       account.Id,
			Name:     account.Name,
			Money:    money,
			Currency: currencySelection.Text(),
		})
	}
	return nil
}
