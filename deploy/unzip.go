package deploy

import (
	"fmt"
	"github.com/lvyong1985/go-jarvis/funcs"
	"path/filepath"
)

type Unzip struct {
	path string
	file string
	name string "解压文件"
}

func (l *Unzip) getName() string {
	return l.name
}

func (l *Unzip) exec(s *funcs.SSH) error {
	return s.ExecMulti("cd "+l.path, decompressFileCommand(l.file))
}

func (l *Unzip) String() string {
	return fmt.Sprintf("~~~shell \n cd %s;%s \n~~~", l.path, decompressFileCommand(l.file))
}

func decompressFileCommand(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".gz":
		return "tar xzf " + filepath.Base(filename)
	case ".tar":
		return "tar xf " + filepath.Base(filename)
	case ".zip":
		return "unzip -qo " + filepath.Base(filename)
	default:
		return ""
	}
}
