package app

func navigateBackAScreen(m *Model) {
	if len(m.history) > 0 {
		lastIndex := len(m.history) - 1
		lastPage := m.history[lastIndex]
		m.history = m.history[:lastIndex]

		m.current = lastPage
	}
}
