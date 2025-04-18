package postgres

import (
	"code-exec/pkg"
	"context"

	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var _ pkg.TransactionRepo = &repository{}

type transactionLogMessagesReader struct {
	selectQuery            *bun.SelectQuery
	transactionLogMessages *[]pkg.TransactionLogMessage
}

func (r *repository) ReadTransactionLogMessages() pkg.TransactionLogMessagesReader {
	var transactionLogMessages []pkg.TransactionLogMessage = []pkg.TransactionLogMessage{}
	return &transactionLogMessagesReader{selectQuery: r.db.NewSelect().Model(&transactionLogMessages), transactionLogMessages: &transactionLogMessages}
}

func (r *transactionLogMessagesReader) Execute(ctx context.Context) ([]pkg.TransactionLogMessage, error) {
	err := r.selectQuery.Scan(ctx)
	return *r.transactionLogMessages, err
}

func (r *transactionLogMessagesReader) ExecuteWithCount(ctx context.Context) ([]pkg.TransactionLogMessage, int, error) {
	count, err := r.selectQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	transactionLogMessages, err := r.Execute(ctx)
	return transactionLogMessages, count, err
}

func (r *transactionLogMessagesReader) ExecuteOne(ctx context.Context) (pkg.TransactionLogMessage, error) {
	err := r.selectQuery.Limit(1).Scan(ctx)
	if err != nil {
		return pkg.TransactionLogMessage{}, err
	}
	if len(*r.transactionLogMessages) == 0 {
		return pkg.TransactionLogMessage{}, pkg.ErrNotFound
	}
	return (*r.transactionLogMessages)[0], err
}

func (r *transactionLogMessagesReader) TeamID(teamID uuid.UUID) pkg.TransactionLogMessagesReader {
	r.selectQuery = r.selectQuery.Relation("Transaction.Blockchain").
		Where("transaction__blockchain.team_id = ?", teamID)
	return r
}

func (r *transactionLogMessagesReader) BlockchainID(blockchainID uuid.UUID) pkg.TransactionLogMessagesReader {
	r.selectQuery = r.selectQuery.Relation("Transaction").
		Where("transaction.blockchain = ?", blockchainID)
	return r
}

func (r *transactionLogMessagesReader) Paginate(page int, limit int) pkg.TransactionLogMessagesReader {
	offset := (page - 1) * limit
	r.selectQuery = r.selectQuery.Offset(offset).Limit(limit)
	return r
}

func (r *transactionLogMessagesReader) OrderCreatedAt(order string) pkg.TransactionLogMessagesReader {
	r.selectQuery = r.selectQuery.Order("transaction_log_message.created_at " + order)
	return r
}

func (r *transactionLogMessagesReader) Between(start, end time.Time) pkg.TransactionLogMessagesReader {
	r.selectQuery = r.selectQuery.Where("transaction_log_message.created_at BETWEEN ? AND ?", start, end)
	return r
}
