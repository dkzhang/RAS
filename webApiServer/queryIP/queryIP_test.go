package queryIP

import (
	"encoding/json"
	"testing"
)

func TestQueryIP(t *testing.T) {
	r, err := queryIP("27.187.116.121")
	if err != nil {
		t.Errorf("queryIP error: %v", err)
	}
	t.Logf("queryIP success")

	ipInfoStruct := QueryIpApiResult{}
	if err := json.Unmarshal(r, &ipInfoStruct); err != nil {
		t.Errorf("ip info unmarshal error: %v", err)
	} else {
		//正常情况
		t.Logf("ip info unmarshal reslut: %v", ipInfoStruct)
	}
}
