package client

import (
	"time"

	"github.com/armon/mdns"
	"github.com/davecgh/go-spew/spew"
)

type ChromecastInfo struct {
	Host string
	Port int
}

func SearchChromecast() chan ChromecastInfo {
	entries := make(chan *mdns.ServiceEntry, 4)
	go func() {
		mdns.Query(&mdns.QueryParam{
			Service: "_googlecast._tcp.local",
			Domain:  "local",
			Timeout: 30 * time.Second,
			Entries: entries,
		})
	}()

	ret := make(chan ChromecastInfo)
	go func() {
		for entry := range entries {
			spew.Dump(entry)
			ret <- ChromecastInfo{
				Host: entry.Addr.String(),
				Port: entry.Port,
			}
		}
	}()
	return ret
}
