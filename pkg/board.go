package board

type Options struct{}

type Square int

type Move struct{}

type Board interface {
	GetListOfMoves(options Options) []Move
	MakeMove(Move) error
	GetBoardString() string
}

type Sajcredez struct {
	board [7][7]int
}
