package main

import (
	"fmt"
	"github.com/daviddengcn/go-villa"
	"log"
	"math"
)

type ghostRanker struct {
	scores   map[int64]float64
	games    map[int64]int64
	numGames int
	sqe      float64
}

func (r *ghostRanker) Print(names map[int64]string) {
	fmt.Println(math.Sqrt(r.sqe / float64(r.numGames)))
	keys := make([]int64, 0, len(r.scores))
	for k := range r.scores {
		keys = append(keys, k)
	}
	villa.SortF(len(keys),
		func(i, j int) bool { return r.scores[keys[i]] > r.scores[keys[j]] },
		func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	for i, k := range keys {
		name, ok := names[k]
		if !ok {
			continue
		}
		fmt.Println(i, ",", name, ",", k, ",", r.scores[k], ",", r.games[k])
	}
}

func (r *ghostRanker) SetRating(player int64, rating float64) {
	r.scores[player] = rating
}

func (r *ghostRanker) Rating(player int64) float64 {
	rating, ok := r.scores[player]
	if !ok {
		rating = float64(100)
	}
	return rating
}

func (r *ghostRanker) AddGame(players []*Player) {
	if len(players) != 7 {
		log.Fatal("Incorrect number of players")
	}

	r.numGames++
	winSize := 0
	sumRatings := float64(0)
	for _, player := range players {
		if player.IsWinner() {
			winSize++
		}
		r.games[player.userID]++
		sumRatings += r.Rating(player.userID)
	}
	if winSize == 0 {
		log.Fatal("Nobody won this game!")
	}
	for i, _ := range players {
		actualOutcome := float64(0)
		if players[i].IsWinner() {
			switch players[i].potType {
			case "Winner-takes-all":
				actualOutcome = float64(1) / float64(winSize)
			default:
				log.Fatal("Pot type", players[i].potType, "not implemented")
			}
		}

		kFactor := 40.0 //sumRatings / 17.5
		expectedOutcome := r.Rating(players[i].userID) / sumRatings

		newRating := r.Rating(players[i].userID) + (kFactor * (actualOutcome - expectedOutcome))
		r.SetRating(players[i].userID, newRating)

		r.sqe += math.Pow(actualOutcome-expectedOutcome, 2)
	}
}

func (r *ghostRanker) NumGames() int {
	return r.numGames
}

func NewGhostRanker() Ranker {
	var r ghostRanker
	r.scores = make(map[int64]float64)
	r.games = make(map[int64]int64)
	return &r
}
