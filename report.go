package feidee

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

//对账报表
type CompareInfo struct {
	Balance    float32
	DayBalance float32
	Money      struct {
		Claims float32
		Income float32
		In     float32
		Debet  float32
		Out    float32
		Payout float32
	}
	Date DateInfo
}

//报表页响应
type CompareReportResponse struct {
	PageInfo
	List []CompareInfo
}

//获取所有对账报表
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

//获取单页对账报表
func (cli *Client) CompareReportByPage(accountId int, begin, end time.Time, page int) (CompareReportResponse, error) {
	data := url.Values{}
	data.Set("m", "compare")
	data.Set("page", strconv.Itoa(page))
	data.Set("cardId", strconv.Itoa(accountId))
	data.Set("endDate", end.Format("2016.01.02"))
	data.Set("beginDate", begin.Format("2016.01.02"))

	resp, err := cli.httpClient.PostForm(BaseUrl+"/money/report.rmi", data)
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
