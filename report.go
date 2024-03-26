package feidee

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// 对账报表
type CompareInfo struct {
	Balance    float32 //当前余额
	DayBalance float32 //当日余额（当日收入-当日支出）
	Money      struct {
		Claims float32
		Income float32
		In     float32 //流入
		Debet  float32
		Out    float32 //留出
		Payout float32
	}
	Date DateInfo
}

// 报表页响应
type CompareReportResponse struct {
	PageInfo
	List []CompareInfo
}

// 获取所有对账报表
func (cli *Client) CompareReport(accountId int, begin, end time.Time) ([]CompareInfo, error) {
	pageCount := 1
	list := []CompareInfo{}
	for page := 1; page <= pageCount; page += 1 {
		info, err := cli.CompareReportByPage(accountId, begin, end, page)
		if err != nil {
			return nil, err
		}
		pageCount = info.PageCount
		list = append(list, info.List...)
	}

	return list, nil
}

// 获取单页对账报表
func (cli *Client) CompareReportByPage(accountId int, begin, end time.Time, page int) (CompareReportResponse, error) {
	data := url.Values{}
	data.Set("m", "compare")
	data.Set("page", strconv.Itoa(page))
	data.Set("cardId", strconv.Itoa(accountId))
	data.Set("endDate", end.Format("2006.01.02"))
	data.Set("beginDate", begin.Format("2006.01.02"))

	resp, err := cli.PostForm(BaseUrl+"/report.rmi", data)
	if err != nil {
		return CompareReportResponse{}, fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	var respInfo CompareReportResponse
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return CompareReportResponse{}, fmt.Errorf("读取出错: %s", err)
	}

	return respInfo, nil
}

// 日常收支报表
type DailyReportList []struct {
	IdName
	Total float32
	List  []struct {
		IdName
		Amount float32
	} `json:"c"`
}

type DailyReport struct {
	InAmount  float32
	OutAmount float32

	Symbol string
	MaxIn  float32 `json:"maxI"`
	MaxOut float32 `json:"maxO"`

	InList  DailyReportList `json:"inlst"`
	OutList DailyReportList `json:"outlst"`
}

// 日常收支报表
func (cli *Client) DailyReport(begin, end time.Time, params url.Values) (DailyReport, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("m", "daily")
	params.Set("endDate", end.Format("2006.01.02"))
	params.Set("beginDate", begin.Format("2006.01.02"))

	resp, err := cli.PostForm(BaseUrl+"/report.rmi", params)
	if err != nil {
		return DailyReport{}, fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	var respInfo DailyReport
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return DailyReport{}, fmt.Errorf("读取出错: %s", err)
	}

	return respInfo, nil
}

type AssetReportList []struct {
	Name   string
	Symbol string
	Amount float32
}

type AssetReport struct {
	InAbsAmount  float32
	InAmount     float32
	InList       AssetReportList `json:"inlst"`  //原始币种金额
	InListR      AssetReportList `json:"inlstr"` //外币自动折算成人民币
	OutAbsAmount float32
	OutAmount    float32
	OutList      AssetReportList `json:"outlst"`  //原始币种金额
	OutListR     AssetReportList `json:"outlstr"` //外币自动折算成人民币
	Symbol       string
}

// 资产负债表
func (cli *Client) AssetReport() (AssetReport, error) {
	params := url.Values{"m": {"asset"}}

	resp, err := cli.PostForm(BaseUrl+"/report.rmi", params)
	if err != nil {
		return AssetReport{}, fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	var respInfo AssetReport
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return AssetReport{}, fmt.Errorf("读取出错: %s", err)
	}

	return respInfo, nil
}
