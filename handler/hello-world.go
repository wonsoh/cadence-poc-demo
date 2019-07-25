package handler

import (
	"context"
	"fmt"

	zapfx "code.uber.internal/go/zapfx.git"
	helloworldgen "code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world"
	"code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world/helloworldserver"
	"go.uber.org/zap"
)

// NewHelloWorld creates the impl for the HelloWorld service in hello_world.thrift.
func NewHelloWorld(logger *zap.SugaredLogger) helloworldserver.Interface {
	return &helloWorld{logger: logger}
}

type helloWorld struct {
	logger *zap.SugaredLogger
}

func (h *helloWorld) Hello(ctx context.Context, request *helloworldgen.HelloRequest) (*helloworldgen.HelloResponse, error) {
	message := fmt.Sprintf("Hello, %v!", request.GetName())
	h.logger.Infow("hello called", zapfx.Trace(ctx), "message", message)

	return &helloworldgen.HelloResponse{Message: &message}, nil
}
