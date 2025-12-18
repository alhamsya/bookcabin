package xhtpp

type ThirdParty struct {
	Service    string `json:"service"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
	Request    string `json:"request"`
	Response   string `json:"response"`
}
