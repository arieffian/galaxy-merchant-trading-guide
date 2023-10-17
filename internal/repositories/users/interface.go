package users

import (
	"context"

	"github.com/arieffian/roman-alien-currency/internal/models"
)

type GetUserByIdParams struct {
	UserId string
}

type GetUserByIdResponse struct {
	User models.User
}

type UserInterface interface {
	GetUserById(ctx context.Context, p GetUserByIdParams) (*GetUserByIdResponse, error)
}
