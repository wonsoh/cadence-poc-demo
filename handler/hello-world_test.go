package handler

import (
	"context"
	"fmt"
	"testing"

	helloworldgen "code.uber.internal/wonsoh/hello-world/.gen/go/wonsoh/hello-world/hello_world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHelloWorldHello(t *testing.T) {
	helloWorld := NewHelloWorld(zap.NewNop().Sugar())

	name := "E. Honda"
	message := fmt.Sprintf("Hello, %v!", name)

	response, err := helloWorld.Hello(context.Background(), &helloworldgen.HelloRequest{Name: &name})
	require.NoError(t, err)

	assert.Equal(t, message, response.GetMessage())
}
