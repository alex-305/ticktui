package context

import (
	"github.com/alex-305/ticktui/internal/api"
	"github.com/alex-305/ticktui/internal/config"
)

type AppContext struct {
	APIClient *api.Client
	Config    *config.Config
}
