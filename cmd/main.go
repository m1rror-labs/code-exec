package main

import (
	"code-exec/pkg/app"
	"code-exec/pkg/dependencies/postgres"
	"code-exec/pkg/dependencies/rpcengine"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	env := os.Getenv("ENV")
	repo := postgres.NewRepository(os.Getenv("DATABASE_URL"))
	rpcEngine := rpcengine.New(os.Getenv("ENGINE_URL"))

	app := app.NewApp(env, repo, rpcEngine)

	app.Run()
}
