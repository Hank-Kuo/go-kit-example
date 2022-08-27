package service

import (
	"context"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type priceSvc struct {
	logger log.Logger
}

func NewService(logger log.Logger, ints, chars metrics.Counter) price.PriceService {
	var svc price.PriceService
	{
		svc = &priceSvc{logger: logger}
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}

func (svc *priceSvc) Sum(_ context.Context, price int64, fee int64) (res int64, err error) {
	return price + fee, nil
}

func (svc *priceSvc) Exchange(ctx context.Context, cost int64, currency string) (res int64, err error) {
	return cost, err
}
