package main

import (
	"log"
	"net"
)

const (
	EMPTY      = " "
	BOARD_SIZE = 3
	PLAYER_X   = "x"
	PLAYER_O   = "o"
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

func (m model) handleEnter() model {
	// If it is not this player's turn, no action.
	if m.playerTurn == m.player {
		return m
	}

	// If the cell is not empty, no action.
	if m.board[m.selectedRow][m.selectedColumn] != EMPTY {
		return m
	}

	// Mark cell.
	m.board[m.selectedRow][m.selectedColumn] = m.player

	// Check if the player won.
	if m.isWinner() {
		m.winner = m.playerTurn
	}

	// TODO send message to opponent.
	err := m.sendMove()
	if err != nil {
		log.Fatalf("could not send move to opponent: %v", err)
	}

	// Switch players turn.
	if m.playerTurn == PLAYER_X {
		m.playerTurn = PLAYER_O
	} else {
		m.playerTurn = PLAYER_X
	}

	return m
}

// isWinner checks if the current player is the winner.
func (m *model) isWinner() bool {
	// Check every row.
	for i := 0; i < BOARD_SIZE; i++ {
		won := true
		for j := 0; j < BOARD_SIZE; j++ {
			won = won && m.board[i][j] == m.playerTurn
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
			won = won && m.board[i][j] == m.playerTurn
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
		won1 = won1 && m.board[i][i] == m.playerTurn
		won2 = won2 && m.board[i][BOARD_SIZE-i-1] == m.playerTurn
		if !won1 && !won2 {
			break
		}
	}

	return won1 || won2
}

func (m *model) sendMove() error {
	// TODO: send the move over the TCP connection.
	return nil
}
