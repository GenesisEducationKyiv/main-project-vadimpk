package main

import (
	"github.com/vadimpk/gses-2023/config"
	"github.com/vadimpk/gses-2023/internal/app"
)

func main() {
	cfg := config.Get()
	app.Run(cfg)
}
