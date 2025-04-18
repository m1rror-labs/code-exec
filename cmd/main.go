package main

import (
	"code-exec/pkg/app"
	"code-exec/pkg/dependencies/rpcengine"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	env := os.Getenv("ENV")
	rpcEngine := rpcengine.New(os.Getenv("ENGINE_URL"))

	app := app.NewApp(env, rpcEngine)

	app.Run()
}
