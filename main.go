package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	//How shall we rank players?
	ranking := NewTinRanker()

	// which games shall we accept?
	filter := &GameFilter{[]GameAcceptor{
		&OnlyClassic{},
		&NonLive{},
		&OnlyWTA{},
		&OnlyFullPress{},
	}}

	answer := scan(ranking, filter)
	fmt.Println(answer, "classic games")
}

// The Player structure represents an entry in the Member table
type Player struct {
	variantID, gameID, userID, pot, supplyCenterNo, phaseMinutes, turn, processTime, IsBanned int64
	gameOver, status, potType, pressType                                                      string
}

func (p *Player) IsWinner() bool {
	if p.status == "Won" || p.status == "Drawn" {
		return true
	}
	if p.status == "Survived" && p.potType == "Points-per-supply-center" {
		return true
	}

	return false
}

type Ranker interface {
	AddGame([]*Player)
	NumGames() int
	Rating(int64) float64
	SetRating(int64, float64)
	Print(map[int64]string)
}

//"variantID","gameID","userID","pot","gameOver","status","supplyCenterNo","potType","phaseMinutes","turn","processTime","pressType","IsBanned"
//"1","3","11","231","Won","Defeated","0","Points-per-supply-center","1440","28","1155585385","Regular","0"
func (p *Player) String() string {
	return fmt.Sprintf("\"%d\",\"%d\",\"%d\",\"%d\",\"%s\",\"%s\",\"%d\",\"%s\",\"%d\",\"%d\",\"%d\",\"%s\",\"%d\"",
		p.variantID,
		p.gameID,
		p.userID,
		p.pot,
		p.gameOver,
		p.status,
		p.supplyCenterNo,
		p.potType,
		p.phaseMinutes,
		p.turn,
		p.processTime,
		p.pressType,
		p.IsBanned,
	)
}

var (
	lastProcess int64
)

func parsePlayer(line []string) (*Player, error) {
	variantID, err := strconv.ParseInt(line[0], 10, 64)
	if err != nil {
		return nil, err
	}
	gameID, err := strconv.ParseInt(line[1], 10, 64)
	if err != nil {
		return nil, err
	}
	userID, err := strconv.ParseInt(line[2], 10, 64)
	if err != nil {
		return nil, err
	}
	pot, err := strconv.ParseInt(line[3], 10, 64)
	if err != nil {
		return nil, err
	}
	supplyCenterNo, err := strconv.ParseInt(line[6], 10, 64)
	if err != nil {
		return nil, err
	}
	phaseMinutes, err := strconv.ParseInt(line[8], 10, 64)
	if err != nil {
		return nil, err
	}
	turn, err := strconv.ParseInt(line[9], 10, 64)
	if err != nil {
		return nil, err
	}
	processTime, err := strconv.ParseInt(line[10], 10, 64)
	if err != nil {
		// some entries don't have processTime, so we assume they processed
		// at whatever time the last game did
		processTime = lastProcess
	}
	lastProcess = processTime
	IsBanned, err := strconv.ParseInt(line[12], 10, 64)
	if err != nil {
		return nil, err
	}

	return &Player{
		variantID:      variantID,
		gameID:         gameID,
		userID:         userID,
		pot:            pot,
		gameOver:       line[4],
		status:         line[5],
		supplyCenterNo: supplyCenterNo,
		potType:        line[7],
		phaseMinutes:   phaseMinutes,
		turn:           turn,
		processTime:    processTime,
		pressType:      line[11],
		IsBanned:       IsBanned,
	}, nil
}

func scan(ranking Ranker, filter *GameFilter) int {
	file, err := os.Open("ghostRatingData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// discard header
	scanner.Scan()
	last := ""
	numGames := 0

	var gamePlayers []*Player
	names := make(map[int64]string)
	mode := "games"

	for scanner.Scan() {
		line := scanner.Text()
		A := strings.Split(strings.Replace(line, "\"", "", -1), ",")
		switch mode {
		case "games":
			if A[0] == "1" {
				// we are only interested in classic games
				if A[1] != last {
					// this is a new game
					// first process any old games
					if gamePlayers != nil && filter.Accept(gamePlayers) {
						ranking.AddGame(gamePlayers)
					}
					numGames++
					gamePlayers = make([]*Player, 0, 7)
					last = A[1]
					player, err := parsePlayer(A)
					if err != nil {
						log.Fatal(err, A)
					}

					gamePlayers = append(gamePlayers, player)

				} else {
					player, err := parsePlayer(A)
					if err != nil {
						log.Fatal(err, A)
					}

					gamePlayers = append(gamePlayers, player)

				}
			}
			if A[0] == "id" {
				// we're done with the games
				mode = "players"
			}
		case "players":
			if A[2] == "0" {
				// this user was not banned so we save their name
				uid, err := strconv.ParseInt(A[0], 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				names[uid] = A[1]
			}
		}
	}
	ranking.Print(names)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return ranking.NumGames()
}
