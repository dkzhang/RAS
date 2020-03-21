package applyLogin

import (
	"RAS/myUtils"
	"RAS/personDB"
	"RAS/toVncServer"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

func PostApplyLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	loginInfo := LoginInfo{
		RetCode: 0,
	}

	len := r.ContentLength
	body := make([]byte, len)
	n, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("loginInfo error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		WriteResponse(loginInfo, &w)
		return
	}
	log.Printf("r.Body.Read %d bytes: %s\n", n, body)

	applyInfo := ApplyInfo{}
	err = json.Unmarshal(body, &applyInfo)
	if err != nil {
		log.Printf("ApplyInfo json.Unmarshal error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		WriteResponse(loginInfo, &w)
		return
	}
	log.Printf("User:%s\n", applyInfo.User)
	log.Printf("ISP:%s\n", applyInfo.ISP)
	log.Printf("PCInfo:%s\n", applyInfo.PCInfo)

	//从数据库中查询该用户的服务器信息
	//获取密码
	filename := "/IdKey/VncServer/server1.json"
	idKey, err := myUtils.LoadIdKey(filename)
	if err != nil {
		log.Printf("myUtils.LoadIdKey for server1 error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		WriteResponse(loginInfo, &w)
		return
	}
	loginInfo.ServerInfo = idKey.SecretId + ":25"

	person := personDB.Person{
		UserName:   "张俊",
		Mobile:     "15383026353",
		VncDisplay: 25,
		ServerUser: "vncu25",
	}

	s := toVncServer.DefaultSshServerInfo()
	s.Host = idKey.SecretId
	s.Password = idKey.SecretKey

	passwd, err := toVncServer.ModifyVncPassword(person, s)
	if err != nil {
		log.Printf("ModifyVncPassword error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		WriteResponse(loginInfo, &w)
		return
	}
	log.Printf("ModifyVncPassword result = %s", passwd)

	//发送短信
	msg := DefaultMessageContent()
	phoneNumber := "+86" + person.Mobile
	msg.PhoneNumberSet = []*string{&phoneNumber}
	timeoutStr := "15"
	msg.TemplateParamSet = []*string{&person.UserName, &passwd, &timeoutStr}
	retsult, err := SendSMS(msg)
	if err != nil {
		log.Printf("send short message error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		WriteResponse(loginInfo, &w)
		return
	}
	log.Printf("send short message success: %s", retsult)

	//一切成功
	log.Printf("all process success: %s", retsult)
	WriteResponse(loginInfo, &w)
	return
}

func WriteResponse(loginInfo LoginInfo, w *http.ResponseWriter) {
	//填写响应
	log.Printf("LoginInfo = %v", loginInfo)

	loginInfoJson, err := json.Marshal(loginInfo)
	if err != nil {
		log.Printf("loginInfo error: %v", err)
		(*w).WriteHeader(500)
	}
	fmt.Fprintf(*w, "%s", string(loginInfoJson))
}

type LoginInfo struct {
	ServerInfo string
	RetCode    int
	Msg        string
}

type ApplyInfo struct {
	User   string
	ISP    string
	PCInfo string
}
