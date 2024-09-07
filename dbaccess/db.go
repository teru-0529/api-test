package dbaccess

import (
	"context"
	"net/http"

	"github.com/sethvargo/go-envconfig"
)

// STRUCT:
type DbAccesser struct {
	client        *http.Client
	PostgrestHost string `env:"POSTGREST_API_HOST"`
	ReseterHost   string `env:"DB_RESETER_API_HOST"`
}

func New() *DbAccesser {
	var accessor DbAccesser
	envconfig.Process(context.Background(), &accessor)

	accessor.client = &http.Client{}
	return &accessor
}
