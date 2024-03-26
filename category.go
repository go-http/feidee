package feidee

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 账单科目
type Category struct {
	IdName
	Type   int  //科目类别：支出或收入，参见TranTypeXXX常量
	IsSub  bool //是否是子科目
	SubIds []int
}

// 初始化账本、分类、账户、商家、项目、成员等信息
func (cli *Client) SyncMetaInfo() error {
	resp, err := cli.Get(BaseUrl + "/tally/new.do")
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return fmt.Errorf("读取响应出错: %s", err)
	}

	div := doc.Find("div#filter-bar div.fb-choose")

	//解析科目信息
	var accountBook AccountBook
	accountBook.Categories = parseCategories(div)

	//解析商家、成员、账户、项目信息
	accountBook.Stores = parseIdNames(div, "store")
	accountBook.Members = parseIdNames(div, "member")
	accountBook.Accounts = parseIdNames(div, "account")
	accountBook.Projects = parseIdNames(div, "project")

	cli.AccountBook = accountBook
	return nil
}

// 解析HTML文档生成科目s
func parseCategories(doc *goquery.Selection) []Category {
	categories := []Category{}
	anchors := doc.Find("div#panel-category a")
	for i := range anchors.Nodes {
		anchor := anchors.Eq(i)

		idStr, _ := anchor.Attr("id")

		id, categoryType := categoryIdTypeSplit(idStr)

		category := Category{
			Type:   categoryType,
			IsSub:  !anchor.HasClass("ctit"),
			IdName: IdName{Id: id, Name: anchor.Text()},
		}

		//找出子类ID
		if !category.IsSub {
			subArchorClass := strings.TrimSuffix(idStr, "-a")
			subArchors := anchor.Parent().Find("a." + subArchorClass)

			subIds := []int{}
			for j := range subArchors.Nodes {
				subIdStr, _ := subArchors.Eq(j).Attr("id")
				subId, _ := categoryIdTypeSplit(subIdStr)
				subIds = append(subIds, subId)
			}

			category.SubIds = subIds
		}

		categories = append(categories, category)
	}

	return categories
}

// 解析HTML文档生成类别s
func parseIdNames(doc *goquery.Selection, zone string) []IdName {
	prefix := "c" + strings.Title(zone[:3]) + "-"

	idNames := []IdName{}
	anchors := doc.Find("div#panel-" + zone + " a")
	for i := range anchors.Nodes {
		anchor := anchors.Eq(i)

		idStr, _ := anchor.Attr("id")
		if idStr == prefix+"a" {
			continue
		}

		idStr = strings.TrimSuffix(idStr, "-a")
		idStr = strings.TrimPrefix(idStr, prefix)
		id, _ := strconv.Atoi(idStr)

		idNames = append(idNames, IdName{Id: id, Name: anchor.Text()})
	}

	return idNames
}

// 从cCat-???-??????-a格式的字符串中提取科目ID和类型
func categoryIdTypeSplit(idStr string) (int, int) {
	idStr = strings.TrimSuffix(idStr, "-a")

	var categoryType int

	if strings.HasPrefix(idStr, "cCat-out-") {
		categoryType = TranTypePayout
		idStr = strings.TrimPrefix(idStr, "cCat-out-")
	} else if strings.HasPrefix(idStr, "cCat-in-") {
		categoryType = TranTypeIncome
		idStr = strings.TrimPrefix(idStr, "cCat-in-")
	}

	id, _ := strconv.Atoi(idStr)
	return id, categoryType
}
