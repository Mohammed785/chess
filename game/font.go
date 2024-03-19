package game

import (
	"log"

	"github.com/Mohammed785/chess/resources"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	FiraNormal font.Face
	FiraBig    font.Face
)

func init() {
	tt, err := opentype.Parse(resources.FiraCodeNerdFont_ttf)
	if err != nil {
		log.Fatalln(err)
	}
	const dpi = 72
	FiraNormal, err = opentype.NewFace(tt, &opentype.FaceOptions{
		DPI:     dpi,
		Size:    24,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatalln(err)
	}
	FiraBig, err = opentype.NewFace(tt, &opentype.FaceOptions{
		DPI:     dpi,
		Size:    48,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalln(err)
	}
	FiraBig = text.FaceWithLineHeight(FiraBig, 54)
}
