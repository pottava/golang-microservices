package config

// Config defines the application configurations
type Config struct {
	Name           string `trim:"true"`
	Port           uint16
	LogLevel       int
	Mode           string `trim:"true"`
	StaticFileHost string `trim:"true"`
	StaticFilePath string `trim:"true"`
	AccessLog      bool
	SecuredCookie  bool
}
