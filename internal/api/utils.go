package api

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

const (
	StatusOK      = 200
	BadRequest    = 400
	InvalidBody   = 401
	InvalidTaskId = 402
	SystemError   = 500
)

const (
	XBJFlag         = "xbj_message"
	UserFlag        = "user_message"
	EmbeddingSearch = "product_search_embed"
)
