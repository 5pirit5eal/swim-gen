package models

type Config struct {
	ProjectID string `env:"PROJECT_ID"`
	Region    string `env:"REGION"`
	Model     string `env:"MODEL"`
	APIKey    string `env:"API_KEY"`
	LogLevel  string `env:"LOG_LEVEL"`
	Port      string `env:"PORT"`

	Embedding struct {
		Name  string `env:"EMBEDDING_NAME"`
		Model string `env:"EMBEDDING_MODEL"`
		Size  int    `env:"EMBEDDING_SIZE"`
	}

	DB struct {
		Name         string   `env:"DB_NAME"`
		Instance     string   `env:"DB_INSTANCE"`
		Port         string   `env:"DB_PORT"`
		User         string   `env:"DB_USER"`
		Pass         string   `env:"DB_PASS"`
		PassLocation string   `env:"DB_PASS_LOCATION"`
		Method       DBMethod `env:"DB_METHOD"`
	}
}
