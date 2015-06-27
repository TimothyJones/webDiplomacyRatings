package main

import (
	"fmt"
)

type gameStats struct {
	numGames int
}

func (r *gameStats) Print(names map[int64]string) {
}

func (r *gameStats) SetRating(player int64, rating float64) {
}

func (r *gameStats) Rating(player int64) float64 {
	return 0.0
}

func (r *gameStats) AddGame(players []*Player) {
	r.numGames++

	numWinners := 0
	winner := int64(0)
	for _, player := range players {
		if player.IsWinner() {
			numWinners++
		}
		if player.status == "Won" {
			winner = player.userID
		}
	}

	fmt.Println(players[0].gameID, ",", numWinners, ",", winner, ",", players[0].processTime)
}

func (r *gameStats) NumGames() int {
	return r.numGames
}

func NewStatsRanker() Ranker {
	var r gameStats
	return &r
}
