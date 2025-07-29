package board

import (
	"errors"
	"fmt"
	"strings"
)

type Piece int8

const (
	sajcredezCols = int8(7)
	sajcredezRows = int8(7)
)

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

// piece moves and slides

var kingMoves = []Square{
	{1, 1},
	{1, 0},
	{1, -1},
	{0, 1},
	{0, -1},
	{-1, 1},
	{-1, 0},
	{-1, -1},
}

var knightMoves = []Square{
	{2, 1},
	{-2, 1},
	{-2, -1},
	{2, -1},
	{1, 2},
	{-1, 2},
	{-1, -2},
	{1, -2},
}

var rookSlides = []Square{
	{0, -1},
	{0, 1},
	{1, 0},
	{-1, 0},
}

var bishopSlides = []Square{
	{1, 1},
	{-1, 1},
	{-1, -1},
	{1, -1},
}

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

var default_board = [sajcredezRows][sajcredezCols]Piece{
	{WROOK, WKNIGHT, WBISHOP, WKING, WBISHOP, WKNIGHT, WROOK},
	{WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN, WPAWN},
	{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	{EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY},
	{BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN, BPAWN},
	{BROOK, BKNIGHT, BBISHOP, BKING, BBISHOP, BKNIGHT, BROOK},
}

type Options struct{}

type Color int8

const (
	NO_COLOR Color = 0
	WHITE    Color = 1
	BLACK    Color = 2
)

type Enhance int8

const (
	NO_ENHANCE Enhance = (iota)
	ENHANCE_MOVE
)

type Square struct {
	col int8
	row int8
}

func addSquares(a, b Square) (sum Square) {
	sum = Square{
		col: a.col + b.col,
		row: a.row + b.row,
	}
	return
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
	return NO_COLOR
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

	board         [sajcredezRows][sajcredezCols]Piece
	whiteEnhances int
	blackEnhances int
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
	return boardStringBuilder.String()
}

func ParseMove(str string) (Move, error) {
	move := Move{}
	return move, nil
}

func (s *Sajcredez) inBounds(sq Square) bool {
	// sq.col is between this two values AND
	colsLowEnd := int8(0)
	colsHighEnd := sajcredezCols

	// sq.row is between this two values
	rowsLowEnd := int8(0)
	rowsHighEnd := sajcredezRows

	return colsLowEnd <= sq.col && sq.col < colsHighEnd && rowsLowEnd <= sq.row && sq.row < rowsHighEnd
}

func (s *Sajcredez) CheckMoveLegality(move Move) error {
	const errorHeader = "CheckMoveLegality: "

	// check if from is inbounds
	if !s.inBounds(move.from) {
		return errors.New(errorHeader + "from is out of bounds")
	}
	// check if to is inbounds
	if !s.inBounds(move.to) {
		return errors.New(errorHeader + "to is out of bounds")
	}

	// check if fromPiece is correct
	if move.fromPiece != s.GetPiece(move.from) {
		return errors.New(errorHeader + "fromPiece is not the piece in from")
	}

	// check if toPiece is correct
	if move.toPiece != s.GetPiece(move.to) {
		return errors.New(errorHeader + "toPiece is not the piece in to")
	}

	// check if it is the correct turn to play
	if s.turn != GetPieceColor(s.GetPiece(move.from)) {
		return errors.New(errorHeader + "incorrect turn to play")
	}

	// check if you have enhancement available
	if move.enhancement != NO_ENHANCE {
		var enhancesAvailable int
		switch s.turn {
		case BLACK:
			enhancesAvailable = s.blackEnhances
		case WHITE:
			enhancesAvailable = s.whiteEnhances
		}
		if enhancesAvailable == 0 {
			return errors.New(errorHeader + "not enough enhances available")
		}
	}

	// check if fromPiece and toPiece are of the same color
	ignoreIfFromPieceToPieceSameColor := move.enhancement == ENHANCE_MOVE && (move.fromPiece == WKNIGHT || move.fromPiece == BKNIGHT)
	if !ignoreIfFromPieceToPieceSameColor && GetPieceColor(move.fromPiece) == GetPieceColor(move.toPiece) {
		return errors.New(errorHeader + "fromPiece and toPiece are from the same color")
	}

	// now comes the fun part
	canMove := false
	switch move.enhancement {
	case NO_ENHANCE:
		{
			switch move.fromPiece {
			case WKNIGHT:
				fallthrough
			case BKNIGHT:
				for _, nm := range knightMoves {
					if addSquares(nm, move.from) == move.to {
						canMove = true
					}
				}
			case WKING:
				fallthrough
			case BKING:
				for _, km := range kingMoves {
					if addSquares(km, move.from) == move.to {
						canMove = true
					}
				}
			case WROOK:
				fallthrough
			case BROOK:
				for _, rs := range rookSlides {
					for nextSquare := addSquares(rs, move.from); s.inBounds(nextSquare); nextSquare = addSquares(rs, nextSquare) {
						if GetPieceColor(s.GetPiece(nextSquare)) == GetPieceColor(move.fromPiece) {
							break
						}
						if nextSquare == move.from {
							canMove = true
						}
						if GetPieceColor(s.GetPiece(nextSquare)) != NO_COLOR {
							break
						}
					}
				}
			case WBISHOP:
				fallthrough
			case BBISHOP:
				for _, bs := range bishopSlides {
					for nextSquare := addSquares(bs, move.from); s.inBounds(nextSquare); nextSquare = addSquares(bs, nextSquare) {
						if GetPieceColor(s.GetPiece(nextSquare)) == GetPieceColor(move.fromPiece) {
							break
						}
						if nextSquare == move.from {
							canMove = true
						}
						if GetPieceColor(s.GetPiece(nextSquare)) != NO_COLOR {
							break
						}
					}
				}
			case WPAWN:
				squareUp := Square{col: move.from.col + 1, row: move.from.row}
				if move.to == squareUp && s.GetPiece(move.to) == EMPTY {
					canMove = true
				}
				squarePositiveDiagonal := Square{col: move.from.col + 1, row: move.from.row + 1}
				if move.to == squarePositiveDiagonal && move.toPiece != EMPTY {
					canMove = true
				}
				squareNegativeDiagonal := Square{col: move.from.col + 1, row: move.from.row - 1}
				if move.to == squareNegativeDiagonal && move.toPiece != EMPTY {
					canMove = true
				}

			case BPAWN:
				squareDown := Square{col: move.from.col - 1, row: move.from.row}
				if move.to == squareDown && s.GetPiece(move.to) == EMPTY {
					canMove = true
				}
				squarePositiveDiagonal := Square{col: move.from.col + 1, row: move.from.row + 1}
				if move.to == squarePositiveDiagonal && move.toPiece != EMPTY {
					canMove = true
				}
				squareNegativeDiagonal := Square{col: move.from.col + 1, row: move.from.row - 1}
				if move.to == squareNegativeDiagonal && move.toPiece != EMPTY {
					canMove = true
				}

			default:
				return errors.New(errorHeader + "invalid fromPiece code (probably trying to move a piece from an EMPTY square)")
			}
		}

	case ENHANCE_MOVE:
		{
			switch move.fromPiece {
			case WKNIGHT:
				fallthrough
			case BKNIGHT:
				for _, nm := range knightMoves {
					if addSquares(nm, move.from) == move.to {
						canMove = true
					}
				}
			case WKING:
				fallthrough
			case BKING:
				for _, km := range kingMoves {
					if addSquares(km, move.from) == move.to {
						canMove = true
					}
				}
			case WROOK:
				fallthrough
			case BROOK:
				for _, rs := range rookSlides {
					for nextSquare := addSquares(rs, move.from); s.inBounds(nextSquare); nextSquare = addSquares(rs, nextSquare) {
						if GetPieceColor(s.GetPiece(nextSquare)) == GetPieceColor(move.fromPiece) {
							break
						}
						if nextSquare == move.from {
							canMove = true
						}
					}
				}
			case WBISHOP:
				fallthrough
			case BBISHOP:
				for _, bs := range bishopSlides {

					// make slide go "through" the board
					nextSquare := addSquares(bs, move.from)
					nextSquare.col = (nextSquare.col + sajcredezCols) % sajcredezCols
					nextSquare.row = (nextSquare.row + sajcredezRows) % sajcredezRows

					count := 0
					maxiterations := int(sajcredezRows) * int(sajcredezCols)

					for s.inBounds(nextSquare) {
						if count > maxiterations {
							return errors.New(errorHeader + "infinite iteration on bishop slides (implementation error, this shouldnt happen)")
						}
						if GetPieceColor(s.GetPiece(nextSquare)) == GetPieceColor(move.fromPiece) {
							break
						}
						if nextSquare == move.from {
							canMove = true
						}
						if GetPieceColor(s.GetPiece(nextSquare)) != NO_COLOR {
							break
						}

						nextSquare := addSquares(bs, move.from)
						nextSquare.col = (nextSquare.col + sajcredezCols) % sajcredezCols
						nextSquare.row = (nextSquare.row + sajcredezRows) % sajcredezRows
					}
				}
			case WPAWN:
				squareUp := Square{col: move.from.col + 1, row: move.from.row}
				if move.to == squareUp {
					canMove = true
				}
				squarePositiveDiagonal := Square{col: move.from.col + 1, row: move.from.row + 1}
				if move.to == squarePositiveDiagonal {
					canMove = true
				}
				squareNegativeDiagonal := Square{col: move.from.col + 1, row: move.from.row - 1}
				if move.to == squareNegativeDiagonal {
					canMove = true
				}

			case BPAWN:
				squareDown := Square{col: move.from.col - 1, row: move.from.row}
				if move.to == squareDown {
					canMove = true
				}
				squarePositiveDiagonal := Square{col: move.from.col + 1, row: move.from.row + 1}
				if move.to == squarePositiveDiagonal {
					canMove = true
				}
				squareNegativeDiagonal := Square{col: move.from.col + 1, row: move.from.row - 1}
				if move.to == squareNegativeDiagonal {
					canMove = true
				}

			default:
				return errors.New(errorHeader + "invalid fromPiece code (probably trying to move a piece from an EMPTY square)")
			}
		}
	}
	if !canMove {
		// specialize the error somehow
		error := "move is illegal due to something (improve this message)"
		return errors.New(errorHeader + error)
	}

	return nil
}
