package game

import (
	"fmt"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const SQUARE_SIZE float32 = 100

var rankNames = [8]string{"a", "b", "c", "d", "e", "f", "g", "h"}

// help with en passant
type PawnDoubleMove struct {
	// the player who can capture an en passant
	Player string
	// the position he can capture on
	Pos [2]int
}

type Board struct {
	Squares        [][]Square
	pieces         [2][16]*Piece
	kingsPositions [2]int
	inCheck        string
	PawnDoubleMove	PawnDoubleMove
}

func NewBoard() *Board {
	var c color.Color
	board := new(Board)
	board.inCheck = ""
	board.Squares = make([][]Square, 8)
	w := 0
	b := 0
	for x := 0; x < 8; x++ {
		board.Squares[x] = make([]Square, 8)
		for y := 0; y < 8; y++ {
			if (x+y)%2 == 0 {
				c = color.RGBA{R: 242, G: 232, B: 231, A: 255}
			} else {
				c = color.RGBA{R: 163, G: 82, B: 78, A: 255}
			}
			if piece := NewPiece(y, x); piece != nil {
				if piece.Color == "white" {
					board.pieces[0][w] = piece
					if piece.Notation == 'K' {
						board.kingsPositions[0] = w
					}
					w++
				} else {
					board.pieces[1][b] = piece
					if piece.Notation == 'K' {
						board.kingsPositions[1] = b
					}
					b++
				}
			}
			board.Squares[x][y] = Square{x: x, y: y, color: c}
		}
	}
	return board
}

func (b *Board) DrawSquares(screen *ebiten.Image) {
	checkedKing := [2]int{-1, -1}
	if b.inCheck == "black" {
		checkedKing = b.GetBKingPos()
	} else if b.inCheck == "white" {
		checkedKing = b.GetWKingPos()
	}
	for _, ranks := range b.Squares {
		for _, square := range ranks {
			square.Draw(screen, checkedKing[0] != -1 && square.Pos() == checkedKing)
		}
	}
}
func (b *Board) DrawPieces(screen *ebiten.Image) {
	for color := range b.pieces {
		for i := range b.pieces[color] {
			b.pieces[color][i].Draw(screen)
		}
	}
}

func (b *Board) GetPiece(pos [2]int) *Piece {
	for c := range b.pieces {
		for _, piece := range b.pieces[c] {
			if piece.pos == pos {
				return piece
			}
		}
	}
	return nil
}

func (b *Board) IsChecked(currentPlayer string) bool {
	idx := 0
	if currentPlayer == "white" {
		idx = 1
	}
	for _, piece := range b.pieces[idx] {
		if piece.CanSeeKing(b) {
			b.inCheck = currentPlayer
			return true
		}
	}
	b.inCheck = ""
	return false
}

func (b *Board) IsCheckmated(player string) bool {
	idx := 1
	if player == "white" {
		idx = 0
	}
	for _, piece := range b.pieces[idx] {
		if len(b.ValidateMoves(player, piece, piece.GetAllMoves(b))) > 0 {
			return false
		}
	}
	return true
}

func (b *Board) GetSquare(x, y int) Square {
	return b.Squares[x][y]
}

func (b *Board) ValidateMoves(currentPlayer string, piece *Piece, moves [][2]int) [][2]int {
	validMoves := make([][2]int, 0, len(moves))
	if piece == nil {
		return nil
	}
	for i := range moves {
		if b.ValidateMove(currentPlayer, piece, moves[i]) {
			validMoves = append(validMoves, moves[i])
		}
	}
	return validMoves
}

func (b *Board) ValidateMove(currentPlayer string, piece *Piece, otherPos [2]int) bool {
	otherPiece := b.GetPiece(otherPos)
	piece.Move(otherPos, otherPiece, false)
	defer piece.RevertMove(otherPiece)
	enemyIdx := 1
	if currentPlayer == "black" {
		enemyIdx = 0
	}
	ePieces := b.pieces[enemyIdx]
	for _, enemyPiece := range ePieces {
		if !enemyPiece.captured && enemyPiece.CanSeeKing(b) {
			return false
		}
	}
	return true
}

func (b *Board) GetBKingPos() [2]int {
	return b.pieces[1][b.kingsPositions[1]].Pos()
}

func (b *Board) GetWKingPos() [2]int {
	return b.pieces[0][b.kingsPositions[0]].Pos()
}

func (b *Board) HighlightSquares(screen *ebiten.Image, squares [][2]int) {
	for _, pos := range squares {
		b.GetSquare(pos[0], pos[1]).Highlight(screen)
	}
}

type Square struct {
	x, y  int
	color color.Color
}

func (s Square) Highlight(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(s.x*int(SQUARE_SIZE))+50., float32(s.y*int(SQUARE_SIZE))+50.0, 50, color.RGBA{20, 20, 20, 90}, true)
}

func (s Square) Draw(screen *ebiten.Image, check bool) {
	c := s.color
	if check {
		c = color.RGBA{189, 2, 2, 100}
	}
	vector.DrawFilledRect(screen, float32(s.x)*SQUARE_SIZE, float32(s.y)*SQUARE_SIZE, SQUARE_SIZE, SQUARE_SIZE, c, true)
	if s.x == 0 {
		text.Draw(screen, fmt.Sprintf("%d", 8-s.y), FiraNormal, 0, (s.y*int(SQUARE_SIZE))+25, color.Black)
	}
	if s.y == 7 {
		text.Draw(screen, rankNames[s.x], FiraNormal, (s.x*int(SQUARE_SIZE))+80, 800, color.Black)
	}
}

func (s Square) DrawCheck(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(s.x)*SQUARE_SIZE, float32(s.y)*SQUARE_SIZE, SQUARE_SIZE, SQUARE_SIZE, color.RGBA{189, 2, 2, 100}, true)
}

func (s Square) Pos() [2]int {
	return [2]int{s.x, s.y}
}
