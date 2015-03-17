package client

import (
	"testing"
)

func TestSearchingChromecast(t *testing.T) {
	<-SearchChromecast()
}
