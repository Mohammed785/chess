package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/Mohammed785/chess/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Mode int

const (
	modeGame Mode = iota
	modeGameOver
)

type Game struct {
	board         *game.Board
	audio         game.Audio
	currentPlayer string
	validMoves    [][2]int
	selectedPos   [2]int
	gameMode      Mode
	winner        string
}

func (g *Game) ChangePlayer() {
	if g.currentPlayer == "white" {
		g.currentPlayer = "black"
	} else {
		g.currentPlayer = "white"
	}
	g.validMoves = nil
	g.selectedPos = [2]int{-1, -1}
}

func (g *Game) IsValid(move [2]int) bool {
	for _, v := range g.validMoves {
		if v == move {
			return true
		}
	}
	return false
}

func (g *Game) RestartGame(){
	g.board = game.NewBoard()
	g.currentPlayer = "white"
	g.gameMode = modeGame
	g.winner=""
	g.validMoves = nil
}

func (g *Game) Update() error {
	switch g.gameMode {
	case modeGameOver:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.RestartGame()
		}
	case modeGame:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			x /= int(game.SQUARE_SIZE)
			y /= int(game.SQUARE_SIZE)
			pos := [2]int{x, y}
			piece := g.board.GetPiece([2]int{x, y})

			if g.validMoves != nil && g.IsValid(pos) {
				p := g.board.GetPiece(g.selectedPos)
				if p.Notation==' '&&piece==nil&&x!=p.Pos()[0]{
					piece = g.board.GetPiece(g.board.PawnDoubleMove.Pos)
				}
				moveType := p.Move(pos, piece, true)
				tempPlayer := g.currentPlayer
				tempPos := g.selectedPos
				g.ChangePlayer()
				if p.Notation==' '&&math.Abs(float64(pos[1]-tempPos[1]))==2{
					g.board.PawnDoubleMove = game.PawnDoubleMove{
						Player: g.currentPlayer,
						Pos: pos,
					}
				}else{
					g.board.PawnDoubleMove.Player=""

				}
				cantMove:=g.board.IsCheckmated(g.currentPlayer)
				if g.board.IsChecked(g.currentPlayer) {
					g.playSound("check")
					if cantMove{
						g.winner = tempPlayer
						g.gameMode = modeGameOver
					}
				} else {
					if cantMove{
						g.winner = "draw"
						g.gameMode = modeGameOver
					}
					g.playSound(moveType)
				}
			} else if piece != nil && piece.Color == g.currentPlayer {
				allMoves := piece.GetAllMoves(g.board)
				g.validMoves = g.board.ValidateMoves(g.currentPlayer, piece, allMoves)
				g.selectedPos = pos
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.DrawSquares(screen)
	g.board.DrawPieces(screen)
	g.board.HighlightSquares(screen, g.validMoves)

	if g.gameMode == modeGameOver {
		if g.winner=="draw"{
			text.Draw(screen, "Game ended in a Draw", game.FiraBig, 150, 300, color.Black)
		}else{
			text.Draw(screen, fmt.Sprintf("%s Won!!", g.winner), game.FiraBig, 280, 300, color.Black)
		}
		text.Draw(screen, "Press 'Enter' to start a new game", game.FiraNormal, 150, 400, color.Black)
	}
}
func (g *Game) playSound(sound string) {
	switch sound {
	case "move":
		g.audio.PlayMoveAudio()
	case "capture":
		g.audio.PlayCaptureAudio()
	case "castle":
		g.audio.PlayCastleAudio()
	case "promote":
		g.audio.PlayPromoteAudio()
	case "check":
		g.audio.PlayCheckAudio()
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}

func NewGame() *Game {
	return &Game{board: game.NewBoard(),audio: game.NewAudio(), currentPlayer: "white", winner: "", gameMode: modeGame}
}

func main() {
	ebiten.SetWindowSize(800, 800)
	ebiten.SetWindowTitle("Chess")
	if err := ebiten.RunGame(NewGame()); err != nil {

		log.Fatal(err)
	}
}
