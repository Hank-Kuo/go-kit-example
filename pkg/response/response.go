package response

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Responser interface {
	Response() (res interface{})
}

type BinaryResponser interface {
	GetBinary() (res []byte)
}

type Paging struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}
