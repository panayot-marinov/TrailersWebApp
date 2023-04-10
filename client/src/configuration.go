package src

type Configuration struct {
	Protocol   string `yaml:"protocol"`
	ClientPort int    `yaml:"clientPort"`
	ServerHost string `yaml:"serverHost"`
	ServerPort int    `yaml:"serverPort"`
}
