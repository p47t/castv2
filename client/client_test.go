package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client := Client{Host: "10.0.1.4", Port: 8009}

	err := client.Connect()
	assert.NoError(t, err)

	time.Sleep(3 * time.Second)
}
