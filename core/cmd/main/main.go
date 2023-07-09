package main

import (
	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/core/internal/app"
)

func main() {
	cfg := config.Get(".env")
	app.Run(cfg)
}
