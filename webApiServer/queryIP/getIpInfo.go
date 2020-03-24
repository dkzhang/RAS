package queryIP

import (
	"RAS/myRedis"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"regexp"
	"time"
)

var TheRedis *myRedis.Redis

func GetIpInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ipInfoResult := QueryIpApiResultBody{}
	rAddr := r.RemoteAddr
	fmt.Printf("Remote IP address is: %s \n", rAddr)

	// 先在缓存中查询
	cacheTime := time.Hour * 24
	if TheRedis.IsExist(rAddr) {
		strJson := TheRedis.Get(rAddr).(string)
		log.Printf("find ip info <%s> from cache: %s", rAddr, strJson)
		fmt.Fprintf(w, "%s", strJson)
	}

	//((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)

	reg := regexp.MustCompile(`((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)`)

	if reg.MatchString(rAddr) == true {
		ipAddress := reg.FindString(rAddr)
		log.Printf("Server received your IP address: %s \n", ipAddress)

		ipInfo, err := queryIP(ipAddress)
		if err != nil {
			log.Printf("queryIP error: %v", err)
			ipInfoResult.RetCode = -1
			ipInfoResult.Msg = "IP地址查询错误"
		} else {
			log.Printf("IP info is:\n%s", string(ipInfo))
			ipInfoStruct := QueryIpApiResult{}
			if err := json.Unmarshal(ipInfo, &ipInfoStruct); err != nil {
				log.Printf("ip info unmarshal error: %v", err)
				ipInfoResult.RetCode = -1
				ipInfoResult.Msg = "IP地址查询错误"
			} else {
				//正常情况
				log.Printf("ip info unmarshal reslut: %v", ipInfoStruct)
				ipInfoResult = ipInfoStruct.Body
				ipInfoResult.RetCode = 0
				ipInfoResult.Msg = ""
			}
		}
	} else {
		log.Printf("Server received an illegal IP address: %s", rAddr)
		ipInfoResult.RetCode = -1
		ipInfoResult.Msg = "IP地址查询错误"
	}

	ipInfoResultJson, err := json.Marshal(ipInfoResult)
	if err != nil {
		log.Printf("ipInfoResult json.Marshal error: %v", err)
		w.WriteHeader(500)
	} else {
		strJson := string(ipInfoResultJson)
		fmt.Fprintf(w, "%s", strJson)
		//写入缓存
		TheRedis.Set(rAddr, strJson, cacheTime)
		log.Printf("write ip info <%s> in cache: %s", rAddr, strJson)
	}
	return
}
