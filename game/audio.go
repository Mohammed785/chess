package game

import (
	"log"

	"github.com/Mohammed785/chess/resources"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sample = 44100

type Audio struct {
	ctx           *audio.Context
	movePlayer    *audio.Player
	checkPlayer   *audio.Player
	capturePlayer *audio.Player
	promotePlayer *audio.Player
	castlePlayer  *audio.Player
}

func NewAudio() Audio {
	audio := Audio{ctx: audio.NewContext(sample)}
	audio.capturePlayer = newPlayer(audio.ctx, "capture.wav")
	audio.movePlayer = newPlayer(audio.ctx, "move.wav")
	audio.promotePlayer = newPlayer(audio.ctx, "promote.wav")
	audio.castlePlayer = newPlayer(audio.ctx, "castle.wav")
	audio.checkPlayer = newPlayer(audio.ctx, "check.wav")
	return audio
}

func newPlayer(ctx *audio.Context, audio string) *audio.Player {
	p, err := ctx.NewPlayer(resources.ReadAudio(audio))
	if err != nil {
		log.Fatalf("couldn't load '%s': %s\n", audio, err.Error())
	}
	return p
}

func (a Audio) PlayMoveAudio()error {
	if err:=a.movePlayer.Rewind();err!=nil{
		return err
	}
	a.movePlayer.Play()
	return nil
}

func (a Audio) PlayCaptureAudio()error {
	if err:=a.capturePlayer.Rewind();err!=nil{
		return err
	}
	a.capturePlayer.Play()
	return nil
}

func (a Audio) PlayCastleAudio()error {
	if err:=a.castlePlayer.Rewind();err!=nil{
		return err
	}
	a.castlePlayer.Play()
	return nil
}

func (a Audio) PlayPromoteAudio()error {
	if err:=a.promotePlayer.Rewind();err!=nil{
		return err
	}
	a.promotePlayer.Play()
	return nil
}

func (a Audio) PlayCheckAudio()error {
	if err:=a.checkPlayer.Rewind();err!=nil{
		return err
	}
	a.checkPlayer.Play()
	return nil
}


