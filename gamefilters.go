package main

type GameAcceptor interface {
	Accept([]*Player) bool
}

type GameFilter struct {
	filters []GameAcceptor
}

func (gf *GameFilter) Accept(game []*Player) bool {
	for _, filter := range gf.filters {
		if filter != nil {
			accept := filter.Accept(game)
			if !accept {
				return accept
			}
		}
	}

	return true
}

type OnlyClassic struct{}
type OnlyLive struct{}
type NonLive struct{}
type OnlyWTA struct{}
type OnlyFullPress struct{}
type OnlyGunboat struct{}

func (*OnlyClassic) Accept(players []*Player) bool {
	// Also exclude games that finished before A1903
	if len(players) == 7 &&
		players[0].variantID == 1 &&
		players[0].turn >= 5 {
		return true
	}
	return false
}

func (*OnlyLive) Accept(players []*Player) bool {
	if players[0].phaseMinutes < 60 {
		return true
	}
	return false
}

func (*NonLive) Accept(players []*Player) bool {
	if players[0].phaseMinutes >= 60 {
		return true
	}
	return false
}

func (*OnlyWTA) Accept(players []*Player) bool {
	if players[0].potType == "Winner-takes-all" {
		return true
	}
	return false
}
func (*OnlyFullPress) Accept(players []*Player) bool {
	if players[0].pressType == "Regular" {
		return true
	}
	return false
}
