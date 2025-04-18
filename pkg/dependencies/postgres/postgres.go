package postgres

import (
	"code-exec/pkg"
	"context"
	"database/sql"
	"log"

	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type repository struct {
	db *bun.DB
}

func InitRepository(url string) repository {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url)))
	maxOpenConns := 100
	sqldb.SetMaxOpenConns(maxOpenConns)        // max connections set in supabase
	sqldb.SetMaxIdleConns(30)                  // 30 percent of max connections
	sqldb.SetConnMaxLifetime(15 * time.Minute) // refresh stale connections
	client := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	client.AddQueryHook(&QueryHook{})

	// client.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return repository{
		db: client,
	}
}

func NewRepository(url string) pkg.Repository {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url)))
	maxOpenConns := 100
	sqldb.SetMaxOpenConns(maxOpenConns)        // max connections set in supabase
	sqldb.SetMaxIdleConns(30)                  // 30 percent of max connections
	sqldb.SetConnMaxLifetime(15 * time.Minute) // refresh stale connections
	client := bun.NewDB(sqldb, pgdialect.New(), bun.WithDiscardUnknownColumns())
	client.AddQueryHook(&QueryHook{})

	// client.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	return &repository{
		db: client,
	}
}

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	duration := time.Since(event.StartTime)

	// slow query threshold
	slowThreshold := 5000 * time.Millisecond // 5 seconds

	if duration > slowThreshold {
		log.Printf("Slow query (%.2fms): %s", float64(duration.Milliseconds()), event.Query)
	}

	// log query failure
	if event.Err != nil {
		log.Printf("Query failed: %s | error: %v", event.Query, event.Err)
	}
}
