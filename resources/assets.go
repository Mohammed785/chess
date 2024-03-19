package resources

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

//go:embed images/*
var Images embed.FS

//go:embed audio/*
var Audio embed.FS

//go:embed FiraCodeNerdFont.ttf
var FiraCodeNerdFont_ttf []byte

func ReadImage(f string) *ebiten.Image {
	data, err := Images.ReadFile(path.Join("images", f))
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func ReadAudio(f string) *wav.Stream {
	audio, err := Audio.ReadFile(path.Join("audio", f))
	if err != nil {
		log.Fatal(err)
	}
	d, err := wav.DecodeWithoutResampling(bytes.NewReader(audio))
	if err != nil {
		log.Fatal(err)
	}
	return d
}
