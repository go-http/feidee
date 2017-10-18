package feidee

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"sort"
	"strconv"
	"time"
)

//流水信息
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
	SellerAcount   string //转账时表示转入账户
	SellerAcountId int    //转账时表示转入账户

	ItemAmount     float32 //交易金额，转账时该数值和CurrencyAmount相同
	CurrencyAmount float32 //转账金额

	Relation        string
	CategoryIcon    string
	Url             string
	Content         string
	ImgId           int
	TranId          int
	TransferStoreId int
	Memo            string

	Date DateInfo
}

//生成用于更新的url.Values参数
func (t Tally) ToUpdateParams() url.Values {
	data := url.Values{}

	data.Set("id", strconv.Itoa(t.TranId))

	if t.TranType == TranTypeTransfer {
		data.Set("in_account", strconv.Itoa(t.SellerAcountId))
		data.Set("out_account", strconv.Itoa(t.BuyerAcountId))
	} else {
		data.Set("account", strconv.Itoa(t.Account))
	}

	data.Set("store", strconv.Itoa(t.TransferStoreId))
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

//流水组（通常是按天分组）
type TallyGroup struct {
	IncomeAndPayout
	List []Tally
}

//查询流水接口的响应
type TallyResponseInfo struct {
	IncomeAndPayout
	PageNo    int
	PageCount int
	Symbol    string
	EndDate   string
	BeginDate string
	Groups    []TallyGroup
}

//获取流水，可用参数包括
//    bids账户、cids科目、mids类型、pids项目、sids商户、memids成员，这几个参数都是逗号分割的ID列表
//    order:  排序字段，支持: project_id项目排序、buyer_name账户、item_amount金额、tran_type类型、category_id科目、tran_time时间
//    isDesc: 是否降序，0升序、1降序
//    note:   搜备注关键字
func (cli *Client) TallyList(begin, end time.Time, data url.Values) (TallyResponseInfo, error) {
	//先取出所有页的信息构成一个Slices
	pageCount := 1
	infos := []TallyResponseInfo{}
	for page := 1; page <= pageCount; page += 1 {
		pageInfo, err := cli.TallyListByPage(begin, end, data, page)
		if err != nil {
			return TallyResponseInfo{}, err
		}

		pageCount = pageInfo.PageCount
		infos = append(infos, pageInfo)
	}

	//如果没有账单记录，直接返回
	if len(infos) == 0 {
		return TallyResponseInfo{}, nil
	}

	//由于系统分页是按数量进行，有可能导致同一天分散到两个Group
	//所以把抛弃原来的Group，把Group里的Tally信息重新按照日期填入map中
	//map的key是由年月日组成的整数，例如20170707，这样效率比字符串更高
	tallyMap := map[int][]Tally{}
	for _, info := range infos {
		for _, group := range info.Groups {
			for _, tally := range group.List {
				key := (tally.Date.Year + 1900) * 10000
				key += (tally.Date.Month + 1) * 100
				key += tally.Date.Date

				tallyMap[key] = append([]Tally{tally}, tallyMap[key]...)
			}
		}
	}

	//遍历Tally构成的map，重组Group和ResponseInfo
	dateMax := 0
	dateMin := 99999999
	mergedInfo := TallyResponseInfo{Groups: []TallyGroup{}}
	for t, tallies := range tallyMap {
		if t > dateMax {
			dateMax = t
		}
		if t < dateMin {
			dateMin = t
		}
		group := TallyGroup{List: []Tally{}}
		for _, tally := range tallies {
			if tally.TranType == TranTypePayout {
				group.Payout += tally.ItemAmount
			} else if tally.TranType == TranTypeIncome {
				group.Income += tally.ItemAmount
			}
			group.List = append(group.List, tally)
		}
		mergedInfo.Income += group.Income
		mergedInfo.Payout += group.Payout

		//FIXME:
		//    由于Group重组，所以记录的顺序会变化
		//    这里暂时按照时间排序, 后续改为根据传入参数确定排序方式
		sort.Slice(group.List, func(i, j int) bool { return group.List[i].Date.Time > group.List[j].Date.Time })
		mergedInfo.Groups = append(mergedInfo.Groups, group)
	}

	mergedInfo.BeginDate = fmt.Sprintf("%4d.%02d.%02d", dateMin/10000, (dateMin/100)%100, dateMin%100)
	mergedInfo.EndDate = fmt.Sprintf("%4d.%02d.%02d", dateMax/10000, (dateMax/100)%100, dateMax%100)

	return mergedInfo, nil
}

func (cli *Client) TallyListByPage(begin, end time.Time, data url.Values, page int) (TallyResponseInfo, error) {
	if data == nil {
		data = url.Values{}
	}
	data.Set("opt", "list2")
	data.Set("page", strconv.Itoa(page))

	if !begin.IsZero() {
		data.Set("beginDate", begin.Format("2006.01.02"))
	}

	if !end.IsZero() {
		data.Set("endDate", end.Format("2006.01.02"))
	}

	//部分参数必须要有默认值
	if data.Get("bids") == "" {
		data.Set("bids", "0")
	}
	if data.Get("cids") == "" {
		data.Set("cids", "0")
	}
	if data.Get("mids") == "" {
		data.Set("mids", "0")
	}
	if data.Get("pids") == "" {
		data.Set("pids", "0")
	}
	if data.Get("sids") == "" {
		data.Set("sids", "0")
	}
	if data.Get("memids") == "" {
		data.Set("memids", "0")
	}

	resp, err := cli.httpClient.PostForm(BaseUrl+"/money/tally/new.rmi", data)
	if err != nil {
		return TallyResponseInfo{}, fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	var respInfo TallyResponseInfo
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return respInfo, fmt.Errorf("解析出错: %s", err)
	}

	return respInfo, nil
}

//获取按月汇总的收支情况，key为201707格式
func (cli *Client) MonthIncomeAndPayoutMap(beginYear, endYear int) (map[int]IncomeAndPayout, error) {
	data := url.Values{}
	data.Set("opt", "someYearSum")
	data.Set("endYear", strconv.Itoa(endYear))
	data.Set("beginYear", strconv.Itoa(beginYear))

	resp, err := cli.httpClient.PostForm(BaseUrl+"/money/tally/new.rmi", data)
	if err != nil {
		return nil, fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	respInfo := map[string]map[string]IncomeAndPayout{}
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return nil, fmt.Errorf("解析出错: %s", err)
	}

	infoMap := map[int]IncomeAndPayout{}
	for yearKey, yearInfo := range respInfo {
		for monthKey, monthInfo := range yearInfo {
			year, _ := strconv.Atoi(yearKey)
			month, _ := strconv.Atoi(monthKey)

			key := year*100 + month
			infoMap[key] = monthInfo
		}
	}

	return infoMap, nil
}

//更新交易的接口
func (cli *Client) TallyUpdate(tally Tally, updateData url.Values) error {
	data := tally.ToUpdateParams()
	for k, vv := range updateData {
		data.Del(k)
		for _, v := range vv {
			data.Add(k, v)
		}
	}

	var tranType string
	switch tally.TranType {
	case TranTypePayout:
		tranType = "payout"
	case TranTypeTransfer:
		tranType = "transfer"
	case TranTypeIncome:
		tranType = "income"
	default:
		return fmt.Errorf("未知的交易类型%d", tally.TranType)
	}

	resp, err := cli.httpClient.PostForm(BaseUrl+"/money/tally/"+tranType+".rmi", data)
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	if string(b) == "{result:'ok'}" {
		return nil
	}

	return fmt.Errorf("请求出错: %s", string(b))
}

//添加交易的接口
//TODO:增加对其他交易类型的支持以及错误处理的支持
func (cli *Client) PayoutCreate(tally Tally, when time.Time) error {
	data := url.Values{}

	//添加时始终为0
	data.Set("id", "0")
	data.Set("store", "0") //TODO:用正确的数据填充

	data.Set("memo", tally.Memo)

	data.Set("category", strconv.Itoa(tally.CategoryId))
	data.Set("project", strconv.Itoa(tally.ProjectId))
	data.Set("member", strconv.Itoa(tally.MemberId))
	data.Set("account", strconv.Itoa(tally.Account))

	data.Set("time", when.Format("2006-01-02 15:04"))
	data.Set("price", strconv.FormatFloat(float64(tally.ItemAmount), 'f', -1, 32))

	//TODO: 下面的字段不清楚
	//data.Set("url", "")
	//data.Set("price2", "")
	//data.Set("in_account", "")
	//data.Set("out_account", "")
	//data.Set("debt_account", "")

	resp, err := cli.httpClient.PostForm(BaseUrl+"/money/tally/payout.rmi", data)
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	return fmt.Errorf(string(b))
}
