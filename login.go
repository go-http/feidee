package feidee

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// 自动跳转以刷新认证信息时，最大的递归调用次数
const MaxAuthRedirectCount = 5

// 登录
func (cli *Client) login(email, password string) error {
	vccodeInfo, err := cli.getVccode()
	if err != nil {
		return fmt.Errorf("获取VCCode出错: %s", err)
	}

	err = cli.verifyUser(vccodeInfo, email, password)
	if err != nil {
		return fmt.Errorf("验证用户名密码出错: %s", err)
	}

	err = cli.authRedirect("GET", LoginUrl+"/auth.do", nil, 0)
	if err != nil {
		return fmt.Errorf("请求验证参数出错: %s", err)
	}

	err = cli.SyncAccountBookList()
	if err != nil {
		return fmt.Errorf("获取账本列表失败: %s", err)
	}
	return nil
}

type VCCodeInfo struct {
	VCCode string
	Uid    string
}

// 获取VCCode
func (cli *Client) getVccode() (VCCodeInfo, error) {
	resp, err := cli.Get(LoginUrl + "/login.do?opt=vccode")
	if err != nil {
		return VCCodeInfo{}, fmt.Errorf("请求VCCode出错: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return VCCodeInfo{}, fmt.Errorf("请求VCCode出错: %s", resp.Status)
	}

	var respInfo VCCodeInfo
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return VCCodeInfo{}, fmt.Errorf("解析VCCode响应出错: %s", err)
	}

	if respInfo.VCCode == "" {
		return VCCodeInfo{}, fmt.Errorf("未解析到合适的VCCode")
	}

	return respInfo, nil
}

// 鉴定用户
func (cli *Client) verifyUser(vccodeInfo VCCodeInfo, email, password string) error {
	//密码加密处理
	password = hexSha1(password)
	password = hexSha1(email + password)
	password = hexSha1(password + vccodeInfo.VCCode)

	data := url.Values{}
	data.Set("email", email)
	data.Set("status", "1") //是否保持登录状态: 0不保持、1保持
	data.Set("password", password)
	data.Set("uid", vccodeInfo.Uid)

	resp, err := cli.Get(LoginUrl + "/login.do?" + data.Encode())
	if err != nil {
		return fmt.Errorf("请求出错: %s", err)
	}
	defer resp.Body.Close()

	var respInfo struct{ Status string }
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return fmt.Errorf("解析出错: %s", err)
	}

	switch respInfo.Status {
	case "ok":
		return nil
	case "no":
		return fmt.Errorf("用户名密码错误")
	case "lock":
		return fmt.Errorf("密码错误次数过多十分钟后再试")
	case "lock-status":
		return fmt.Errorf("此帐号存在异常已被锁定，请上官网申述")
	default:
		return fmt.Errorf("未知状态: %s", respInfo.Status)
	}
}

// 自动跟踪认证跳转，完成验证信息刷新
func (cli *Client) authRedirect(method, address string, data url.Values, jumpCount int) error {
	if cli.Verbose {
		log.Println("第", jumpCount, "次认证跳转", method, address, "参数", data)
	}
	if jumpCount > MaxAuthRedirectCount {
		return fmt.Errorf("跳转次数太多")
	}

	var err error
	var resp *http.Response
	if method == "POST" {
		resp, err = cli.PostForm(address, data)
	} else if method == "GET" {
		resp, err = cli.Get(address + "?" + data.Encode())
	} else {
		return fmt.Errorf("未知跳转方法%s", method)
	}
	if err != nil {
		return fmt.Errorf("请求Auth出错: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return fmt.Errorf("读取响应出错: %s", err)
	}

	//解析body标签，如果不需要继续跳转则返回
	onload, _ := doc.Find("body").First().Attr("onload")
	if onload != "document.forms[0].submit()" {
		return nil
	}

	//否则解析跳转表单并递归跳转
	form := doc.Find("form").First()

	formMethod, _ := form.Attr("method")
	formAction, _ := form.Attr("action")

	formData := url.Values{}
	inputs := form.Find("input")
	for i := range inputs.Nodes {
		input := inputs.Eq(i)
		name, _ := input.Attr("name")
		value, _ := input.Attr("value")

		if name != "" {
			formData.Set(name, value)
		}
	}

	return cli.authRedirect(formMethod, formAction, formData, jumpCount+1)
}

// 密码加密,算法来自于: https://www.feidee.com/sso/js/fdLogin.js中的hex_sha1
func hexSha1(input string) string {
	sha1Sum := sha1.Sum([]byte(input))
	return hex.EncodeToString(sha1Sum[:])
}
