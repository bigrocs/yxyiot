package yxyiot

import (
	"fmt"
	"os"
	"testing"

	"github.com/bigrocs/yxyiot/requests"
)

func TestScan(t *testing.T) {
	// 创建连接
	client := NewClient()
	client.Config.AppId = os.Getenv("yxyiot_AppId")
	client.Config.AppSecret = os.Getenv("yxyiot_AppSecret")
	client.Config.Sandbox = false
	// 配置参数
	request := requests.NewCommonRequest()
	request.ApiName = "play"
	request.BizContent = map[string]interface{}{
		"devName":       os.Getenv("yxyiot_devName"),
		"bizType":       "1",
		"money":         "2.5",
		"broadCastType": "1",
	}
	// 请求
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Println(err)
	}
	r, err := response.GetVerifySignDataMap()
	fmt.Println("TestPlay", r, err)
	t.Log(r, err, "|||")
}

// 指令说明
// https://docs.qq.com/sheet/DQkNoTm9uVWFyeEdU?tab=BB08J2
// func TestPrint(t *testing.T) {
// 	// 创建连接
// 	client := NewClient()
// 	client.Config.ClientId = os.Getenv("yxyiot_ClientId")
// 	client.Config.ClientSecret = os.Getenv("yxyiot_ClientSecret")
// 	client.Config.Sandbox = false
// 	// 配置参数
// 	request := requests.NewCommonRequest()
// 	request.ApiName = "print.index"
// 	request.BizContent = map[string]interface{}{
// 		"access_token": os.Getenv("yxyiot_access_token"),
// 		"machine_code": os.Getenv("yxyiot_machine_code"),
// 		"content": `
// 		@@测试
// 		<FS3>测试标题</FS3>
// 		<FH2>asdas</FH2>
// 		`,
// 		"origin_id": uuid.NewV4().String() + "A",
// 	}
// 	// 请求
// 	response, err := client.ProcessCommonRequest(request)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	r, err := response.GetVerifySignDataMap()
// 	fmt.Println("TestPlay", r, err)
// 	t.Log(r, err, "|||")
// }
