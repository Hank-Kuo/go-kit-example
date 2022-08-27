package endpoints

type SumRequest struct {
	Price int64 `json:"price"`
	Fee   int64 `json:"fee"`
}

type ExchangeRequest struct {
	Cost     int64  `json:"cost"`
	Currency string `json:"currency"`
}

func (r SumRequest) validate() error {
	return nil
}
