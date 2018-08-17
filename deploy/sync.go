package deploy

import (
	"fmt"
	"github.com/lvyong1985/go-jarvis/funcs"
)

type Sync struct {
	sourcePath string
	targetPath string
	name       string "同步文件"
}

func (l *Sync) getName() string {
	return l.name
}

func (l *Sync) exec(s *funcs.SSH) error {
	return s.Put(l.sourcePath, l.targetPath)
}

func (l *Sync) String() string {
	return fmt.Sprintf("~~~shell \n sftp> put %s %s \n~~~", l.sourcePath, l.targetPath)
}
