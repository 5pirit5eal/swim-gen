package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/golobby/dotenv"
)

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
		Name         string          `env:"DB_NAME"`
		Instance     string          `env:"DB_INSTANCE"`
		Port         string          `env:"DB_PORT"`
		User         string          `env:"DB_USER"`
		Pass         string          `env:"DB_PASS"`
		PassLocation string          `env:"DB_PASS_LOCATION"`
		Method       models.DBMethod `env:"DB_METHOD"`
	}
}

func LoadConfig(filename string, overwrite bool) (Config, error) {
	var cfg Config
	// Load config from file
	if file, err := os.Open(filename); err == nil {
		if err := dotenv.NewDecoder(file).Decode(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to load config from file: %w", err)
		}
	}
	// Directly load env variables or defaults if not set by file
	if err := set(&cfg, "env", "default", overwrite); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func setField(field reflect.Value, val string) error {
	if !field.CanSet() {
		return fmt.Errorf("can't set value %s", field)
	}

	// Setting values support string and int
	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(val, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(val).Convert(field.Type()))
	case reflect.Float64:
		if val, err := strconv.ParseFloat(val, 64); err == nil {
			field.Set(reflect.ValueOf(val).Convert(field.Type()))
		}
	}

	return nil
}

// set populates the struct fields with values from environment variables or default values.
// It uses the struct field tags to determine which environment variable to use.
func set(ptr any, envTag string, defaultTag string, overwrite bool) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer")
	}

	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		// Skip if value was already set
		if !v.Field(i).IsZero() && !overwrite {
			continue
		}
		defaultVal := t.Field(i).Tag.Get(defaultTag)
		if envVal := t.Field(i).Tag.Get(envTag); envVal != "" {
			val, found := os.LookupEnv(envVal)

			if !found && v.Field(i).IsZero() && defaultVal != "" {
				if err := setField(v.Field(i), defaultVal); err != nil {
					return err
				}
			} else if found {
				if err := setField(v.Field(i), val); err != nil {
					return err
				}
			}
		} else if v.Field(i).IsZero() && defaultVal != "" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}
	return nil
}
