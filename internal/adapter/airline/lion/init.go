package lion

import "github.com/alhamsya/bookcabin/internal/core/domain/config"

type Airline struct {
	CfgUpstream *config.Upstream
}

func New(param *Airline) *Airline {
	return &Airline{
		CfgUpstream: param.CfgUpstream,
	}
}
