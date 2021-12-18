package auth

import (
	"context"

	kc "github.com/ory/kratos-client-go"
)

type customKey int

var identityKey customKey

func NewIdentityContext(ctx context.Context, identity *kc.Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}
func IdentityFromContext(ctx context.Context) (*kc.Identity, bool) {
	identity, ok := ctx.Value(identityKey).(*kc.Identity)
	return identity, ok
}
