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
		"devName":       "bsj00575",
		"bizType":       "2",
		"content":       "张三收款成功3467.91元",
		"money":         "24222.5",
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
func TestPrintPlay(t *testing.T) { // 打印机播报模式
	// 创建连接
	// client := NewClient()
	// client.Config.AppId = os.Getenv("yxyiot_AppId")
	// client.Config.AppSecret = os.Getenv("yxyiot_AppSecret")
	// client.Config.Sandbox = false
	// // 配置参数
	// request := requests.NewCommonRequest()
	// request.ApiName = "print"
	// request.BizContent = map[string]interface{}{
	// 	"devName":   "bsj00576",
	// 	"actWay":    "2",
	// 	"voiceJson": `{"devName":"bsj00575","bizType":"1","money":"4.5","broadCastType":"1"}`,
	// }
	// // 请求
	// response, err := client.ProcessCommonRequest(request)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// r, err := response.GetVerifySignDataMap()
	// fmt.Println("TestPlay", r, err)
	// t.Log(r, err, "|||")
}

func TestPrint(t *testing.T) { // 打印机播报模式
	// 创建连接
	// client := NewClient()
	// client.Config.AppId = os.Getenv("yxyiot_AppId")
	// client.Config.AppSecret = os.Getenv("yxyiot_AppSecret")
	// client.Config.Sandbox = false
	// // 配置参数
	// request := requests.NewCommonRequest()
	// request.ApiName = "print"
	// request.BizContent = map[string]interface{}{
	// 	"devName": "bsj00576",
	// 	"actWay":  "1",
	// 	"data": `格式模板：
	// 	<RS:2><L><CB>####### 6 ######</CB></L><BR>
	// 	<RS:2><C>*沙县小吃*</C><BR>
	// 	<RS:2><CB>--已在线支付--</CB><BR>
	// 	<RS:2><C>----------------</C><BR>
	// 	送达时间: 2021-01-04 13:28:50<BR>
	// 	送达时间: 2021-01-04 12:28:50<BR>
	// 	订单编号: 1200897812792015996<BR>
	// 	-------01号篮子------<BR>
	// 	<L>爆炒肥肠      x2   80.0</L><BR>
	// 	<L>蚂蚁上树      x1   12.3</L><BR>
	// 	[会员减配送费: 0.0]<BR>
	// 	<RS:1>[商家承担的配送费: 1.0]<BR>
	// 	配送费: ￥1.0<BR>
	// 	<B>实付: ￥43.3</B><BR>
	// 	<B>手机号: 13012345678</B><BR>
	// 	<RS:2><B>仲恺高新区惠风西3路1号</B><BR>
	// 	备  注: 不要辣椒<BR>
	// 	发票抬头: 惠州市博实结科技有限公司<BR>
	// 	<LOGO>
	// 	<CUT>`,
	// }
	// // 请求
	// response, err := client.ProcessCommonRequest(request)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// r, err := response.GetVerifySignDataMap()
	// fmt.Println("TestPlay", r, err)
	// t.Log(r, err, "|||")
}
