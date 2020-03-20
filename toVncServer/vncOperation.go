package toVncServer

import (
	"RAS/personDB"
	"fmt"
	"regexp"
)

func ModifyVncPassword(person personDB.Person, server SshServerInfo) (passwd string, err error) {
	cmdFormat := `su - %s -c "/opt/TurboVNC/bin/vncpasswd -o -display :%d"`
	cmd := fmt.Sprintf(cmdFormat, person.ServerUser, person.VncDisplay)
	fmt.Println(cmd)
	result, err := SshTo(server, cmd)

	//test
	//reg := `^Full control one-time password: (\S+)$`
	r := regexp.MustCompile(`Full control one-time password: (\S+)`)
	if err != nil {
		return "", fmt.Errorf("SshTo server %s error:%v", server.Host, err)
	}

	if r.MatchString(result) == false {
		return "", fmt.Errorf("server %s result %s, match failed", server.Host, result)
	}

	ms := r.FindStringSubmatch(result)
	if len(ms) != 2 {
		return "", fmt.Errorf("server %s result %s, match len error", server.Host, result)
	}

	return ms[1], nil
	return
}
