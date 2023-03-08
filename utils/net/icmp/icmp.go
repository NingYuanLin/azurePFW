package icmp

import (
	"github.com/go-ping/ping"
	"time"
)

func CanReach(addr string) bool {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return false
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 3
	err = pinger.Run()
	if err != nil {
		return false
	}
	statistics := pinger.Statistics()
	if statistics.PacketsRecv == 0 {
		return false
	}
	return true
}
