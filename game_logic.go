package main

var player1 = 'O'
var player2 = 'X'

type WinPositions struct {
	Winner    rune        `json:"winner"`
	Positions [3][2]int32 `json:"positions"`
}

func CheckWinner(board [3][3]int32) WinPositions {
	// check winner for rows
	var winPositions WinPositions
	winPositions.Winner = '_'
	for row := 0; row < 3; row++ {
		if board[row][0] == board[row][1] && board[row][1] == board[row][2] {
			winPositions.Positions[0] = [2]int32{int32(row), 0}
			winPositions.Positions[1] = [2]int32{int32(row), 1}
			winPositions.Positions[2] = [2]int32{int32(row), 2}
			if board[row][0] == player1 {
				winPositions.Winner = player1
				return winPositions
			}
			if board[row][0] == player2 {
				winPositions.Winner = player2
				return winPositions
			}
		}
	}
	// check winner for columns
	for col := 0; col < 3; col++ {
		if board[0][col] == board[1][col] && board[1][col] == board[2][col] {
			winPositions.Positions[0] = [2]int32{int32(0), int32(col)}
			winPositions.Positions[1] = [2]int32{int32(1), int32(col)}
			winPositions.Positions[2] = [2]int32{int32(2), int32(col)}
			if board[0][col] == player1 {
				winPositions.Winner = player1
				return winPositions
			}
			if board[0][col] == player2 {
				winPositions.Winner = player2
				return winPositions
			}
		}
	}
	//check diagonals winner
	if board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		winPositions.Positions[0] = [2]int32{0, 0}
		winPositions.Positions[1] = [2]int32{1, 1}
		winPositions.Positions[2] = [2]int32{2, 2}
		if board[0][0] == player1 {
			winPositions.Winner = player1
			return winPositions
		}
		if board[0][0] == player2 {
			winPositions.Winner = player2
			return winPositions
		}
	}

	//check diagonals winner
	if board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		winPositions.Positions[0] = [2]int32{0, 2}
		winPositions.Positions[1] = [2]int32{1, 1}
		winPositions.Positions[2] = [2]int32{2, 0}
		if board[0][2] == player1 {
			winPositions.Winner = player1
			return winPositions
		}
		if board[0][2] == player2 {
			winPositions.Winner = player1
			return winPositions
		}
	}
	return winPositions
}

func MoveLeft(board [3][3]rune, winner rune) bool {
	if winner != '_' {
		return false
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == '_' {
				return true
			}
		}
	}
	return false
}
