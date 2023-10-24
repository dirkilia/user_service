package types

type Config struct {
	Server struct {
		Port string `envconfig:"SERVER_PORT"`
		Host string `envconfig:"SERVER_HOST"`
	}
	Database struct {
		Host     string `envconfig:"DB_HOST"`
		Port     string `envconfig:"DB_PORT"`
		Username string `envconfig:"DB_USER"`
		Password string `envconfig:"DB_PASSWORD"`
		Name     string `envconfig:"DB_NAME"`
		Driver   string `envconfig:"DB_DRIVER"`
	}
}
