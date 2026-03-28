package screens

type ChangeScreenMsg struct {
	NewScreen Screen
}

type GoBackScreenMsg struct{}

type TaskDeletedMsg struct{ err error }
