package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	BOARD_SIZE = 3
	BLANK_LINE = "    |   |\n"
	BREAK_LINE = "----+---+----\n"
)

func initialModel() model {
	board := make([][]string, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		board[i] = make([]string, BOARD_SIZE)
		for j := 0; j < BOARD_SIZE; j++ {
			// Set an empty character for each cell on the board.
			board[i][j] = " "
		}
	}

	return model{
		board:          board,
		selectedRow:    0,
		selectedColumn: 0,
		xTurn:          true,
		winner:         "",
	}
}

func (m model) Init() tea.Cmd {
	return receiveMove
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Was it a key press?
	case tea.KeyMsg:

		// Which key was pressed?
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.selectedRow > 0 {
				m.selectedRow -= 1
			}
		case "down":
			if m.selectedRow < BOARD_SIZE-1 {
				m.selectedRow += 1
			}
		case "left":
			if m.selectedColumn > 0 {
				m.selectedColumn -= 1
			}
		case "right":
			if m.selectedColumn < BOARD_SIZE-1 {
				m.selectedColumn += 1
			}
		case "enter":
			m = m.handleEnter()
			// If there is a winner, end the game.
			if m.winner != "" {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n"

	for i := 0; i < BOARD_SIZE; i++ {
		s += BLANK_LINE + " "
		for j := 0; j < BOARD_SIZE; j++ {
			// If the cell is selected, show ">".
			if i == m.selectedRow && j == m.selectedColumn {
				s += ">"
			} else {
				s += " "
			}
			// Then show the character in the cell.
			s += m.board[i][j]

			if j != BOARD_SIZE-1 {
				s += " |"
			} else {
				s += "\n"
			}
		}
		s += BLANK_LINE

		if i != BOARD_SIZE-1 {
			s += BREAK_LINE
		} else {
			s += "\n"
		}
	}

	// If there is a winner, add a line of text.
	if m.winner != "" {
		s += m.winner + " wins!\n"
	}

	return s
}

// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Alas, there's been an error: %v", err)
// 		os.Exit(1)
// 	}
// }

// var (
// 	waitForPlayer = flag.Bool("wait", false, "Wait for player to join you.")
// 	ipAddress     = flag.String("ip", "", "IPv4 address to connect to the other player.")
// )

func main() {
	waitForPlayer := flag.Bool("wait", false, "Wait for player to join you.")
	ipAddress := flag.String("ip", "", "IPv4 address to connect to the other player.")
	port := flag.String("port", "8080", "Port on which to listen to.")

	flag.Parse()

	if *waitForPlayer {
		ln, err := net.Listen("tcp", ":"+*port)
		if err != nil {
			log.Fatalf("Failed to listen on port %v: %v", port, err)
		}

		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Failed to accept a connection: %v", err)
		}

		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("Could not read from the connection: %v", err)
		}
		fmt.Printf("Command received: %s\n", string(buffer[:len]))

		conn.Write([]byte("Aye aye Captain!"))
		conn.Close()
	} else {
		conn, err := net.Dial("tcp", *ipAddress+":"+*port)
		if err != nil {
			log.Fatalf("Could not connect to %s:%s: %v", *ipAddress, *port, err)
		}

		_, err = conn.Write([]byte("We attack at dawn!"))
		if err != nil {
			log.Fatalf("Could not send message: %v", err)
		}

		buffer := make([]byte, 1024)
		len, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("Could not read from connection: %v", err)
		}

		fmt.Printf("Response received: %s\n", string(buffer[:len]))
	}
}
