package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	BLANK_LINE = "    |   |\n"
	BREAK_LINE = "----+---+----\n"
	ENTER      = "enter"
	UP         = "up"
	LEFT       = "left"
	RIGHT      = "right"
	DOWN       = "down"
	CTRLC      = "ctrl+c"
	QUIT       = "q"
)

func (m model) Init() tea.Cmd {
	return createReceiveMove(*m.conn)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Was it a key press?
	case tea.KeyMsg:

		// Which key was pressed?
		switch msg.String() {
		case CTRLC, QUIT:
			return m, tea.Quit

		case UP:
			if m.selectedRow > 0 {
				m.selectedRow -= 1
			}
		case DOWN:
			if m.selectedRow < BOARD_SIZE-1 {
				m.selectedRow += 1
			}
		case LEFT:
			if m.selectedColumn > 0 {
				m.selectedColumn -= 1
			}
		case RIGHT:
			if m.selectedColumn < BOARD_SIZE-1 {
				m.selectedColumn += 1
			}
		case ENTER:
			m = m.handleEnter()
			// If there is a winner, end the game.
			if m.winner != "" {
				return m, tea.Quit
			}
		}

		// Return the modified model and receiveMove function.
		return m, createReceiveMove(*m.conn)
	case moveMessage:
		commandParts := strings.Split(",", msg.command)
		switch commandParts[0] {
		case ENTER:
			// TODO handle enter again but check if it's his turn!
			// maybe another handle method?
			// TODO parse to get the opponents cursor position
			m.handleOpponentEnter(commandParts[0], commandParts[1], commandParts[2], commandParts[3])

		default:
			log.Fatalf("Unknown key %q", commandParts[0])
		}

		return m, nil
	}

	// TODO see if this is necessary and if it can be removed.
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

func main() {
	waitForPlayer := flag.Bool("wait", false, "Wait for player to join you.")
	ipAddress := flag.String("ip", "", "IPv4 address to connect to the other player.")
	port := flag.String("port", "8080", "Port on which to listen to.")

	flag.Parse()

	// Set logs to tictactoe.log.
	f, err := tea.LogToFile("tictactoe.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// Declare a connection and a player string.
	var conn net.Conn
	var player string
	// Wait for a connection from the opponent.
	if *waitForPlayer {
		ln, err := net.Listen("tcp", ":"+*port)
		if err != nil {
			log.Fatalf("Failed to listen on port %v: %v", port, err)
		}
		log.Printf("Listening on port %v\n", *port)

		conn, err = ln.Accept()
		if err != nil {
			log.Fatalf("Failed to accept a connection: %v", err)
		}
		log.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())

		// Read which player the opponent chooses to be.
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			log.Fatalf("Could not read from the connection: %v", err)
		}

		player := string(buffer[:length])
		log.Printf("The opponent chose to be %s\n", player)

	} else {
		// Dial a connection to the waiting opponent.
		var err error
		conn, err = net.Dial("tcp", *ipAddress+":"+*port)
		if err != nil {
			log.Fatalf("Could not connect to %s:%s: %v", *ipAddress, *port, err)
		}

		// Choose between X and O randomly.
		src := rand.NewSource(time.Now().UnixNano())
		r := rand.New(src)
		if r.Int()%2 == 0 {
			player = PLAYER_X
		} else {
			player = PLAYER_O
		}
		log.Printf("Randomly chose to be %s\n", player)

		// Send the choice to the opponent.
		_, err = conn.Write([]byte(player))
		if err != nil {
			log.Fatalf("could not send player choice to waiting opponent: %v", err)
		}
		log.Println("Successfully sent the choice to waiting opponent.")
	}

	// Start the program.
	p := tea.NewProgram(NewModel(&conn, "X"))
	if _, err := p.Run(); err != nil {
		log.Printf("Alas, there's been an error: %v\n", err)
	}

	conn.Close()
}
