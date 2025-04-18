package codeexec

import (
	"code-exec/pkg"
	"context"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type ExecuteCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

type LogWithUrl struct {
	ID                   uuid.UUID       `json:"id"`
	CreatedAt            time.Time       `json:"created_at"`
	Url                  string          `json:"url"`
	TransactionSignature string          `json:"transaction_signature"`
	Log                  string          `json:"log"`
	Index                int             `json:"index"`
	Transaction          pkg.Transaction `json:"transaction"`
}

func RunCode(ctx context.Context, code string, codeExecutor pkg.CodeExecutor) (string, error) {
	return codeExecutor.ExecuteCode(code)
}

func parseEngineUrl(code string) (uuid.UUID, bool) {
	re := regexp.MustCompile(`(https://engine\.mirror\.ad/rpc/|http://localhost:8899/rpc/)([0-9a-fA-F-]{36})`)
	match := re.FindStringSubmatch(code)
	if len(match) < 3 {
		return uuid.Nil, false
	}

	uuidStr := match[2]
	engineID, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, false
	}
	return engineID, true
}
