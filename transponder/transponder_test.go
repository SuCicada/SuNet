package transponder

import (
	"testing"
)

func TestClient(t *testing.T) {
	Client("win.ip", 7001,
		4000, 1234)
}

func TestClient2(t *testing.T) {
	Client("win.ip", 7001,
		22, 9090)
}

func TestServer(t *testing.T) {
	Server(7001)
}
