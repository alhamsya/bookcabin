package response

type Common struct {
	HttpCode int `json:"-"`

	Message string   `json:"message"`
	Data    any      `json:"data"`
	Errors  []string `json:"errors"`
}
