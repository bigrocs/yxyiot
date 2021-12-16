/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless resuired by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package responses

import (
	"encoding/json"

	"github.com/clbanning/mxj"

	"github.com/bigrocs/yxyiot/config"
	"github.com/bigrocs/yxyiot/requests"
)

// CommonResponse 公共回应
type CommonResponse struct {
	Config      *config.Config
	Request     *requests.CommonRequest
	httpContent []byte
	json        string
}

type Map *mxj.Map

// NewCommonResponse 创建新的请求返回
func NewCommonResponse(config *config.Config, request *requests.CommonRequest) (response *CommonResponse) {
	c := &CommonResponse{}
	c.Config = config
	c.Request = request
	return c
}

// GetHttpContentJson 获取 JSON 数据
func (res *CommonResponse) GetHttpContentJson() string {
	return res.json
}

// GetHttpContentMap 获取 MAP 数据
func (res *CommonResponse) GetHttpContentMap() (mxj.Map, error) {
	return mxj.NewMapJson([]byte(res.json))
}

// GetVerifySignDataMap 获取 GetVerifySignDataMap 校验后数据数据
func (res *CommonResponse) GetVerifySignDataMap() (m mxj.Map, err error) {
	return res.GetHttpContentMap()
}

// SetHttpContent 设置请求信息
func (res *CommonResponse) SetHttpContent(httpContent []byte, dataType string) {
	res.httpContent = httpContent
	switch dataType {
	case "xml":
		mv, _ := mxj.NewMapXml(res.httpContent) // unmarshal
		var str interface{}
		if _, ok := mv["xml"]; ok { //去掉 xml 外层
			str = mv["xml"]
		} else {
			str = mv
		}
		jsonStr, _ := json.Marshal(str)
		res.json = string(jsonStr)
	case "string":
		res.json = string(res.httpContent)
	}
}
