package funcs

import (
	"testing"
	"os"
)

func TestExecCmd_Exec(t *testing.T) {
	cmd := NewExecCmd("/Users/lvyong/Code/Go/src/github.com/lvyong1985/go-jarvis/tmp", os.Stdout)
	cmd.Exec("git clone --depth 1 git@192.168.131.32:develop/FOOTSTONE/GateWay/Code/api-gateway-manager.git")
}
