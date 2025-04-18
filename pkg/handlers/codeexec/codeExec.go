package codeexechandlers

import (
	"code-exec/pkg"
	"code-exec/pkg/services/codeexec"
	"log"

	"github.com/gin-gonic/gin"
)

func ExecuteTypescript(c *gin.Context, deps pkg.Dependencies) {
	var request codeexec.ExecuteCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	output, logs, err := codeexec.RunCode(c, request.Code, deps.TsRuntime, deps.Repo)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "output": output, "logs": logs})
		return
	}

	c.JSON(200, gin.H{"output": output, "logs": logs})
}

func ExecuteRust(c *gin.Context, deps pkg.Dependencies) {
	var request codeexec.ExecuteCodeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	output, logs, err := codeexec.RunCode(c, request.Code, deps.RustRuntime, deps.Repo)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": err.Error(), "output": output, "logs": logs})
		return
	}

	c.JSON(200, gin.H{"output": output, "logs": logs})
}

func LoadTest(c *gin.Context, deps pkg.Dependencies) {
	if err := codeexec.LoadTestCodeExec(c, 10, 10); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Load test completed successfully"})
}
