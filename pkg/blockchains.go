package pkg

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Blockchain struct {
	bun.BaseModel  `bun:"table:blockchains"`
	ID             uuid.UUID  `json:"id,omitempty" bun:"type:uuid,default:uuid_generate_v4(),pk"`
	CreatedAt      time.Time  `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	AirdropKeypair []byte     `json:"-" bun:",nullzero,notnull"`
	TeamID         uuid.UUID  `json:"team_id" bun:"type:uuid,notnull"`
	Label          *string    `json:"label" bun:",notnull"`
	Expiry         *time.Time `json:"expiry"`
}

type SolanaAccount struct {
	Address    string `json:"address"`
	Lamports   uint   `json:"lamports"`
	Data       []byte `json:"data"`
	Owner      string `json:"owner"`
	Executable bool   `json:"executable"`
	RentEpoch  uint   `json:"rentEpoch"`
}

type RpcEngine interface {
	CreateBlockchain(ctx context.Context, apiKey uuid.UUID, user_id *string, config *uuid.UUID) (uuid.UUID, error)
	DeleteBlockchain(ctx context.Context, apiKey uuid.UUID, id uuid.UUID) error
	ExpireBlockchains(ctx context.Context) error

	SetAccounts(ctx context.Context, blockchainID uuid.UUID, accounts []SolanaAccount, label *string, token_mint_auth *string) error
	LoadProgram(ctx context.Context, blockchainID uuid.UUID, programID string, programBinary []byte) error
}
