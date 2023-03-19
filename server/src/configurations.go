package src

type DbConfiguration struct {
	Host           string `yaml:"host"`
	Name           string `yaml:"name"`
	User           string `yaml:"user"`
	Pass           string `yaml:"pass"`
	Port           int    `yaml:"port"`
	ConnectTimeout int    `yaml:"connectTimeout"`
	SslMode        string `yaml:"sslMode"`
}

type MailApiConfiguration struct {
	AccessTokenPrivateKeyPath  string `yaml:"accessTokenPrivateKeyPath"`
	AccessTokenPublicKeyPath   string `yaml:"accessTokenPublicKeyPath"`
	RefreshTokenPrivateKeyPath string `yaml:"refreshTokenPrivateKeyPath"`
	RefreshTokenPublicKeyPath  string `yaml:"refreshTokenPublicKeyPath"`
	JwtExpiration              int    `yaml:"jwtExpiration"` // in minutes
	SendGridApiKey             string `yaml:"sendGridApiKey"`
	MailVerifCodeExpiration    int    `yaml:"mailVerifCodeExpiration"` // in hours
	PassResetCodeExpiration    int    `yaml:"passResetCodeExpiration"` // in minutes
	MailVerifTemplateID        string `yaml:"mailVerifTemplateID"`
	PassResetTemplateID        string `yaml:"passResetTemplateID"`
}

// Configurations wraps all the config variables required by the auth service
type Configuration struct {
	ServerAddress string
	DbConfig      DbConfiguration      `yaml:"dbConfiguration"`
	MailApiConfig MailApiConfiguration `yaml:"mailApiConfiguration"`
}
