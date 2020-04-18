package queryIP

import "testing"

func TestQueryIP(t *testing.T) {
	r, err := queryIP(" 27.187.116.121")
	if err != nil {
		t.Errorf("queryIP error: %v", err)
	}
	t.Logf("queryIP success: %v", r)
}
