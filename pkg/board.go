package board

type Piece int8

const (
	EMPTY Piece = iota
	WPAWN
	WKNIGHT
	WBISHOP
	WROOK
	WKING
	BPAWN
	BKNIGHT
	BBISHOP
	BROOK
	BKING
)

var default_board = [7][7]Piece{
	{WROOK, WKNIGHT, WBISHOP, WKING, WBISHOP, WKNIGHT, WROOK},
	{WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN},
	{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	{BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN},
	{BROOK, BKNIGHT, BBISHOP, BKING, BBISHOP, BKNIGHT, BROOK},
}

type Options struct{}

type Color bool

const (
	WHITE Color = false
	BLACK Color = true
)

type Enhance int8

const (
	NO_ENHANCE Enhance = (iota)
	ENHANCE_MOVEMENT
)

type Square struct {
	col int8
	row int8
}

type Move struct {
	from      Square
	to        Square
	fromPiece Piece
	toPiece   Piece

	enhancement Enhance
	promotion   Piece
}

type Board interface {
	GetBoardString() string
	GetHistoryString() string

	GetPieceColor(Piece) Color
	GetPiece(Square) Piece
	SquareIsThreattened(Color) []Square
	CheckMoveLegality(Move) error

	GetListOfMoves() []Move
	MakeMove(Move) error
}

type Sajcredez struct {
	History []Move

	board         [7][7]Piece
	whiteEnhances Enhance
	blackEnhances Enhance
	turn          Color

	BlackKingCastle  bool
	BlackQueenCastle bool
	WhiteKingCastle  bool
	WhiteQueenCastle bool

	// Make move skips move validation
	// in order to prevent having to check if a move is legal
	// for engines or uis
	skipValidation bool
}
