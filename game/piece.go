package game

import (
	"fmt"
	"math"

	"github.com/Mohammed785/chess/resources"
	"github.com/hajimehoshi/ebiten/v2"
)

type MoveType int

type Piece struct {
	Notation rune
	Color    string
	pos      [2]int
	oldPos   [2]int
	hasMoved bool
	captured bool
	Img      *ebiten.Image
}

func NewPiece(y, x int) *Piece {
	var pieceImg *ebiten.Image
	if x > 1 && x < 6 {
		return nil
	}
	color := "black"
	notation := ' '

	if x > 5 {
		color = "white"
	}
	if x == 1 || x == 6 {
		pieceImg = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_pawn.png", color)))
	} else if y == 0 || y == 7 {
		pieceImg = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_rook.png", color)))
		notation = 'R'
	} else if y == 1 || y == 6 {
		pieceImg = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_knight.png", color)))
		notation = 'N'
	} else if y == 2 || y == 5 {
		pieceImg = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_bishop.png", color)))
		notation = 'B'
	} else if y == 3 {
		pieceImg = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_queen.png", color)))
		notation = 'Q'
	} else if y == 4 {
		pieceImg = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_king.png", color)))
		notation = 'K'
	}
	return &Piece{
		Img:      pieceImg,
		Color:    color,
		Notation: notation,
		pos:      [2]int{y, x},
		oldPos:   [2]int{-1, -1},
		hasMoved: false,
		captured: false,
	}
}

func (p *Piece) Draw(screen *ebiten.Image) {
	if !p.captured {
		op := &ebiten.DrawImageOptions{}
		tx := float64(p.pos[0]*int(SQUARE_SIZE)) + float64(p.Img.Bounds().Min.X)
		ty := float64(p.pos[1]*int(SQUARE_SIZE)) + float64(p.Img.Bounds().Min.Y)
		op.GeoM.Translate(tx, ty)
		screen.DrawImage(p.Img, op)
	}
}
func (p *Piece) Pos() [2]int {
	return p.pos
}

func (p *Piece) RevertMove(other *Piece) {
	temp := p.pos
	p.pos = p.oldPos
	p.oldPos = temp
	if other != nil {
		other.captured = false
		other.RevertMove(nil)
	}
}

func handleCastling(king, rook *Piece) {
	king.hasMoved = true
	rook.hasMoved = true
	var newPositions [2][2]int
	if king.Color == "white" {
		if rook.Pos() == [2]int{7, 7} {
			newPositions = [2][2]int{{6, 7}, {5, 7}}
		} else {
			newPositions = [2][2]int{{2, 7}, {3, 7}}
		}
	} else {
		if rook.Pos() == [2]int{7, 0} {
			newPositions = [2][2]int{{6, 0}, {5, 0}}
		} else {
			newPositions = [2][2]int{{2, 0}, {3, 0}}
		}
	}
	king.SetNewPos(newPositions[0])
	rook.SetNewPos(newPositions[1])

}
func (p *Piece) SetNewPos(newPos [2]int) {
	p.oldPos = p.pos
	p.pos = [2]int{newPos[0], newPos[1]}
}

func (p *Piece) Move(newPos [2]int, other *Piece, confirm bool) string {
	if other != nil && p.Color == other.Color && confirm { //the only possible way for this is castling
		handleCastling(p, other)
		return "castle"
	}
	moveType := "move"
	p.oldPos = p.pos
	p.pos = newPos 
	if other != nil {
		other.captured = true
		other.oldPos = other.pos
		other.pos = [2]int{-1, -1}
		moveType = "capture"
	}
	if confirm {
		p.hasMoved = confirm
		if p.Notation == ' ' && (newPos[1] == 0 || newPos[1] == 7) {
			p.Img = ebiten.NewImageFromImage(resources.ReadImage(fmt.Sprintf("%v_queen.png", p.Color)))
			p.Notation = 'Q'
			moveType = "promote"
		}
	}
	return moveType
}

func (p *Piece) CanSeeKing(b *Board) bool {
	enemyKing := b.GetBKingPos()
	if p.Color == "black" {
		enemyKing = b.GetWKingPos()
	}
	moves := p.GetAllMoves(b)
	for _, pos := range moves {
		if pos == enemyKing {
			return true
		}
	}
	return false
}

func (p *Piece) GetAllMoves(b *Board) [][2]int {

	switch p.Notation {
	case ' ':
		return PawnMoves(b, p)
	case 'N':
		return KnightMoves(b, p)
	case 'B':
		return BishopMoves(b, p)
	case 'R':
		return RookMoves(b, p)
	case 'Q':
		return QueenMoves(b, p)
	case 'K':
		return KingMoves(b, p)
	}
	return nil
}

func validatePos(pos [2]int) bool {
	if pos[0] < 0 || pos[0] > 7 || pos[1] < 0 || pos[1] > 7 {
		return false
	}
	return true
}

