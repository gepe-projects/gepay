package main

import (
	"github.com/ilhamgepe/gepay/internal/wire"
)

func main() {
	app := wire.InitializeApps()
	app.Server.Run()
}
