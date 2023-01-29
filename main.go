package main

import (
	"xl-app/svc"
	"xl-app/ui"
)

func main() {
	app := ui.NewApp(svc.NewDynamoService())
	app.RenderApp()
}
