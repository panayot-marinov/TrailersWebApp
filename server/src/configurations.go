package src

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
	ServerAddress              string
	DBHost                     string
	DBName                     string
	DBUser                     string
	DBPass                     string
	DBPort                     string
	DBConn                     string
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int // in minutes
	SendGridApiKey             string
	MailVerifCodeExpiration    int // in hours
	PassResetCodeExpiration    int // in minutes
	MailVerifTemplateID        string
	PassResetTemplateID        string
}
