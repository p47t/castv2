package main

import "simplypatrick.com/castv2/client"

func main() {
	ret := client.SearchChromecast()
	<-ret
}
