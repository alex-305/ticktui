package screens

type ChangeScreenMsg struct {
	NewScreen Screen
}

type ChangeScreenMsgNoHistory struct {
	NewScreen Screen
}

type GoBackScreenMsg struct{}
