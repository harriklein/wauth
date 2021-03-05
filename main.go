package main

import (
	"github.com/harriklein/wauth/app"
	"github.com/harriklein/wauth/config"
	"github.com/harriklein/wauth/log"
)

func main() {
	log.Init()

	app.Init()

	app.MapUrls()

	//datasource.Init()

	app.RunServer(config.BindAddress)

	app.Finish()
}
