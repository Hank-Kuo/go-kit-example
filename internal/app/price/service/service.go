package service

import (
	"context"

	priceUtils "github.com/Hank-Kuo/go-kit-example/internal/app/price/utils"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Service interface {
	Sum(ctx context.Context, price, fee int64) (int64, error)
	Exchange(ctx context.Context, cost int64, currency string) (int64, error)
}

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
)

type priceSvc struct {
	logger log.Logger
}

func NewService(logger log.Logger, ints, chars metrics.Counter) Service {
	var svc Service
	{
		svc = &priceSvc{logger: logger}
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}

func (svc *priceSvc) Sum(_ context.Context, price int64, fee int64) (res int64, err error) {
	if price <= 0 || fee <= 0 {
		return 0, priceUtils.ErrLessZeroe
	}
	if (fee > 0 && price > (intMax-fee)) || (fee < 0 && price < (intMin-fee)) {
		return 0, priceUtils.ErrIntOverflow
	}
	return price + fee, nil
}

func (svc *priceSvc) Exchange(ctx context.Context, cost int64, currency string) (res int64, err error) {
	return cost, err
}
