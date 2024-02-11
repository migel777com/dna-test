package service

import (
	"context"
	"dna-test/models"
	"errors"
	"fmt"
	"time"
)

func (s *service) CreateAccount(ctx context.Context, account *models.Account) error {
	if account == nil {
		return errors.New("nil account")
	}

	err := s.Db.Create(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FreezeAccount(ctx context.Context, id string) error {
	var filter models.FilterParams
	filter.Filter = fmt.Sprintf(`id = '%v'`, id)

	var account models.Account
	err := s.Db.Get(ctx, filter, &account)
	if err != nil {
		return err
	}

	if account.IsDisabled {
		return errors.New("account already frozen")
	}

	account.IsDisabled = true
	err = s.Db.Update(ctx, filter, &account)
	if err != nil {
		return err
	}

	err = s.Cache.SetHash(ctx, models.AccountKey+id, account, time.Hour*24)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAccount(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	err := s.Cache.GetHash(ctx, models.AccountKey+id, &account)
	if models.IsErrNotFound(err) {
		var filter models.FilterParams
		filter.Filter = fmt.Sprintf(`id = '%v'`, id)

		err = s.Db.Get(ctx, filter, &account)
		if err != nil {
			return nil, err
		}

		err = s.Cache.SetHash(ctx, models.AccountKey+id, account, time.Hour*24)
		if err != nil {
			return nil, err
		}
	}

	return &account, nil
}
