package common

import (
	"fmt"
	"time"

	"github.com/bigrocs/yxyiot/config"
	"github.com/bigrocs/yxyiot/requests"
	"github.com/bigrocs/yxyiot/responses"
	"github.com/bigrocs/yxyiot/util"
	uuid "github.com/satori/go.uuid"
)

// Common 公共封装
type Common struct {
	Config   *config.Config
	Requests *requests.CommonRequest
}
type Api struct {
	Name string
	URL  string
}

var apiList = []Api{
	{
		Name: "play",
		URL:  "/openApi/dev/controlDevice.json",
	}, {
		Name: "print",
		URL:  "/print/index",
	},
}

// Action 创建新的公共连接
func (c *Common) Action(response *responses.CommonResponse) (err error) {
	return c.Request(response)
}

// APIBaseURL 默认 API 网关
func (c *Common) APIBaseURL() string { // TODO(): 后期做容灾功能
	con := c.Config
	if con.Sandbox { // 沙盒模式
		return "https://ioe.car900.com/v1"
	}
	return "https://ioe.car900.com/v1"
}

// Request 执行请求
// AppCode           string `json:"app_code"`             //API编码
// AppId             string `json:"app_id"`               //应用ID
// UniqueNo          string `json:"unique_no"`            //私钥
// PrivateKey        string `json:"private_key"`          //私钥
// yxyiotPublicKey string `json:"lin_shang_public_key"` //临商银行公钥
// MsgId             string `json:"msg_id"`               //消息通讯唯一编号，每次调用独立生成，APP级唯一
// Signature         string `json:"Signature"`            //签名值
// Timestamp         string `json:"timestamp"`            //发送请求的时间，格式"yyyy-MM-dd HH:mm:ss"
// NotifyUrl         string `json:"notify_url"`           //工商银行服务器主动通知商户服务器里指定的页面http/https路径。
// BizContent        string `json:"biz_content"`          //业务请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
// Sandbox           bool   `json:"sandbox"`              // 沙盒
func (c *Common) Request(response *responses.CommonResponse) (err error) {
	con := c.Config
	req := c.Requests
	apiUrl := ""
	for _, api := range apiList {
		if api.Name == req.ApiName {
			apiUrl = c.APIBaseURL() + api.URL
		}
	}
	// 构建配置参数
	params := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"appId":     con.AppId,
		"requestId": uuid.NewV4().String(),
		"userCode":  con.AppId,
	}
	for k, v := range req.BizContent {
		params[k] = v
	}
	param := util.FormatParam(params, con.AppSecret)
	token := util.Md5([]byte(param)) // 开发签名
	params["token"] = token
	fmt.Println(param, token)
	res, err := util.PostForm(apiUrl, param)
	if err != nil {
		return err
	}
	response.SetHttpContent(res, "string")
	return
}
