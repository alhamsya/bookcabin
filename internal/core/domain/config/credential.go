package config

type Credential struct {
	Redis CredentialRedis `mapstructure:"redis"`
}

type CredentialRedis struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
