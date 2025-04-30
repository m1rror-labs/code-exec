package anchor

import (
	"code-exec/pkg"
)

type runtime struct {
}

func NewRuntime() pkg.ProgramBuilder {
	return &runtime{}
}
