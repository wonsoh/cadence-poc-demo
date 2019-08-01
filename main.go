package main

import (
	"code.uber.internal/devexp/cadencefx.git"
	"code.uber.internal/go/dosafx.git"
	"code.uber.internal/go/uberfx.git"
	"code.uber.internal/go/yarpcfx.git"
	"code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world/helloworldfx"
	_ "code.uber.internal/wonsoh/hello-world/cadence/activities"
	"code.uber.internal/wonsoh/hello-world/cadence/worker"
	_ "code.uber.internal/wonsoh/hello-world/cadence/workflows"
	"code.uber.internal/wonsoh/hello-world/entities"
	"code.uber.internal/wonsoh/hello-world/handler"
	"go.uber.org/fx"
)

func main() {
	// New to Fx? Brush up at t.uber.com/fx.
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		entities.Module,
		dosafx.Module,
		uberfx.Module,
		yarpcfx.Module,
		worker.Module,
		cadencefx.Module,
		fx.Provide(
			helloworldfx.Server(),
			handler.NewHelloWorld,
		),
	)
}
