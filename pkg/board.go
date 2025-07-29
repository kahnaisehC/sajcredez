package board

import (
	"fmt"
	"strings"
)

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

var pieceToChar = map[Piece]rune{
	// white pieces: ♕
	// black pieces: ♛
	EMPTY:   '_',
	WPAWN:   '♙',
	WKNIGHT: '♘',
	WBISHOP: '♗',
	WROOK:   '♖',
	WKING:   '♔',
	BPAWN:   '♟',
	BKNIGHT: '♞',
	BBISHOP: '♝',
	BROOK:   '♜',
	BKING:   '♚',
}

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

// NOTE: if it is EMPTY returns White by default
func GetPieceColor(p Piece) Color {
	if p >= WPAWN && p <= WKING {
		return WHITE
	} else if p >= BPAWN && p <= BKING {
		return BLACK
	}
	return WHITE
}

type Board interface {
	GetBoardString() string // Done
	GetHistoryString() string

	GetPiece(Square) Piece // Done
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

func CreateSajcredez() Sajcredez {
	s := Sajcredez{
		board:            default_board,
		whiteEnhances:    0,
		blackEnhances:    0,
		turn:             WHITE,
		BlackKingCastle:  true,
		BlackQueenCastle: true,
		WhiteKingCastle:  true,
		WhiteQueenCastle: true,
	}
	return s
}

func (s *Sajcredez) GetPiece(sq Square) Piece {
	return s.board[sq.row][sq.col]
}

// WARNING: This function is really slow due to constant string concatenation
func (s *Sajcredez) getBoardStringNoBuilder() string {
	boardString := ""
	boardString += fmt.Sprintf("whiteEnhances: %d blackEnhances: %d\ncastlingRights: ", s.whiteEnhances, s.blackEnhances)

	if s.WhiteKingCastle {
		boardString += "K"
	} else {
		boardString += "-"
	}

	if s.WhiteQueenCastle {
		boardString += "Q"
	} else {
		boardString += "-"
	}

	if s.BlackKingCastle {
		boardString += "k"
	} else {
		boardString += "-"
	}

	if s.BlackQueenCastle {
		boardString += "q"
	} else {
		boardString += "-"
	}
	boardString += "\n"
	boardString += "turn: "

	switch s.turn {
	case WHITE:
		boardString += "WHITE\n"
	case BLACK:
		boardString += "BLACK\n"

	}

	boardString += "\t"

	for col := range s.board {
		boardString += fmt.Sprintf("%d\t", col+1)
	}
	boardString += "\n"
	for row := range s.board {
		boardString += fmt.Sprintf("%d\t", row+1)
		for _, piece := range s.board[row] {
			boardString += fmt.Sprintf("%c\t", pieceToChar[piece])
		}
		boardString += "\n"
	}
	return boardString
}

func (s *Sajcredez) GetBoardString() string {
	boardStringBuilder := strings.Builder{}
	// the size of the initial boardString in bytes (using
	// piece emotes instead of chars) is 258
	estimatedBoardStringSize := 260
	boardStringBuilder.Grow(estimatedBoardStringSize)

	fmt.Fprintf(&boardStringBuilder, "whiteEnhances: %d blackEnhances: %d\ncastlingRights: ", s.whiteEnhances, s.blackEnhances)

	if s.WhiteKingCastle {
		boardStringBuilder.WriteByte('K')
	} else {
		boardStringBuilder.WriteByte('-')
	}

	if s.WhiteQueenCastle {
		boardStringBuilder.WriteByte('Q')
	} else {
		boardStringBuilder.WriteByte('-')
	}

	if s.BlackKingCastle {
		boardStringBuilder.WriteByte('k')
	} else {
		boardStringBuilder.WriteByte('-')
	}

	if s.BlackQueenCastle {
		boardStringBuilder.WriteByte('q')
	} else {
		boardStringBuilder.WriteByte('-')
	}
	boardStringBuilder.Write([]byte("\nturn: "))

	switch s.turn {
	case WHITE:
		boardStringBuilder.Write([]byte("WHITE\n"))
	case BLACK:
		boardStringBuilder.Write([]byte("BLACK\n"))

	}

	boardStringBuilder.WriteByte('\t')

	for col := range s.board {
		boardStringBuilder.WriteString(fmt.Sprintf("%d\t", col+1))
	}
	boardStringBuilder.WriteByte('\n')
	for row := range s.board {
		boardStringBuilder.WriteString(fmt.Sprintf("%d\t", row+1))
		for _, piece := range s.board[row] {
			boardStringBuilder.WriteRune(pieceToChar[piece])
			boardStringBuilder.WriteByte('\t')
		}
		boardStringBuilder.WriteByte('\n')
	}
	fmt.Println(len(boardStringBuilder.String()))
	return boardStringBuilder.String()
}
