package icmp

import "testing"

func TestCanReach(t *testing.T) {
	addr := "8.8.8.8"
	ok := CanReach(addr)
	if ok == false {
		t.Fatal("无法icmp")
	}
}
