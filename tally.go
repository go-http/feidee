package feidee

import (
	"fmt"
	"net/url"
	"strconv"
)

// 流水信息
type Tally struct {
	Account        int
	BuyerAcount    string //支出、收入时为交易账户，转账时为转出账户
	BuyerAcountId  int    //支出、收入时为交易账户，转账时为转出账户
	CategoryId     int
	CategoryName   string
	MemberId       int
	MemberName     string
	TranType       int
	TranName       string
	ProjectId      int
	ProjectName    string
	StoreId        int    //商户ID
	SellerAcount   string //转账时表示转入账户
	SellerAcountId int    //转账时表示转入账户

	ItemAmount     float32 //交易金额，转账时该数值和CurrencyAmount相同
	CurrencyAmount float32 //转账金额

	Relation     string
	CategoryIcon string
	Url          string
	Content      string
	ImgId        int
	TranId       int
	Memo         string

	Date DateInfo
}

// 生成用于更新的url.Values参数
func (t Tally) ToUpdateParams() url.Values {
	data := url.Values{}

	data.Set("id", strconv.Itoa(t.TranId))

	if t.TranType == TranTypeTransfer {
		data.Set("in_account", strconv.Itoa(t.SellerAcountId))
		data.Set("out_account", strconv.Itoa(t.BuyerAcountId))
	} else {
		data.Set("account", strconv.Itoa(t.Account))
	}

	data.Set("store", strconv.Itoa(t.StoreId))

	data.Set("category", strconv.Itoa(t.CategoryId))
	data.Set("project", strconv.Itoa(t.ProjectId))
	data.Set("member", strconv.Itoa(t.MemberId))

	data.Set("memo", t.Memo)
	data.Set("url", t.Url)

	strTime := fmt.Sprintf("%d.%d.%d %d:%d:%d", 1900+t.Date.Year, t.Date.Month+1, t.Date.Date, t.Date.Hours, t.Date.Minutes, t.Date.Seconds)
	data.Set("time", strTime)

	data.Set("price", fmt.Sprintf("%.2f", t.ItemAmount))
	data.Set("price2", fmt.Sprintf("%.2f", t.CurrencyAmount))

	return data
}
