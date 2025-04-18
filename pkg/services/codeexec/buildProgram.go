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

func BuildAndLoadProgram(
	ctx context.Context,
	code string,
	programID string,
	blockchainID uuid.UUID,
	programBuilder pkg.ProgramBuilder,
	rpcEngine pkg.RpcEngine,
) error {
	programBinary, err := programBuilder.BuildProgram(code)
	if err != nil {
		return err
	}
	if err := rpcEngine.LoadProgram(ctx, blockchainID, programID, programBinary); err != nil {
		log.Println("error loading program", err)
		return err
	}
	return nil
}
