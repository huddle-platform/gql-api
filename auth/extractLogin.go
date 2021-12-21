package auth

import (
	"context"

	kc "github.com/ory/kratos-client-go"
)

type Identity struct {
	kc.Identity
	traits map[string]interface{}
}
type customKey int

var identityKey customKey

func NewIdentityContext(ctx context.Context, identity *kc.Identity) context.Context {
	return context.WithValue(ctx, identityKey, identity)
}
func IdentityFromContext(ctx context.Context) (*Identity, bool) {
	identity, ok := ctx.Value(identityKey).(*kc.Identity)
	if ok {
		traits, ok := identity.Traits.(map[string]interface{})
		if ok {
			return &Identity{Identity: *identity, traits: traits}, ok
		}
	}
	return nil, ok
}

func (i *Identity) GetTrait(key string) (string, bool) {
	val, ok := i.traits[key].(string)
	return val, ok
}
