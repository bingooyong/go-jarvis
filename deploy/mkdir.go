package deploy

import (
	"fmt"
	"github.com/lvyong1985/go-jarvis/funcs"
)

type Mkdir struct {
	path string
	name string "创建目录"
}

func (l *Mkdir) getName() string {
	return l.name
}

func (l *Mkdir) exec(s *funcs.SSH) error {
	return s.ExecCmd("mkdir -p " + l.path)
}

func (l *Mkdir) String() string {
	return fmt.Sprintf("~~~shell \n mkdir -p %s \n~~~", l.path)
}
