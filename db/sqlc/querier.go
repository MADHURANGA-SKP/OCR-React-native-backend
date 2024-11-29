// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateImageConversion(ctx context.Context, arg CreateImageConversionParams) (ImageConversion, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUsers(ctx context.Context, arg CreateUsersParams) (User, error)
	DeleteImageConversion(ctx context.Context, conversionID int64) error
	DeleteUsers(ctx context.Context, userID int64) error
	GetImageConversionByID(ctx context.Context, conversionID int64) (ImageConversion, error)
	GetImageConversionsByUser(ctx context.Context, userID int32) ([]ImageConversion, error)
	GetSession(ctx context.Context, sessionID uuid.UUID) (Session, error)
	GetUser(ctx context.Context, userID int64) (GetUserRow, error)
	GetUserID(ctx context.Context, userID int64) (GetUserIDRow, error)
	GetUsers(ctx context.Context, userName string) (User, error)
	UpdateImageConversion(ctx context.Context, arg UpdateImageConversionParams) (ImageConversion, error)
	UpdateUsers(ctx context.Context, arg UpdateUsersParams) (User, error)
}

var _ Querier = (*Queries)(nil)
