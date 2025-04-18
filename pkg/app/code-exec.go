package app

import (
	codeexechandlers "code-exec/pkg/handlers/codeexec"

	"github.com/gin-gonic/gin"
)

func (a *App) AttachCodeExecRoutes() {
	a.engine.POST("/code-exec/typescript", func(c *gin.Context) {
		codeexechandlers.ExecuteTypescript(c, a.deps)
	})
	a.engine.POST("/code-exec/rust", func(c *gin.Context) {
		codeexechandlers.ExecuteRust(c, a.deps)
	})

	a.engine.POST("/code-exec/programs/anchor", func(c *gin.Context) {
		codeexechandlers.BuildAndDeployAnchor(c, a.deps)
	})
}
