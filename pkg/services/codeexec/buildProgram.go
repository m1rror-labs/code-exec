package codeexec

import (
	"code-exec/pkg"
	"context"
	"log"

	"github.com/google/uuid"
)

type BuildProgramRequest struct {
	Code         string    `json:"code" binding:"required"`
	ProgramID    string    `json:"program_id" binding:"required"`
	BlockchainID uuid.UUID `json:"blockchain_id" binding:"required"`
}

type BuildAndTestProgramRequest struct {
	Code         string    `json:"code" binding:"required"`
	ProgramID    string    `json:"program_id" binding:"required"`
	BlockchainID uuid.UUID `json:"blockchain_id" binding:"required"`
	TestCode     string    `json:"test_code" binding:"required"`
}

func BuildAndLoadProgram(
	ctx context.Context,
	code string,
	programID string,
	blockchainID uuid.UUID,
	programBuilder pkg.ProgramBuilder,
	rpcEngine pkg.RpcEngine,
) error {
	programBinary, _, err := programBuilder.BuildProgram(code, true)
	if err != nil {
		return err
	}
	if err := rpcEngine.LoadProgram(ctx, blockchainID, programID, programBinary); err != nil {
		log.Println("error loading program", err)
		return err
	}
	return nil
}

func BuildAndTestProgram(
	ctx context.Context,
	code string,
	programID string,
	blockchainID uuid.UUID,
	testCode string,
	programBuilder pkg.ProgramBuilder,
	rpcEngine pkg.RpcEngine,
) (string, error) {
	programBinary, codeID, err := programBuilder.BuildProgram(code, false)
	if err != nil {
		return "", err
	}
	if err := rpcEngine.LoadProgram(ctx, blockchainID, programID, programBinary); err != nil {
		log.Println("error loading program", err)
		return "", err
	}
	result, err := programBuilder.TestCode(testCode, blockchainID, codeID)
	if err != nil {
		log.Println("error testing program", err)
		return "", err
	}

	return result, nil
}
