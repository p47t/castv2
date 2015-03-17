package client

import (
	"time"

	"github.com/yinghau76/mdns"
)

type ChromecastInfo struct {
	Host string
	Port int
}

// https://github.com/jloutsenhizer/CR-Cast/wiki/Chromecast-Implementation-Documentation-WIP
func SearchChromecast() <-chan ChromecastInfo {
	entries := make(chan *mdns.ServiceEntry, 4)
	go func() {
		mdns.Query(&mdns.QueryParam{
			Service: "_googlecast._tcp",
			Domain:  "local",
			Timeout: 30 * time.Second,
			Entries: entries,
		})
	}()

	ret := make(chan ChromecastInfo)
	go func() {
		for entry := range entries {
			ret <- ChromecastInfo{
				Host: entry.Addr.String(),
				Port: entry.Port,
			}
		}
	}()
	return ret
}
