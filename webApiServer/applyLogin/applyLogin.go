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

	applyInfo := ApplyInfo{}

	if extractApplyInfo(r, &applyInfo, &loginInfo).RetCode != 0 {
		writeResponse(loginInfo, &w)
		return
	}

	//从数据库中查询该用户
	//该用户是否存在，该用户是否有分配服务器

	//从数据库中查询该服务器

	//填写服务器地址+vnc桌面号
	//修改密码

	//获取密码
	filename := "/IdKey/VncServer/server1.json"
	idKey, err := myUtils.LoadIdKey(filename)
	if err != nil {
		log.Printf("myUtils.LoadIdKey for server1 error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		writeResponse(loginInfo, &w)
		return
	}
	loginInfo.ServerInfo = idKey.SecretId + ":25"

	//personDB.QueryPerson(theDB)
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
		writeResponse(loginInfo, &w)
		return
	}
	log.Printf("ModifyVncPassword result = %s", passwd)

	//发送短信
	if sendShortMessage(person, passwd, &loginInfo).RetCode != 0 {
		writeResponse(loginInfo, &w)
		return
	}

	//一切成功
	log.Printf("all process success")
	writeResponse(loginInfo, &w)
	return
}

func extractApplyInfo(r *http.Request, pApplyInfo *ApplyInfo, pLoginInfo *LoginInfo) *LoginInfo {
	contentLength := r.ContentLength
	body := make([]byte, contentLength)
	n, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("loginInfo error: %v", err)
		pLoginInfo.RetCode = -1
		pLoginInfo.Msg = ""
		return pLoginInfo
	}
	log.Printf("r.Body.Read %d bytes: %s\n", n, body)

	err = json.Unmarshal(body, pApplyInfo)
	if err != nil {
		log.Printf("ApplyInfo json.Unmarshal error: %v", err)
		pLoginInfo.RetCode = -1
		pLoginInfo.Msg = ""
		return pLoginInfo
	}
	log.Printf("User:%s\n", pApplyInfo.User)
	log.Printf("ISP:%s\n", pApplyInfo.ISP)
	log.Printf("PCInfo:%s\n", pApplyInfo.PCInfo)
	return pLoginInfo
}

func sendShortMessage(person personDB.Person, passwd string, pLoginInfo *LoginInfo) *LoginInfo {
	msg := DefaultMessageContent()
	phoneNumber := "+86" + person.Mobile
	msg.PhoneNumberSet = []*string{&phoneNumber}
	timeoutStr := "15"
	msg.TemplateParamSet = []*string{&person.UserName, &passwd, &timeoutStr}
	result, err := SendSMS(msg)
	if err != nil {
		log.Printf("send short message error: %v", err)
		pLoginInfo.RetCode = -1
		pLoginInfo.Msg = ""
		return pLoginInfo
	}
	log.Printf("send short message success: %s", result)
	return pLoginInfo
}

func writeResponse(loginInfo LoginInfo, w *http.ResponseWriter) {
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
