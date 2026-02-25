package main

import (
	"server/internal/app"
	"server/internal/core"
	"server/internal/core/shared/types"
)

var (
	version     = "0.0.0"
	buildNumber = "0"
)

func main() {
	info := types.App{
		Name:        "Warehouse Management System",
		Version:     version,
		BuildNumber: buildNumber,
	}
	hooks := make(chan types.Hook)
	context := core.NewContext(info, hooks)
	application := app.New(context, hooks)
	if err := application.Run(); err != nil {
		context.Logger().Fatal(err)
	}
}
