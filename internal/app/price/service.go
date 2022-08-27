package price

import "context"

type PriceService interface {
	Sum(ctx context.Context, price, fee int64) (int64, error)
	Exchange(ctx context.Context, cost int64, currency string) (int64, error)
}

type Service interface {
	Sum(ctx context.Context, price, fee int64) (int64, error)
	Exchange(ctx context.Context, cost int64, currency string) (int64, error)
}
