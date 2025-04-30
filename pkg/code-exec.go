package pkg

import "github.com/google/uuid"

type CodeExecutor interface {
	ExecuteCode(code string) (string, error)
}

type ProgramBuilder interface {
	BuildProgram(code string, deleteArtifacts bool) ([]byte, string, error)
	TestCode(code string, blockchainID uuid.UUID, codeID string) (string, error)
}

type Err string

func (e Err) Error() string {
	return string(e)
}

const (
	ErrUnauthorized   = Err("Unauthorized")
	ErrTooManyApiKeys = Err("Too many api keys")
	ErrHttpRequest    = Err("HTTP request error")
	ErrNoApiKey       = Err("No api key")
	ErrNotFound       = Err("Not found")

	ErrInvalidPubkey       = Err("Invalid pubkey")
	ErrInvalidSignature    = Err("Invalid signature")
	ErrAccountNotFound     = Err("Account not found")
	ErrSettingAccount      = Err("Error setting account")
	ErrTransactionNotFound = Err("Transaction not found")
	ErrNoAccounts          = Err("No accounts")
)
