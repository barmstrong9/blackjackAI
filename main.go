package main

import (
	"fmt"

	"github.com/barmstrong9/deck"

	"github.com/barmstrong9/blackjackAI/blackjack"
)

type basicAI struct{
	score int
	seen int
	decks int
}

func (ai *basicAI) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}
	trueScore := ai.score / ((ai.decks * 52 - ai.seen) / 52)
	switch {
	case trueScore >=14:
		return 10000
	case trueScore <=8:
		return 100
	}
	return 100
}

func (ai *basicAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	dScore := blackjack.Score(dealer)
	cardScore := blackjack.Score(hand[0])
	if len(hand) == 2{
		if hand[0] == hand[1]{
		switch {
		case cardScore == 11 || cardScore == 8:
			return blackjack.MoveSplit
		case (cardScore ==2 || cardScore ==3) && (dScore >= 4 && dScore <=7):
			return blackjack.MoveSplit
		case cardScore ==4 && (dScore == 5 || dScore == 6):
			return blackjack.MoveSplit
		case cardScore ==6 && (dScore >= 2 && dScore <= 6):
			return blackjack.MoveSplit
		case cardScore ==7 && (dScore >= 2 && dScore <= 7):
			return blackjack.MoveSplit
		case cardScore ==9 && (dScore >= 2 && dScore <= 9):
			return blackjack.MoveSplit
		}
		}
		switch {
		case (score == 9 && !blackjack.Soft(hand...))&&(dScore >= 3 && dScore <= 6):
			return blackjack.MoveDouble
		case (score == 10 && !blackjack.Soft(hand...))&&(dScore != 10 && dScore != 11):
			return blackjack.MoveDouble
		case (score == 11 && !blackjack.Soft(hand...))&& dScore != 11:
			return blackjack.MoveDouble
		case ((score ==13 || score ==14) && blackjack.Soft(hand...) == true) && (dScore == 5 || dScore ==6):
			return blackjack.MoveDouble
		case ((score == 15 || score ==16) && blackjack.Soft(hand...) == true) && (dScore >= 4 &&  dScore <=6):
			return blackjack.MoveDouble
		case ((score == 17 || score ==18) && blackjack.Soft(hand...) == true) && (dScore >=3 && dScore <=6):
			return blackjack.MoveDouble
		case score <= 11 && !blackjack.Soft(hand...):
			return blackjack.MoveHit
		case score == 12 && !blackjack.Soft(hand...) && (dScore >= 4 && dScore <= 6):
			return blackjack.MoveStand
		case score == 12 && !blackjack.Soft(hand...):
			return blackjack.MoveHit
		case (score >= 13 && score <= 16)&& !blackjack.Soft(hand...) && (dScore >= 2 && dScore <= 6):
			return blackjack.MoveStand
		case (score >= 13 && score <= 16)&& !blackjack.Soft(hand...):
			return blackjack.MoveHit
		case score >= 17 && !blackjack.Soft(hand...):
			return blackjack.MoveStand
		case score <= 17 && blackjack.Soft(hand...):
			return blackjack.MoveHit
		case (score == 18 && blackjack.Soft(hand...)) && (dScore >= 9 && dScore <=11):
			return blackjack.MoveHit
		case score == 18 && blackjack.Soft(hand...):
			return blackjack.MoveStand
		case score >= 19 && blackjack.Soft(hand...):
			return blackjack.MoveStand
		}
		 
	}
	return blackjack.MoveStand
}

func (ai *basicAI) Result(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer{
		ai.count(card)
	}
	for _, hand := range hands{
		for _, card := range hand{
			ai.count(card)
		}
	}
}

func (ai *basicAI) count(card deck.Card){
	score := blackjack.Score(card)
		switch {
		case score >= 10:
			ai.score--
		case score <= 6:
			ai.score++
		}
		ai.seen++
}

func main() {
	opts := blackjack.Options{
		Decks:           4,
		Hands:           90000,
		BlackJackPayout: 1.5,
	}
	game := blackjack.New(opts)
	//use winnings := game.Play(blackjack.HumanAI()) for non-simulation
	//use winnings := game.Play(&basicAI{
	//									decks: 4
	//									}) for non-simulation
	winnings := game.Play(&basicAI{
		decks: 4,
	})
	fmt.Println("winnings:", winnings)
}
