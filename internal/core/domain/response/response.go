package modelResponse

type Common struct {
	HttpCode int `json:"-"`

	Data   any      `json:"data"`
	Errors []string `json:"errors"`
}
