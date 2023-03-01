package main

import (
	"context"
	"xl-app/db"
	"xl-app/svc"
	"xl-app/ui"
)

func main() {
	ctx := context.TODO()
	db := db.New(ctx)
	app := ui.NewApp(svc.NewDynamoService(db))
	app.RenderApp()
}
