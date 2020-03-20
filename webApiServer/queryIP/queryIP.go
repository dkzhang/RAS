package queryIP

import (
	"RAS/myUtils"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"strings"
	"time"
)

func calcAuthorization(source string, secretId string, secretKey string) (auth string, datetime string, err error) {
	timeLocation, _ := time.LoadLocation("Etc/GMT")
	datetime = time.Now().In(timeLocation).Format("Mon, 02 Jan 2006 15:04:05 GMT")
	signStr := fmt.Sprintf("x-date: %s\nx-source: %s", datetime, source)

	// hmac-sha1
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	auth = fmt.Sprintf("hmac id=\"%s\", algorithm=\"hmac-sha1\", headers=\"x-date x-source\", signature=\"%s\"",
		secretId, sign)

	return auth, datetime, nil
}

func urlencode(params map[string]string) string {
	var p = gourl.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}

func queryIP(ipAddress string) (ipInfo []byte, err error) {
	//Load IdKey from file
	filename := "/IdKey/QueryIP/config.json"
	idKey, err := myUtils.LoadIdKey(filename)
	if err != nil {
		return nil, fmt.Errorf("queryIP Load IdKey from file error: %v", err)
	}

	// 云市场分配的密钥Id
	secretId := idKey.SecretId
	// 云市场分配的密钥Key
	secretKey := idKey.SecretKey
	source := "market-22fhhos72"

	/////////////////////////////////////////////
	// 签名
	auth, datetime, _ := calcAuthorization(source, secretId, secretKey)

	// 请求方法
	method := "GET"
	// 请求头
	headers := map[string]string{"X-Source": source, "X-Date": datetime, "Authorization": auth}

	// 查询参数
	queryParams := make(map[string]string)
	queryParams["ip"] = ipAddress
	// body参数
	bodyParams := make(map[string]string)

	// url参数拼接
	url := "https://service-53e769xh-1255468759.ap-shanghai.apigateway.myqcloud.com/release/iparea"
	if len(queryParams) > 0 {
		url = fmt.Sprintf("%s?%s", url, urlencode(queryParams))
	}

	bodyMethods := map[string]bool{"POST": true, "PUT": true, "PATCH": true}
	var body io.Reader = nil
	if bodyMethods[method] {
		body = strings.NewReader(urlencode(bodyParams))
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("make http.NewRequest error: %v", err)
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do http.NewRequest error: %v", err)
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %v", err)
	}
	return bodyBytes, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//unmarshal JSON
type QueryIpApiResult struct {
	Error string               `json:"showapi_res_error"`
	ID    string               `json:"showapi_res_id"`
	Code  int                  `json:"showapi_res_code"`
	Body  QueryIpApiResultBody `json:"showapi_res_body"`
}

type QueryIpApiResultBody struct {
	Isp         string `json:"isp"`
	Ip          string `json:"ip"`
	Region      string `json:"region"`
	Lnt         string `json:"lnt"`
	County      string `json:"county"`
	EnNameShort string `json:"en_name_short"`
	Lat         string `json:"lat"`
	City        string `json:"city"`
	CityCode    string `json:"city_code"`
	Country     string `json:"country"`
	Continents  string `json:"continents"`
	EnName      string `json:"en_name"`
	RetCode     int    `json:"ret_code"`
	Msg         string
}

/*
{
"showapi_res_error": "",
"showapi_res_id": "f997f533fcf04879b6cb90357453227c",
"showapi_res_code": 0,
"showapi_res_body": {"isp":"电信","ip":"27.187.113.186","region":"河北","lnt":"115.482331","county":"","en_name_short":"CN","lat":"38.867657","city":"保定","city_code":"130600","country":"中国","continents":"亚洲","en_name":"China","ret_code":0}
}
*/

///////////////////////////////////////////////////////////////////////////////////////////////////
