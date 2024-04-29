package repository

import (
	"context"
)

var _ AuthRepo = (*BobAuthRepo)(nil)

type (
	AuthRepo interface {
		GetLogin(ctx context.Context) error
	}
	BobAuthRepo struct {
	}
)

func (r BobAuthRepo) GetLogin(ctx context.Context) error {
	return nil
}

func NewBobAuthRepo() AuthRepo {

	return &BobAuthRepo{}
}
