package main

import (
	"fmt"
	"github.com/daviddengcn/go-villa"
	"log"
	"math"
)

type tinRanker struct {
	scores   map[int64]float64
	games    map[int64]int64
	wins     map[int64]int64
	draws    map[int64]int64
	losses   map[int64]int64
	numGames int
	sqe      float64
	battles  int64
}

func (r *tinRanker) Print(names map[int64]string) {
	fmt.Println(math.Sqrt(r.sqe / float64(r.battles)))
	keys := make([]int64, 0, len(r.scores))
	for k := range r.scores {
		keys = append(keys, k)
	}
	villa.SortF(len(keys),
		func(i, j int) bool { return r.scores[keys[i]] > r.scores[keys[j]] },
		func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })

	numPlayers := 0
	for _, k := range keys {
		name, ok := names[k]
		if !ok {
			continue
		}
		numPlayers++
		fmt.Println(numPlayers, ",", name, ",", k, ",", r.scores[k], ",", r.games[k])
	}
}

func (r *tinRanker) SetRating(player int64, rating float64) {
	r.scores[player] = rating
}

func (r *tinRanker) Rating(player int64) float64 {
	rating, ok := r.scores[player]
	if !ok {
		rating = float64(1000)
	}
	return rating
}

func (r *tinRanker) AddGame(players []*Player) {
	if len(players) != 7 {
		log.Fatal("Incorrect number of players")
	}

	r.numGames++

	deltaRatings := make([]float64, 7)

	comparisons := 0
	for i, _ := range players {
		for j, _ := range players {
			if i == j {
				continue
			}
			if !players[i].IsWinner() && !players[j].IsWinner() {
				// players who've drawn as losers aren't considered
				continue
			}
			comparisons++
		}
	}
	for i, _ := range players {
		r.games[players[i].userID]++
		for j, _ := range players {
			if i == j {
				continue
			}
			actualOutcome := 0.0
			if !players[i].IsWinner() && !players[j].IsWinner() {
				continue
			}
			if players[i].IsWinner() && players[j].IsWinner() {
				// if we both won, then we drew
				actualOutcome = 0.5
			} else if players[i].IsWinner() {
				switch players[i].potType {
				case "Winner-takes-all":
					actualOutcome = 1.0
				case "Points-per-supply-center":
					actualOutcome = 1.0
				default:
					log.Fatal("Pot type", players[i].potType, "not implemented")
				}
			}
			kFactor := 40.0 / float64(comparisons)
			expectedOutcome := 1.0 / (1.0 + math.Pow(10, (r.Rating(players[i].userID)-r.Rating(players[j].userID))/400.0))
			change := (kFactor * (actualOutcome - expectedOutcome))
			deltaRatings[i] += change
			r.sqe += math.Pow(actualOutcome-expectedOutcome, 2)
			r.battles++
		}

	}
	for i, rating := range deltaRatings {
		r.SetRating(players[i].userID, r.Rating(players[i].userID)+rating)
	}

}

func (r *tinRanker) NumGames() int {
	return r.numGames
}

func NewTinRanker() Ranker {
	var r tinRanker
	r.scores = make(map[int64]float64)
	r.games = make(map[int64]int64)
	return &r
}
