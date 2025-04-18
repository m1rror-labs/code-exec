package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestReadLogs(t *testing.T) {
	t.Skip()
	godotenv.Load("../../../.env")
	url := os.Getenv("DATABASE_URL")
	rep := InitRepository(url)

	logs, err := rep.ReadTransactionLogMessages().
		TeamID(uuid.MustParse("2b55bd68-dee3-40b6-ac3f-0d71f5519679")).
		BlockchainID(uuid.MustParse("41335c33-c715-4a07-9c55-4818ef900a97")).
		Paginate(1, 10).
		OrderCreatedAt("DESC").
		Execute(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Fatal(logs)
}
