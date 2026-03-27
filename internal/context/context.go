package context

import "github.com/alex-305/ticktui/internal/services"

type AppContext struct {
	TaskService *services.TaskService
}
