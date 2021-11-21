package resolvers

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	pool *pgxpool.Pool
}

func NewResolver(connstring string) (*Resolver, error) {
	db, err := pgxpool.Connect(context.Background(), connstring)
	if err != nil {
		return nil, err
	}
	return &Resolver{
		pool: db,
	}, nil
}

func NewDemoResolver() *Resolver {
	return &Resolver{}
}
