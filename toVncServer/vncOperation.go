package toVncServer

import (
	"RAS/myUtils/genPasswd"
	"RAS/personDB"
	"fmt"
	"log"
)

func ModifyVncPassword(person personDB.Person, server SshServerInfo) (passwd string, err error) {
	genPasswd.RandomSeed()
	passwd = genPasswd.GeneratePasswd(8, genPasswd.FlagNumber)

	cmdFormat := `/root/ras/setvncpasswd %s %s`
	cmd := fmt.Sprintf(cmdFormat, person.ServerUser, passwd)

	result, err := SshTo(server, cmd)
	if err != nil {
		return "", fmt.Errorf("SshTo server %s error:%v", server.Host, err)
	}
	log.Printf("SshTo server %s success:%v", server.Host, result)
	return passwd, nil
}
