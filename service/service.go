package service

import (
	"context"
	"dna-test/models"
)

type Service interface {
	CreateAccount(ctx context.Context, account *models.Account) error
	FreezeAccount(ctx context.Context, id string) error
	GetAccount(ctx context.Context, id string) (*models.Account, error)
}

type service struct {
	Db    models.DbClient
	Cache models.CacheClient
}

func NewService(db models.DbClient, cache models.CacheClient) Service {
	return &service{
		Db:    db,
		Cache: cache,
	}
}
