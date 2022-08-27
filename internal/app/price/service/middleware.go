package service

import (
	"context"

	"github.com/Hank-Kuo/go-kit-example/internal/app/price"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

type Middleware func(price.PriceService) price.PriceService

type loggingMiddleware struct {
	logger log.Logger
	next   price.PriceService
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next price.PriceService) price.PriceService {
		return loggingMiddleware{logger, next}
	}
}

func (mw loggingMiddleware) Sum(ctx context.Context, price, fee int64) (v int64, err error) {
	defer func() {
		mw.logger.Log("method", "Sum", "price", price, "fee", fee, "v", v, "err", err)
	}()
	return mw.next.Sum(ctx, price, fee)
}

func (mw loggingMiddleware) Exchange(ctx context.Context, a int64, b string) (v int64, err error) {
	defer func() {
		mw.logger.Log("method", "Concat", "a", a, "b", b, "v", v, "err", err)
	}()
	return mw.next.Exchange(ctx, a, b)
}

type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  price.PriceService
}

func InstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
	return func(next price.PriceService) price.PriceService {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

func (mw instrumentingMiddleware) Sum(ctx context.Context, price, fee int64) (int64, error) {
	v, err := mw.next.Sum(ctx, price, fee)
	mw.ints.Add(float64(v))
	return v, err
}

func (mw instrumentingMiddleware) Exchange(ctx context.Context, cost int64, currency string) (int64, error) {
	v, err := mw.next.Exchange(ctx, cost, currency)
	mw.ints.Add(float64(v))
	return v, err
}
