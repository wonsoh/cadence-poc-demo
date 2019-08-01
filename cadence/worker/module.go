package worker

import (
	"context"
	"time"

	"github.com/uber-go/dosa"
	"go.uber.org/fx"
)

// private type -- serves as key for context
type contextKey int

const (
	// PerpetualTimeout denotes "de facto" perpetual time out for long-running activities
	PerpetualTimeout = time.Hour * 24 * 365 * 30

	// DOSAClientContextKey is the key for retrieving DOSA client in background context
	DOSAClientContextKey = contextKey(iota)
)

type params struct {
	fx.In

	DOSAClient dosa.Client
}

type result struct {
	fx.Out

	ActivitiesContext context.Context `name:"activities_context"`
}

func registerDependencies(p params) result {
	ctx := context.WithValue(context.Background(), DOSAClientContextKey, p.DOSAClient)
	return result{
		ActivitiesContext: ctx,
	}
}

// RetrieveDOSAClientFromContext retries DOSA client from the parent context
func RetrieveDOSAClientFromContext(ctx context.Context) dosa.Client {
	return ctx.Value(DOSAClientContextKey).(dosa.Client)
}

// Module exports any cadence dependencies
var Module = fx.Provide(registerDependencies)
