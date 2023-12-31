package main

import (
	"github.com/vadimpk/gses-2023/crypto/config"
	"github.com/vadimpk/gses-2023/crypto/internal/app"
)

func main() {
	cfg := config.Get(".env")
	app.Run(cfg)
}
