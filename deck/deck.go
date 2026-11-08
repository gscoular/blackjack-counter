package deck

import (
	"math/rand"
	"time"
)

type Card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
	Score int    `json:"score"`
}

type Shoe struct {
	Cards        []Card `json:"-"`
	RunningCount int    `json:"running_count"`
}

func NewShoe(numDecks int) *Shoe {

	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	var cards []Card
	for i := 0; i < numDecks; i++ {
		for _, suit := range suits {
			for _, value := range values {
				cards = append(cards, Card{
					Suit:  suit,
					Value: value,
					Score: getCountValue(value),
				})
			}
		}
	}

	shoe := Shoe{Cards: cards, RunningCount: 0}
	shoe.Shuffle()
	return &shoe
}

func (s *Shoe) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(s.Cards), func(i, j int) {
		s.Cards[i], s.Cards[j] = s.Cards[j], s.Cards[i]
	})
}

func (s *Shoe) Deal() (Card, bool) {
	if len(s.Cards) == 0 {
		return Card{}, false
	}
	card := s.Cards[0]
	s.Cards = s.Cards[1:]
	s.RunningCount = card.Score
	return card, true

}

func getCountValue(value string) int {
	switch value {
	case "2", "3", "4", "5", "6":
		return 1
	case "10", "J", "Q", "K", "A":
		return -1
	default:
		return 0
	}
}
