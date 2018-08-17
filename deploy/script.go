package deploy

import (
	"fmt"
	"github.com/lvyong1985/go-jarvis/funcs"
	"strings"
)

type Script struct {
	workPath   string
	execScript string
	name       string "执行脚本"
}

var (
	exportCommand = "export TERM=xterm"
	sourceCommand = "source /etc/profile"
	bashProfile   = "source ~/.bash_profile"
	bashRc        = "source ~/.bashrc"
)

func (l *Script) getName() string {
	return l.name
}

func (l *Script) exec(s *funcs.SSH) error {
	if strings.TrimSpace(l.execScript) == "" {
		return nil
	}
	currentFolder := "cd " + l.workPath
	return s.ExecMulti(exportCommand, sourceCommand, bashProfile, bashRc, currentFolder, strings.TrimSpace(l.execScript), "sleep 1")
}

func (l *Script) String() string {
	if strings.TrimSpace(l.execScript) == "" {
		return ""
	}
	return fmt.Sprintf("~~~shell \n cd %s;%s; \n~~~", l.workPath, l.execScript)
}
