package applyLogin

import (
	"RAS/myRedis"
	"RAS/myUtils"
	"RAS/personDB"
	"RAS/serverDB"
	"RAS/toVncServer"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var TheDB *sqlx.DB
var TheRedis *myRedis.Redis

func PostApplyLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var applyInfo ApplyInfo
	var loginInfo LoginInfo

	applyInfo, loginInfo = extractApplyInfo(r)
	if loginInfo.RetCode != 0 {
		writeResponse(loginInfo, &w)
		return
	}

	//防频繁登录：先在redis缓存中进行检索
	blockTime := time.Minute * 3
	if TheRedis.IsExist(applyInfo.User) {
		loginInfo = LoginInfo{
			ServerInfo: "",
			RetCode:    -2,
			Msg: fmt.Sprintf("检测到您短时间内频繁登录，请%d分钟后再试!\n",
				myUtils.FloatToInt(blockTime.Minutes())),
		}
		writeResponse(loginInfo, &w)
		return
	}

	//从数据库中查询该用户
	//该用户是否存在, 该用户是否有分配服务器
	var person personDB.Person
	person, loginInfo = queryPerson(applyInfo.User, TheDB)
	if loginInfo.RetCode != 0 {
		writeResponse(loginInfo, &w)
		return
	}

	//从数据库中查询该服务器
	var server serverDB.ServerInfo
	server, loginInfo = queryServer(person.ServerName, TheDB)
	if loginInfo.RetCode != 0 {
		writeResponse(loginInfo, &w)
		return
	}

	// 查询ssh服务器，并解析端口号
	sshServer := toVncServer.DefaultSshServerInfo()
	hostPort := server.SshIP
	// Regular expression for IPv4 addresses
	regex := `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)?(:([0-9]{1,5}))?$`

	match, _ := regexp.MatchString(regex, hostPort)

	if !match {
		fmt.Println("Invalid IPv4 address or port")
		log.Printf("Invalid IPv4 address or port: %s", hostPort)
		loginInfo := LoginInfo{
			RetCode: -1,
			Msg:     "服务器SSH地址及端口号格式错误，请联系管理员！",
		}
		writeResponse(loginInfo, &w)
		return
	}

	re := regexp.MustCompile(`:([0-9]{1,5})$`)
	matches := re.FindStringSubmatch(hostPort)

	host := hostPort
	port := "22" // default port

	if len(matches) > 1 {
		host = hostPort[:len(hostPort)-len(matches[0])]
		port = matches[1]
	}

	sshServer.Host = host
	sshServer.Port, _ = strconv.Atoi(port)
	sshServer.Password = server.Password
	log.Printf("sshServer ip = %s, port = %d", sshServer.Host, sshServer.Port)

	//修改并获取密码
	passwd, err := toVncServer.ModifyVncPassword(person, sshServer)
	if err != nil {
		log.Printf("ModifyVncPassword error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		writeResponse(loginInfo, &w)
		return
	}
	log.Printf("ModifyVncPassword passwd = %s", passwd)

	//发送短信
	timeout := time.Minute * 10
	if sendShortMessage(person, passwd, timeout).RetCode != 0 {
		writeResponse(loginInfo, &w)
		return
	}

	//超时重置密码
	time.AfterFunc(timeout, func() {
		passwd, err := toVncServer.ModifyVncPassword(person, sshServer)
		if err != nil {
			log.Printf("timeout reset ModifyVncPassword error: %v", err)
		}
		log.Printf("timeout reset ModifyVncPassword passwd = %s", passwd)
	})

	//防频繁登录，在redis中记录当前用户
	rAddr := r.RemoteAddr
	fmt.Printf("Remote IP address is: %s \n", rAddr)
	TheRedis.Set(applyInfo.User, rAddr, blockTime)

	//一切成功
	log.Printf("all process success")
	//填写服务器地址+vnc桌面号
	loginInfo.ServerInfo = fmt.Sprintf("%s:%d", server.IP, person.VncDisplay)
	loginInfo.RetCode = 0
	writeResponse(loginInfo, &w)
	return
}

func extractApplyInfo(r *http.Request) (applyInfo ApplyInfo, loginInfo LoginInfo) {
	contentLength := r.ContentLength
	body := make([]byte, contentLength)
	n, err := r.Body.Read(body)
	if err != nil && err != io.EOF {
		log.Printf("loginInfo error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		return
	}
	log.Printf("r.Body.Read %d bytes: %s\n", n, body)

	err = json.Unmarshal(body, &applyInfo)
	if err != nil {
		log.Printf("ApplyInfo json.Unmarshal error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		return
	}
	log.Printf("User:%s\n", applyInfo.User)
	log.Printf("ISP:%s\n", applyInfo.ISP)
	log.Printf("PCInfo:%s\n", applyInfo.PCInfo)
	return
}

func queryPerson(userID string, db *sqlx.DB) (person personDB.Person, loginInfo LoginInfo) {
	var err error
	person, err = personDB.QueryPerson(userID, db)
	if err != nil {
		log.Printf("personDB.QueryPerson =%s error: %v", userID, err)
		loginInfo.RetCode = -1
		loginInfo.Msg = "无法找到该用户！"
		return
	}

	if len(person.ServerName) == 0 {
		log.Printf("person %s has no server resource", person.UserID)
		loginInfo.RetCode = -1
		loginInfo.Msg = "该用户无权访问服务器资源！"
		return
	}

	log.Printf("person %s info query success: %v", person.UserID, person)
	return
}

func queryServer(serverName string, db *sqlx.DB) (server serverDB.ServerInfo, loginInfo LoginInfo) {
	var err error
	server, err = serverDB.QueryServer(serverName, db)
	if err != nil {
		log.Printf("serverDB.QueryServer =%s error: %v", serverName, err)
		loginInfo.RetCode = -1
		loginInfo.Msg = "无法找到为您配置的服务器，请联系管理员！"
		return
	}
	return
}

func sendShortMessage(person personDB.Person, passwd string, timeout time.Duration) (loginInfo LoginInfo) {
	msg := DefaultMessageContent()
	phoneNumber := "+86" + person.Mobile
	msg.PhoneNumberSet = []*string{&phoneNumber}
	timeoutStr := fmt.Sprintf("%d", myUtils.FloatToInt(timeout.Minutes()))
	msg.TemplateParamSet = []*string{&person.UserName, &passwd, &timeoutStr}
	result, err := SendSMS(msg)
	if err != nil {
		log.Printf("send short message error: %v", err)
		loginInfo.RetCode = -1
		loginInfo.Msg = ""
		return
	}
	log.Printf("send short message success: %s", result)
	return
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
