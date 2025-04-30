package codeexechandlers

import (
	"code-exec/pkg"
	"code-exec/pkg/services/codeexec"
	"log"

	"github.com/gin-gonic/gin"
)

func BuildAndDeployAnchor(c *gin.Context, deps pkg.Dependencies) {
	var request codeexec.BuildProgramRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := codeexec.BuildAndLoadProgram(
		c,
		request.Code,
		request.ProgramID,
		request.BlockchainID,
		deps.AnchorRuntime,
		deps.RpcEngine,
	)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Program built and loaded successfully"})
}

func BuildAndTestAnchor(c *gin.Context, deps pkg.Dependencies) {
	var request codeexec.BuildAndTestProgramRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := codeexec.BuildAndTestProgram(
		c,
		request.Code,
		request.ProgramID,
		request.BlockchainID,
		request.TestCode,
		deps.AnchorRuntime,
		deps.RpcEngine,
	)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Program built, loaded, and tested successfully", "result": result})
}
