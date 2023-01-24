package main

import (
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

type moveMessage struct{ command string }

type errMsg struct{ err error }

func (e *errMsg) Error() string {
	return e.err.Error()
}

func createReceiveMove(conn net.Conn) func() tea.Msg {
	return func() tea.Msg {
		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)
		if err != nil {
			return errMsg{err: err}
		}

		return moveMessage{command: string(buffer[:len])}
	}
}

// func receiveMove() tea.Msg {
// 	buffer := make([]byte, 1024)
// 	len, err := conn.Read(buffer)
// 	if err != nil {
// 		return errMsg{err: err}
// 	}

// 	return moveMessage{command: string(buffer[:len])}
// }
