package toVncServer

/*
func TestSshTo(t *testing.T) {
	serverInfo := DefaultSshServerInfo()

	serverInfo.Host = "192.168.10.28"
	serverInfo.User = "root"
	serverInfo.Password = "Zhang111111"
	serverInfo.Type = "password"

	cmd := `runuser -l dkzhang -c "/opt/TurboVNC/bin/vncpasswd -o -display :1"`
	//cmd := `runuser -l dkzhang -c "/opt/TurboVNC/bin/vncserver -kill :1"`
	result, err := SshTo(serverInfo, cmd)
	if err != nil {
		t.Errorf("ssh error: %v", err)
	} else {
		t.Logf("ssh result ==>\n%s", result)
	}
}
*/
