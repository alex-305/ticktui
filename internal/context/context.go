package context

import (
	"github.com/alex-305/ticktui/internal/config"
	api "github.com/alex-305/ticktui/pkg/ticktickapi"
)

type AppContext struct {
	APIClient *api.Client
	Config    *config.Config
}
