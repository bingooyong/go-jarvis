package deploy

import "github.com/lvyong1985/go-jarvis/funcs"

type Command interface {
	getName() string
	exec(s *funcs.SSH) error
}
