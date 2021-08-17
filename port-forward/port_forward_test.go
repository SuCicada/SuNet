package port_forward

import (
	"os"
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	os.Args = strings.Split(". -r ubuntu.wsl2:22 -p 2222", " ")
	Start()
}
