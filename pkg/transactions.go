package pkg

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel   `bun:"table:transactions"`
	ID              uuid.UUID   `json:"id" bun:"type:uuid,default:uuid_generate_v4(),pk"`
	CreatedAt       time.Time   `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	Version         string      `json:"version" bun:",notnull"`
	RecentBlockhash []byte      `json:"recent_blockhash" bun:",notnull"`
	Slot            int         `json:"slot" bun:",notnull"`
	BlockchainID    uuid.UUID   `json:"blockchain_id" bun:"blockchain,type:uuid,notnull"`
	Signature       string      `json:"signature" bun:",notnull"`
	Blockchain      *Blockchain `json:"blockchain,omitempty" bun:"rel:has-one,join:blockchain=id"`

	Logs []TransactionLogMessage `json:"logs" bun:"rel:has-many,join:signature=transaction_signature"`
}

type TransactionLogMessage struct {
	bun.BaseModel        `bun:"table:transaction_log_messages"`
	ID                   uuid.UUID   `json:"id" bun:"type:uuid,default:uuid_generate_v4(),pk"`
	CreatedAt            time.Time   `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	TransactionSignature string      `json:"transaction_signature" bun:",notnull"`
	Log                  string      `json:"log" bun:",notnull"`
	Index                int         `json:"index" bun:",notnull"`
	Transaction          Transaction `json:"transaction" bun:"rel:has-one,join:transaction_signature=signature"`
}

type TransactionRepo interface {
	// ReadTransaction() TransactionReader

	ReadTransactionLogMessages() TransactionLogMessagesReader
}

type TransactionLogMessagesReader interface {
	Execute(ctx context.Context) ([]TransactionLogMessage, error)
	ExecuteWithCount(ctx context.Context) ([]TransactionLogMessage, int, error)
	ExecuteOne(ctx context.Context) (TransactionLogMessage, error)

	TeamID(teamID uuid.UUID) TransactionLogMessagesReader
	BlockchainID(blockchainID uuid.UUID) TransactionLogMessagesReader
	Paginate(page int, limit int) TransactionLogMessagesReader
	Between(start time.Time, end time.Time) TransactionLogMessagesReader

	OrderCreatedAt(order string) TransactionLogMessagesReader
}
