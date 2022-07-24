package api

import (
	json2 "encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
)

func StoreMove(playerMove *PlayerMove) [3][3]rune {
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

func DeleteBoard(key string) {
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
