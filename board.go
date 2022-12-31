package main

type model struct {
	board          [][]string
	selectedRow    int
	selectedColumn int
	xTurn          bool
	winner         string
}

func (m model) handleEnter() model {
	// If the cell is not empty, no action.
	if m.board[m.selectedRow][m.selectedColumn] != " " {
		return m
	}

	// Mark cell according to the players turn.
	if m.xTurn {
		m.board[m.selectedRow][m.selectedColumn] = "x"
	} else {
		m.board[m.selectedRow][m.selectedColumn] = "o"
	}

	// Check if the player won.
	if isWinner(m) {
		if m.xTurn {
			m.winner = "X"
		} else {
			m.winner = "O"
		}
	}

	// Switch players turn.
	m.xTurn = !m.xTurn

	return m
}

func isWinner(m model) bool {
	// Check if the current player is the winner.
	var player string
	if m.xTurn {
		player = "x"
	} else {
		player = "o"
	}

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
