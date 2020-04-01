package toVncServer

import (
	"RAS/personDB"
	"testing"
)

func TestModifyVncPassword(t *testing.T) {
	p := personDB.Person{
		UserID:       "",
		UserName:     "",
		Department:   "",
		Mobile:       "",
		ServerName:   "",
		VncDisplay:   25,
		ServerUser:   "vncu25",
		ServerPasswd: "",
	}
	s := DefaultSshServerInfo()
	s.Host = ""
	s.Password = ""

	passwd, err := ModifyVncPassword(p, s)
	if err != nil {
		t.Errorf("ModifyVncPassword error: %v", err)
	}
	t.Logf("ModifyVncPassword success: %s", passwd)

}
