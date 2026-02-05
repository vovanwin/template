package main

import (
	"go.uber.org/fx"
)

func inject() fx.Option {
	options := []fx.Option{}
	return fx.Options(options...)
}

func main() {
	app := fx.New(inject())
	app.Run()
}
