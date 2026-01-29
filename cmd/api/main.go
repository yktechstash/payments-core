package main

import (
	"github.com/payments-core/internal/config"
	"github.com/payments-core/internal/router"
)

func main() {
	config.MustInit()
	router.Start()
}
