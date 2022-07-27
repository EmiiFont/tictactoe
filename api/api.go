package api

import (
	json2 "encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

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
	p.Player = strings.ToUpper(p.Player)
	return nil
}

type WinnerResult struct {
	Winner    WinPositions `json:"winner"`
	Board     [3][3]rune   `json:"board"`
	MovesLeft bool         `json:"movesLeft"`
	RoomId    uuid.UUID    `json:"roomId"`
	Turn      int32        `json:"turn"`
}

func ReceiveMove(w http.ResponseWriter, r *http.Request) {
	data := &PlayerMove{}
	if err := render.Bind(r, data); err != nil {
		fmt.Errorf("%f", err)
		return
	}
	boardState := StoreMove(data)
	winner := CheckWinner(boardState)
	movesLeft := MoveLeft(boardState, winner.Winner)
	if !movesLeft {
		DeleteBoard(data.RoomId)
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

func GetBoard(w http.ResponseWriter, r *http.Request) {
	board := [3][3]rune{
		{'_', '_', '_'},
		{'_', '_', '_'},
		{'_', '_', '_'},
	}
	roomId, _ := uuid.NewUUID()
	newBoardRoom := WinnerResult{Board: board, Winner: WinPositions{Winner: '_'}, MovesLeft: true, RoomId: roomId, Turn: player1}
	jsonBoard, _ := json2.Marshal(newBoardRoom)

	w.Write(jsonBoard)
}

func NewGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(uuid.NewString()))
}
