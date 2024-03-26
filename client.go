package feidee

import (
	"net/http"
	"net/http/cookiejar"
	"os"
)

// 请求的基础链接地址
const (
	BaseUrl  = "https://www.sui.com"
	LoginUrl = "https://login.sui.com"
)

// 交易类型
const (
	TranTypePayout   = 1 //支出
	TranTypeTransfer = 2 //转账
	TranTypeIncome   = 5 //收入
)

// 多个响应使用的分页结构
type PageInfo struct {
	PageCount int
	PageNo    int
}

// 多个响应使用的日期结构
type DateInfo struct {
	Year           int //从1900算起第N年
	Month          int //从0开始的月份
	Date           int //日期
	Day            int //星期
	Hours          int //时北京时间
	Minutes        int //分
	Seconds        int //秒
	Time           int //Unix时间戳* 1000
	TimezoneOffset int //与UTC时间的相差的小时数
}

// 包含ID、Name两个属性的结构
type IdName struct {
	Id   int
	Name string
}

// 包含收入、支出两个属性的结构
type IncomeAndPayout struct {
	Income float32
	Payout float32
}

// 执行操作的Feidee客户端
type Client struct {
	httpClient *http.Client
	Verbose    bool
	AccountBook
	AccountBookList []IdName
}

// 创建客户端
func New(email, password string) (*Client, error) {
	cookieJar, _ := cookiejar.New(nil)
	cli := &Client{httpClient: &http.Client{Jar: cookieJar}}
	err := cli.login(email, password)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

// 从环境变量中创建一个新客户端
// 用户名、密码对应的环境变量分别为FEIDEE_USERNAME、FEIDEE_PASSWORD
func NewFromEnv() (*Client, error) {
	username := os.Getenv("FEIDEE_USERNAME")
	password := os.Getenv("FEIDEE_PASSWORD")
	return New(username, password)
}
