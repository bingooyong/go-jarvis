package deploy

import (
	"fmt"
	"github.com/lvyong1985/go-jarvis/funcs"
)

type Link struct {
	sourcePath string
	targetPath string
	name       string "创建软链接"
}

func (l *Link) getName() string {
	return l.name
}

func (l *Link) exec(s *funcs.SSH) error {
	common := "ln -vnsf " + l.sourcePath + " " + l.targetPath
	return s.ExecCmd(common)
}

func (l *Link) String() string {
	return fmt.Sprintf("~~~shell \n ln -vnsf %s %s \n~~~", l.sourcePath, l.targetPath)
}
