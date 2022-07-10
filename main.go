package main

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"strings"
)

var player1 = 'O'
var player2 = 'X'

type WinPositions struct {
	Winner    rune        `json:"winner"`
	Positions [3][2]int32 `json:"positions"`
}

func checkWinner(board [3][3]int32) WinPositions {
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

func moveLeft(board [3][3]rune, winner rune) bool {
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

type PlayerMove struct {
	Player string `json:"player"`
	Column int32  `json:"column"`
	Row    int32  `json:"row"`
	RoomId string `json:"roomId"`
}

func (p PlayerMove) Bind(r *http.Request) error {
	if p.Player == "" {
		return errors.New("missing required Player field")
	}
	if p.Row < 0 && p.Row >= 3 {
		return errors.New("row is out of the board, cannot find cell")
	}

	if p.Column < 0 && p.Column >= 3 {
		return errors.New("column is out of the board, cannot find cell")
	}
	p.Player = strings.ToUpper(p.Player) // unset the protected ID
	return nil
}

type WinnerResult struct {
	Winner    WinPositions `json:"winner"`
	Board     [3][3]rune   `json:"board"`
	MovesLeft bool         `json:"movesLeft"`
	RoomId    uuid.UUID    `json:"roomId"`
	Turn      int32        `json:"turn"`
}

func receiveMove(w http.ResponseWriter, r *http.Request) {
	data := &PlayerMove{}
	if err := render.Bind(r, data); err != nil {
		fmt.Errorf("%f", err)
		return
	}
	fmt.Println(data)
	boardState := storeMove(data)
	winner := checkWinner(boardState)
	movesLeft := moveLeft(boardState, winner.Winner)
	if !movesLeft {
		deleteBoard(data.RoomId)
	}
	uud, _ := uuid.Parse(data.RoomId)
	var turn rune
	if data.Player == string(player1) {
		turn = player2
	} else {
		turn = player1
	}
	winnerResult := WinnerResult{Winner: winner, Board: boardState, MovesLeft: movesLeft, RoomId: uud, Turn: turn}
	returnedWinner, _ := json2.Marshal(winnerResult)

	w.Header().Set("Content-Type", "application/json")
	w.Write(returnedWinner)
}

func deleteBoard(key string) {
	db, err := bolt.Open("boards.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("boards"))
		if err != nil {
			return nil
		}
		return b.Delete([]byte(key))
	})
}

func storeMove(playerMove *PlayerMove) [3][3]rune {
	db, err := bolt.Open("boards.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var boardToSave [3][3]rune
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("boards"))
		if err != nil {
			return nil
		}
		boardState := b.Get([]byte(playerMove.RoomId))
		if boardState == nil {
			boardToSave = [3][3]rune{
				{'_', '_', '_'},
				{'_', '_', '_'},
				{'_', '_', '_'},
			}
		} else {
			marErr := json2.Unmarshal(boardState, &boardToSave)
			if marErr != nil {
				return nil
			}
		}
		var symbol rune
		if playerMove.Player == "O" {
			symbol = player1
		} else {
			symbol = player2
		}
		boardToSave[playerMove.Column][playerMove.Row] = symbol
		fmt.Println(boardToSave)
		jsonData, _ := json2.Marshal(boardToSave)

		return b.Put([]byte(playerMove.RoomId), jsonData)
	})

	if err != nil {
		fmt.Println(err)
	}
	return boardToSave
}

func getBoard(w http.ResponseWriter, r *http.Request) {
	board := [3][3]rune{
		{'_', '_', '_'},
		{'_', '_', '_'},
		{'_', '_', '_'},
	}
	roomId, _ := uuid.NewUUID()
	newBoardRoom := WinnerResult{Board: board, Winner: WinPositions{Winner: '_'}, MovesLeft: true, RoomId: roomId, Turn: player1}
	jsonBoard, _ := json2.Marshal(newBoardRoom)
	fmt.Println(roomId)
	w.Write(jsonBoard)
}

func main() {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))
	r.Use(middleware.Logger)
	r.Get("/", getBoard)
	r.Post("/move", receiveMove)
	http.ListenAndServe(":8080", r)
}
