package feidee

import (
	//"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"io/ioutil"
)

// 周期账手动入账
func (cli *Client) BillEntry(id int, day time.Time, money float64) (string, error) {
	data := url.Values{}
	data.Set("opt", "entry")
	data.Set("id", strconv.Itoa(id))
	data.Set("date", day.Format("2006.1.2"))
	data.Set("money", strconv.FormatFloat(money, 'f', -1, 32))

	fmt.Printf("%#v", data)

	resp, err := cli.PostForm(BaseUrl+"/bill/index.rmi", data)
	if err != nil {
		return "", fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	//返回结果是{result:'false'}或者{result:money}
	//不是合法的JSON格式
	//var respInfo struct{
	//	Id int
	//	ErrorInfo string
	//	Result string
	//}
	//err = json.NewDecoder(resp.Body).Decode(&respInfo)

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("读取出错: %s", err)
	}

	str := string(b)

	if str == "{result:'false'}" {
		return "", fmt.Errorf("失败: %s", str)
	}

	return str, nil
}
