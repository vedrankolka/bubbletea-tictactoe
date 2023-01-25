package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	EMPTY      = " "
	BOARD_SIZE = 3
	PLAYER_X   = "X"
	PLAYER_O   = "O"
	TIE        = "tie"
)

type model struct {
	// board contains the game state.
	board [][]string

	// selectedRow is the currently row of the cursor.
	selectedRow int

	// selectedColumn is the currently column of the cursor.
	selectedColumn int

	// winner notes the winner of the game.
	winner string

	// conn is the TCP connection to the opponent.
	conn *net.Conn

	// player notes if the player is X or O.
	player string

	// playerTurn notes which player's turn it is.
	playerTurn string
}

func NewModel(conn *net.Conn, player string) model {
	board := make([][]string, BOARD_SIZE)
	for i := 0; i < BOARD_SIZE; i++ {
		board[i] = make([]string, BOARD_SIZE)
		for j := 0; j < BOARD_SIZE; j++ {
			// Set an empty character for each cell on the board.
			board[i][j] = EMPTY
		}
	}

	return model{
		board:          board,
		selectedRow:    0,
		selectedColumn: 0,
		winner:         "",
		conn:           conn,
		player:         player,
		playerTurn:     PLAYER_X,
	}
}

func (m model) HandleMyEnter() (model, error) {
	m = m.handlePlayerEnter(m.player, m.selectedRow, m.selectedColumn)

	// Send move message to opponent.
	err := m.sendMove(ENTER)
	return m, err
}

func (m model) HandleOpponentEnter(key, opponent, selectedRowStr, selectedColStr string) (model, error) {
	selectedRow, err := strconv.Atoi(selectedRowStr)
	if err != nil {
		return m, err
	}

	selectedCol, err := strconv.Atoi(selectedColStr)
	if err != nil {
		return m, err
	}

	return m.handlePlayerEnter(opponent, selectedRow, selectedCol), nil
}

// handlePlayerEnter handles an Enter from the given player.
func (m model) handlePlayerEnter(player string, row int, col int) model {
	log.Printf("Handling move of player %s on turn %s at cursor [%d, %d]\n",
		player, m.playerTurn, row, col,
	)

	// If it's not the player's turn, no action.
	if player != m.playerTurn {
		log.Printf("Ignoring %s's move as it's %s's turn.\n", player, m.playerTurn)
		return m
	}

	// If the cell is not empty, no action.
	if m.board[row][col] != EMPTY {
		log.Printf("Ignoring move as cell [%d, %d] is not empty.\n", row, col)
		return m
	}

	// Mark cell.
	m.board[row][col] = player
	log.Printf("%s marked cell [%d, %d]\n", player, row, col)

	// Check if the player won.
	if m.isWinner(player) {
		m.winner = player
		log.Printf("Marked that %s won.\n", player)
	}

	// Check if it's a tie.
	if m.isTie() {
		m.winner = TIE
		log.Printf("Marked the game as a tie.\n")
	}

	// Switch player's turn.
	if m.playerTurn == PLAYER_X {
		m.playerTurn = PLAYER_O
	} else {
		m.playerTurn = PLAYER_X
	}

	return m
}

// isWinner checks if the given player is the winner.
func (m *model) isWinner(player string) bool {
	// Check every row.
	for i := 0; i < BOARD_SIZE; i++ {
		won := true
		for j := 0; j < BOARD_SIZE; j++ {
			won = won && m.board[i][j] == player
			if !won {
				break
			}
		}

		if won {
			return true
		}
	}

	// Check every column.
	for j := 0; j < BOARD_SIZE; j++ {
		won := true
		for i := 0; i < BOARD_SIZE; i++ {
			won = won && m.board[i][j] == player
			if !won {
				break
			}
		}

		if won {
			return true
		}
	}

	// Check diagonals.
	won1 := true
	won2 := true
	for i := 0; i < BOARD_SIZE; i++ {
		won1 = won1 && m.board[i][i] == player
		won2 = won2 && m.board[i][BOARD_SIZE-i-1] == player
		if !won1 && !won2 {
			break
		}
	}

	return won1 || won2
}

// isTie checks if the board results in a tie.
// TODO actually check that, now it checks if the board is filled
// and there is no winner.
func (m *model) isTie() bool {
	boardFull := true
	for i := 0; i < BOARD_SIZE; i++ {
		for j := 0; j < BOARD_SIZE; j++ {
			boardFull = boardFull && (m.board[i][j] != EMPTY)
			if !boardFull {
				break
			}
		}
	}

	// Claim it's a tie if the board is full and there is no winner yet.
	return boardFull && m.winner == ""
}

func (m *model) sendMove(key string) error {
	move := fmt.Sprintf("%s,%s,%d,%d", key, m.player, m.selectedRow, m.selectedColumn)
	_, err := (*m.conn).Write([]byte(move))
	return err
}