func PawnMoves(board *Board, piece *Piece) [][2]int {
	moves := make([][2]int, 0, 3)
	y := -1
	if piece.Color == "black" {
		y = 1
	}
	for x := -1; x < 2; x++ {
		if pos := [2]int{piece.pos[0] + x, piece.pos[1] + y}; validatePos(pos) {
			targetPiece := board.GetPiece(pos)
			if (x != 0 && (targetPiece != nil && targetPiece.Color != piece.Color)) ||
				(x == 0 && targetPiece == nil) {
				moves = append(moves, pos)
			}
		}
	}
	// 2 square moves
	if pos := [2]int{piece.pos[0], piece.pos[1] + (y * 2)}; validatePos(pos) && !piece.hasMoved && board.GetPiece(pos) == nil {
		moves = append(moves, pos)
	}
	pieceY := piece.Pos()[1]
	if pieceY == 3 && piece.Color == "white" && board.PawnDoubleMove.Player == "white" {
		if board.PawnDoubleMove.Pos == [2]int{piece.pos[0] + 1, pieceY} || board.PawnDoubleMove.Pos == [2]int{piece.pos[0] - 1, pieceY} {
			moves = append(moves, [2]int{board.PawnDoubleMove.Pos[0], pieceY + y})
		}
	} else if pieceY == 4 && piece.Color == "black" && board.PawnDoubleMove.Player == "black" {
		if board.PawnDoubleMove.Pos == [2]int{piece.pos[0] + 1, pieceY} || board.PawnDoubleMove.Pos == [2]int{piece.pos[0] - 1, pieceY} {
			moves = append(moves, [2]int{board.PawnDoubleMove.Pos[0], pieceY + y})
		}
	}
	return moves
}

func KnightMoves(board *Board, piece *Piece) [][2]int {
	moves := make([][2]int, 0, 4)
	candidates := [][2]int{{1, 2}, {-1, 2}, {1, -2}, {-1, -2}, {2, -1}, {2, 1}, {-2, 1}, {-2, -1}}

	for _, c := range candidates {
		if pos := [2]int{piece.pos[0] + c[0], piece.pos[1] + c[1]}; validatePos(pos) {
			targetPiece := board.GetPiece(pos)
			if targetPiece == nil || targetPiece.Color != piece.Color {
				moves = append(moves, pos)
			}
		}
	}
	return moves
}

func BishopMoves(board *Board, piece *Piece) [][2]int {
	moves := make([][2]int, 0, 4)
	blocked := [4]int8{0, 0, 0, 0}
	validateMove := func(x, y, i int) {
		if blocked[i] == 1 {
			return
		}
		if pos := [2]int{piece.pos[0] + x, piece.pos[1] + y}; validatePos(pos) {
			targetPiece := board.GetPiece(pos)
			if targetPiece != nil {
				blocked[i] = 1
			}
			if targetPiece == nil || targetPiece.Color != piece.Color {
				moves = append(moves, pos)
			}
		}
	}
	for x := 1; x < 8; x++ {
		validateMove(x, x, 0)
		validateMove(x, -x, 1)
		validateMove(-x, x, 2)
		validateMove(-x, -x, 3)
	}
	return moves
}

func RookMoves(board *Board, piece *Piece) [][2]int {
	moves := make([][2]int, 0, 4)
	blocked := [4]int8{0, 0, 0, 0}
	validateMove := func(x, y, i int) {
		if blocked[i] == 1 {
			return
		}
		if pos := [2]int{piece.pos[0] + x, piece.pos[1] + y}; validatePos(pos) {
			targetPiece := board.GetPiece(pos)
			if targetPiece != nil {
				blocked[i] = 1
			}
			if targetPiece == nil || targetPiece.Color != piece.Color {
				moves = append(moves, pos)
			}
		}
	}
	for x := 1; x < 8; x++ {
		validateMove(x, 0, 0)
		validateMove(-x, 0, 1)
		validateMove(0, -x, 2)
		validateMove(0, x, 3)
	}
	return moves
}
func QueenMoves(board *Board, piece *Piece) [][2]int {
	moves := RookMoves(board, piece)
	moves = append(moves, BishopMoves(board, piece)...)
	return moves
}

func KingMoves(board *Board, king *Piece) [][2]int {
	moves := make([][2]int, 0, 4)
	candidates := [][2]int{{1, 1}, {1, 0}, {1, -1}, {0, -1}, {0, 1}, {-1, 1}, {-1, -1}, {-1, 0}}
	for _, c := range candidates {
		if pos := [2]int{king.pos[0] + c[0], king.pos[1] + c[1]}; validatePos(pos) {
			targetPiece := board.GetPiece(pos)
			if targetPiece == nil || targetPiece.Color != king.Color {
				moves = append(moves, pos)
			}
		}
	}
	if !king.hasMoved && king.Color != board.inCheck {
		var kingSideRook, queenSideRook *Piece

		emptyPath := func(from, to [2]int) bool {
			start := int(math.Min(float64(from[0]+1), float64(to[0]+1)))
			for x := start; x <= start; x++ {
				if board.GetPiece([2]int{x, to[1]}) != nil {
					return false
				}

			}
			return true
		}
		if king.Color == "white" {
			kingSideRook = board.GetPiece([2]int{7, 7})
			queenSideRook = board.GetPiece([2]int{0, 7})
		} else {
			kingSideRook = board.GetPiece([2]int{7, 0})
			queenSideRook = board.GetPiece([2]int{0, 0})
		}
		if kingSideRook != nil && !kingSideRook.hasMoved {
			if emptyPath(king.Pos(), kingSideRook.Pos()) {
				moves = append(moves, kingSideRook.Pos())
			}
		}
		if queenSideRook != nil && !queenSideRook.hasMoved {
			if emptyPath(king.Pos(), queenSideRook.Pos()) {
				moves = append(moves, queenSideRook.Pos())
			}
		}
	}
	return moves
}
