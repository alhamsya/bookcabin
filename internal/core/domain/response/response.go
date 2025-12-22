package modelResponse

type Common struct {
	HttpCode int `json:"-"`

	Data     any             `json:"data"`
	Errors   []string        `json:"errors,omitempty"`
	Metadata *CommonMetadata `json:"metadata,omitempty"`
}

type CommonMetadata struct {
	TotalResult int `json:"total_result"`
}
