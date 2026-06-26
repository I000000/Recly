//go:generate mockery --name TokenRepository --output ../../mocks --outpkg mocks --case underscore
package domain

import "context"

type TokenRepository interface {
	StoreRefreshToken(ctx context.Context, rt *RefreshToken) error
	GetRefreshToken(ctx context.Context, id string) (*RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, id string) error
}
