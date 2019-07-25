package main

import (
	uberfx "code.uber.internal/go/uberfx.git"
	yarpcfx "code.uber.internal/go/yarpcfx.git"
	"code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world/helloworldfx"
	"code.uber.internal/wonsoh/hello-world/handler"
	"go.uber.org/fx"
)

func main() {
	// New to Fx? Brush up at t.uber.com/fx.
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		uberfx.Module,
		yarpcfx.Module,
		fx.Provide(
			helloworldfx.Server(),
			handler.NewHelloWorld,
		),
	)
}
