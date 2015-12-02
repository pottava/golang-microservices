package config

// Config defines the application configurations
type Config struct {
	Name            string `trim:"true"`
	Host            string `trim:"true"`
	Port            uint16
	LogLevel        int
	AccessLog       bool
	CorsMethods     string `trim:"true"`
	CorsOrigin      string `trim:"true"`
	SecuredCookie   bool
	TwitterKey      string `trim:"true"`
	TwitterSecret   string `trim:"true"`
	TwitterCallback string `trim:"true"`
}
