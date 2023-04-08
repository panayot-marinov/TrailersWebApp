package src

type Configuration struct {
	Protocol   string `yaml:"protocol"`
	ClientPort int    `yaml:"clientPort"`
	ServerPort int    `yaml:"serverPort"`
}
