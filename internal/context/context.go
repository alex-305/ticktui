package context

import (
	"github.com/alex-305/ticktui/internal/api"
)

type AppContext struct {
	APIClient *api.Client
}
