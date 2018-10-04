package blackjack

import (
	"errors"

	"github.com/barmstrong9/deck"
)

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type state int8

//Options sets up the different parts of blackjack
type Options struct {
	Decks           int
	Hands           int
	BlackJackPayout float64
}

//New sets up default settings for blackjack
func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackJackPayout == 0.0 {
		opts.BlackJackPayout = 1.5
	}
	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackJackPayout
	return g
}

//Game contains all of the different data types needed for blackjack
type Game struct {
	nDecks          int
	nHands          int
	blackjackPayout float64

	state state
	deck  []deck.Card

	player    []hand
	handIdx   int
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handIdx].cards
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("it currently isn't any players turn")
	}
}

type hand struct {
	cards []deck.Card
	bet   int
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	if bet < 100 {
		panic("bet must be at least 100")
	}
	g.playerBet = bet
}

func deal(g *Game) {
	playerHand := make([]deck.Card, 0, 5)
	g.handIdx = 0
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		playerHand = append(playerHand, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.player = []hand{
		{
			cards: playerHand,
			bet:   g.playerBet,
		},
	}
	g.state = statePlayerTurn
}

//Play is the main part of the game
func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3
	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = deck.NewDeck(deck.MultiDeck(g.nDecks), deck.Shuffle)
			shuffled = true
		}

		bet(g, ai, shuffled)
		deal(g)
		if Blackjack(g.dealer...) {
			endRound(g, ai)
			continue
		}

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				//noop
			default:
				panic(err)
			}
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endRound(g, ai)
	}
	return g.balance
}

var (
	errBust = errors.New("hand score exceeds 21")
)

//Move moves state
type Move func(*Game) error

//MoveHit moves the game forward so that a new card is drawn
func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

//MoveSplit allows the user to split their hand if they have to cards of the same rank
func MoveSplit(g *Game) error {
	cards := g.currentHand()

	if len(*cards) != 2 {
		return errors.New("you can only split with 2 cards in your hand")
	}
	if (*cards)[0].Rank != (*cards)[1].Rank {
		return errors.New("both cards must have the same rank")
	}
	g.player = append(g.player, hand{
		cards: []deck.Card{(*cards)[1]},
		bet:   g.player[g.handIdx].bet,
	})
	g.player[g.handIdx].cards = (*cards)[:1]
	return nil
}

//MoveDouble allows the player to double down their hand.
func MoveDouble(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

//MoveStand moves the game state foward so that it can end.
func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}
	if g.state == statePlayerTurn {
		g.handIdx++
		if g.handIdx >= len(g.player) {
			g.state++
		}
		return nil
	}
	return errors.New("Invalid State")
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

//Score works out what the max score of the current cards is, useful if aces are present
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			//Currently, ace = 1, we change it to 11, only if the score is 11 or below
			return minScore + 10
		}
	}
	return minScore
}

//Soft is for a score of 17 if an ace is present
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score

}

//Blackjack works out if either the ai or user has blackjack
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func endRound(g *Game, ai AI) {
	dScore := Score(g.dealer...)
	dBlackJack := Blackjack(g.dealer...)
	allHands := make([][]deck.Card, len(g.player))
	for hi, hand := range g.player {
		allHands[hi] = hand.cards
		cards := hand.cards
		pScore, pBlackJack := Score(cards...), Blackjack(cards...)
		winnings := hand.bet
		switch {
		case pBlackJack && dBlackJack:
			winnings = 0
		case dBlackJack:
			winnings = -winnings
		case pBlackJack:
			winnings = int(float64(winnings) * g.blackjackPayout)
		case pScore > 21:
			winnings = -winnings
		case dScore > 21:
			//win
		case pScore > dScore:
			//win
		case dScore > pScore:
			winnings = -winnings
		case dScore == pScore:
			winnings = 0
		}
		g.balance += winnings
	}
	ai.Result(allHands, g.dealer)
	if Blackjack(g.dealer...) {

	}
	g.player = nil
	g.dealer = nil
}
