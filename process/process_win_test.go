package process

import (
	"fmt"
	"testing"
)

func TestFindPortPid(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
	}{
		{"2222", args{2222}},
	}
	for _, tt := range tests {
		res := FindPortPid(tt.args.port)
		fmt.Println(res)
	}
}

func TestKillPid(t *testing.T) {
	pid := FindPortPid(2222)
	KillPid(pid)
}
