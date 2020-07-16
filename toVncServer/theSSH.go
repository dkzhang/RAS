package toVncServer

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

type SshServerInfo struct {
	Host     string
	User     string
	Password string
	Type     string
	KeyPath  string
	Port     int

	Timeout time.Duration
}

func DefaultSshServerInfo() SshServerInfo {
	return SshServerInfo{
		Host:     "",
		User:     "root",
		Password: "",
		Type:     "password",
		KeyPath:  "",
		Port:     22,
		Timeout:  time.Minute,
	}
}

func SshTo(server SshServerInfo, cmd string) (result string, err error) {
	config := &ssh.ClientConfig{
		Timeout:         server.Timeout,
		User:            server.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if server.Type == "password" {
		config.Auth = []ssh.AuthMethod{ssh.Password(server.Password)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(server.KeyPath)}
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", server.Host, server.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("ssh.Dial to generate sshClient error: %v", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("create sshClient.NewSession error: %v", err)
	}
	defer session.Close()

	//执行远程命令
	combo, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("session.CombinedOutput error: %v", err)
	}

	return string(combo), nil
}

// unaccomplished！！！
func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath := ""

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
