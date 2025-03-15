package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	Simg            []ebiten.Image
	buttonState     int
	bX, bY          int
	bHeight, bWidth int
	label           string
}

func (btn *Button) isMouseOverPos() bool {
	mX, mY := ebiten.CursorPosition()
	if ((mX < btn.bX+btn.bWidth) && (mX > btn.bX)) && ((mY < btn.bY+btn.bHeight) && (mY > btn.bY)) {
		return true
	}
	return false
}
func (btn *Button) Update(g *Game) {
	//ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
		btn.buttonState = 2
	} else if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
		btn.buttonState = 1
	} else {
		btn.buttonState = 0
	}
}
func (btn *Button) Draw(screen *ebiten.Image, g *Game) {
	//w, h := btn.Simg[btn.buttonState].Bounds().Dx(), btn.Simg[btn.buttonState].Bounds().Dy()
	g.op.GeoM.Reset()
	scaleX, scaleY := 2.0, 2.0
	//g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	//g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	g.op.GeoM.Scale(scaleX, scaleY)
	g.op.GeoM.Translate(float64(btn.bX), float64(btn.bY))
	screen.DrawImage(&btn.Simg[btn.buttonState], &g.op)
	g.op.GeoM.Reset()
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(btn.bX+(fontSize0/2)), float64(btn.bY+(fontSize0/2)))
	op.ColorScale.ScaleWithColor(color.RGBA{250, 250, 250, 255})
	op.LineSpacing = fontSize0
	text.Draw(screen, btn.label, textFaceMono, op)
}
