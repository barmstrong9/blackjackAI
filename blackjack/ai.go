package blackjack

import (
	"fmt"
	"strings"

	"github.com/barmstrong9/deck"
)
//AI sets up the different states for the AI
type AI interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Result(hands [][]deck.Card, dealer []deck.Card)
}
type dealerAI struct{}

func (ai dealerAI) Bet(shuffled bool) int {
	return 1
}

func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Result(hands [][]deck.Card, dealer []deck.Card){
	//nothing
}
//HumanAI returns the humanAI struct
func HumanAI() AI{
	return humanAI{}
}

type humanAI struct{}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled{
		fmt.Println("The deck was just shuffled.")
	}
	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("\nYour Cards:", hand)
		fmt.Println("Dealer's Cards:", dealer)

		fmt.Println("What will you do? (h)it, (s)tand, (d)ouble, s(p)lit")
		var input string
		fmt.Scanf("%s\n", &input)
		lowerInput := strings.ToLower(input)
		switch lowerInput {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("Invalid Option:", lowerInput)
		}
	}
}

func (ai humanAI) Result(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:")
	for _, h := range hands {
		fmt.Println(" ", h)
	}
	fmt.Println("Dealer:", dealer)
}