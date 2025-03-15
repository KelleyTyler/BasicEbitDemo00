package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	Simg            []ebiten.Image
	buttonState     int
	bX, bY          int
	bHeight, bWidth int
	label           string
	isLocking       bool
	isLocked        bool
	scaleX, scaleY  float32
	msg             string
}

func (btn *Button) isMouseOverPos() bool {
	mX, mY := ebiten.CursorPosition()
	if ((mX < btn.bX+btn.bWidth) && (mX > btn.bX)) && ((mY < btn.bY+btn.bHeight) && (mY > btn.bY)) {
		return true
	}
	return false
}
func (btn *Button) Update() {
	//ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	//btn.msg = fmt.Sprintf("%s %d\n%d %t", btn.label, len(btn.Simg), btn.buttonState, btn.isMouseOverPos())
	if btn.isLocking {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			btn.buttonState = 2
			btn.isLocked = !btn.isLocked
			//fmt.Printf("THE STATUS IS %t %t for %s\n", btn.isLocking, btn.isLocked, btn.label)
		} else if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			btn.buttonState = 1
			//fmt.Printf("THE STATUS IS %t %t for %s\n", btn.isLocking, btn.isLocked, btn.label)
		} else {
			if btn.isLocked {
				btn.buttonState = 1
			} else {
				btn.buttonState = 0
			}
		}
	} else {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			btn.buttonState = 2
			//fmt.Printf("THE STATUS IS %d for %s\n", btn.buttonState, btn.label)
		} else if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			btn.buttonState = 1
		} else {
			btn.buttonState = 0
		}
	}
}

func (btn *Button) GetCurrState() int {
	var temp int
	if btn.isLocking {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			temp = 2
			btn.isLocked = !btn.isLocked
			//fmt.Printf("THE STATUS IS %t %t for %s\n", btn.isLocking, btn.isLocked, btn.label)
		} else if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			temp = 1
			//fmt.Printf("THE STATUS IS %t %t for %s\n", btn.isLocking, btn.isLocked, btn.label)
		} else {
			if btn.isLocked {
				btn.buttonState = 1
			} else {
				btn.buttonState = 0
			}
		}
	} else {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			temp = 2
			//fmt.Printf("THE STATUS IS %d for %s\n", btn.buttonState, btn.label)
		} else if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
			temp = 1
		} else {
			temp = 0
		}
	}
	return temp
}

func (btn *Button) Draw(screen *ebiten.Image) {
	//w, h := btn.Simg[btn.buttonState].Bounds().Dx(), btn.Simg[btn.buttonState].Bounds().Dy()
	var op ebiten.DrawImageOptions
	op.GeoM.Reset()
	//scaleX, scaleY := 2.0, 2.0
	//g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	//g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	op.GeoM.Scale(float64(btn.scaleX), float64(btn.scaleY))
	op.GeoM.Translate(float64(btn.bX), float64(btn.bY))
	screen.DrawImage(&btn.Simg[btn.buttonState], &op)
	vector.StrokeLine(screen, float32(btn.bX), float32(btn.bY), float32(btn.bX+btn.bWidth), float32(btn.bY), 2.0, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, float32(btn.bX+btn.bWidth), float32(btn.bY), float32(btn.bX+btn.bWidth), float32(btn.bY+btn.bHeight), 2.0, color.RGBA{255, 0, 0, 255}, true)
	op.GeoM.Reset()
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(float64(btn.bX+(fontSize0/2)), float64(btn.bY+(fontSize0/2)))
	op2.ColorScale.ScaleWithColor(color.RGBA{250, 250, 250, 255})
	op2.LineSpacing = fontSize0
	btn.msg = fmt.Sprintf("%s %d\n%d %t %t", btn.label, len(btn.Simg), btn.buttonState, btn.isMouseOverPos(), ebiten.IsMouseButtonPressed(ebiten.MouseButton0))
	text.Draw(screen, btn.msg, textFaceMono, op2)
}

func (btn *Button) DrawBtnPnl(pnl *ButtonPanel, num int) {
	//w, h := btn.Simg[btn.buttonState].Bounds().Dx(), btn.Simg[btn.buttonState].Bounds().Dy()
	var op ebiten.DrawImageOptions
	op.GeoM.Reset()
	//scaleX, scaleY := 2.0, 2.0
	//g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	//g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	newX, newY := btn.bX-pnl.panelLocX, btn.bY-pnl.panelLocY
	op.GeoM.Scale(float64(btn.scaleX), float64(btn.scaleY))
	op.GeoM.Translate(float64(newX), float64(newY))
	// if btn.buttonState < 2 {
	// 	btn.buttonState++
	// } else {
	// 	btn.buttonState = 0
	// }
	//btn.Update()
	pnl.bckGroundImg.DrawImage(&btn.Simg[num], &op)
	vector.StrokeLine(pnl.bckGroundImg, float32(newX), float32(newY), float32(newX+btn.bWidth), float32(newY), 2.0, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(pnl.bckGroundImg, float32(newX+btn.bWidth), float32(newY), float32(newX+btn.bWidth), float32(newY+btn.bHeight), 2.0, color.RGBA{255, 0, 0, 255}, true)
	op.GeoM.Reset()
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(float64(newX+(fontSize0/2)), float64(newY+(fontSize0/2)))
	op2.ColorScale.ScaleWithColor(color.RGBA{250, 250, 250, 255})
	op2.LineSpacing = fontSize0
	btn.msg = fmt.Sprintf("%s %d\n%d %t", btn.label, len(btn.Simg), btn.buttonState, btn.isMouseOverPos())
	text.Draw(pnl.bckGroundImg, btn.msg, textFaceMono, op2)
}

func (btn *Button) init(lbl string, imgs []ebiten.Image, locX int, locY int, w int, h int) {
	btn.label = lbl
	btn.Simg = imgs
	btn.bX = locX
	btn.bY = locY
	btn.bHeight = h
	btn.bWidth = w
	btn.buttonState = 0
	btn.isLocking = false
	btn.scaleX = 2.0
	btn.scaleY = 2.0
	btn.msg = ""
	// btn.isLocked
}

func getNewButton(lbl string, imgs []ebiten.Image, locX int, locY int, w int, h int) *Button {
	var temp Button
	temp.init(lbl, imgs, locX, locY, w, h)
	return &temp
}

type ButtonPanel struct {
	panelName    string
	bckGroundImg *ebiten.Image
	panelSizeX   int
	panelSizeY   int
	panelLocX    int
	panelLocY    int
	buttons      []Button
	btnStates    []int
	defColor     color.RGBA
	active       bool
	Visible      bool
}
type ButtonPanelOptions struct {
	btnSpacingVert int
	btnSpacingHorz int
	marginH        int
	marginV        int
	borderColor    color.RGBA
}

func (btnPnl *ButtonPanel) InitDefault(pnlName string, xSize int, ySize int, xLoc int, yLoc int, color0 color.RGBA, isActive bool, isVis bool, btnNames []string, bImgs []ebiten.Image) {
	btnPnl.panelName = pnlName
	btnPnl.panelLocX = xLoc
	btnPnl.panelLocY = yLoc
	btnPnl.panelSizeX = xSize
	btnPnl.panelSizeY = ySize
	btnPnl.bckGroundImg = ebiten.NewImage(xSize, ySize)
	btnPnl.defColor = color0
	btnPnl.bckGroundImg.Fill(color0)
	btnPnl.active = isActive
	btnPnl.Visible = isVis

	btnW, btnH := 64, 32
	btnMarginX := 10
	for i, s := range btnNames {
		btnPnl.buttons = append(btnPnl.buttons, *getNewButton(s, bImgs, btnPnl.panelLocX+btnMarginX, btnPnl.panelLocY+25+((btnH+8)*i), btnW, btnH))
		//btnPnl.buttons[len(btnPnl.buttons)-1].isLocking = true
		btnPnl.btnStates = append(btnPnl.btnStates, 0)
		// btnPnl.buttons = append(btnPnl.buttons, *getNewButton(s, bImgs, btnPnl.panelLocX+btnMarginX, btnPnl.panelLocY+10+((btnH+8)*i), btnW, btnH))

	}
}
func (btnPnl *ButtonPanel) Draw(screen *ebiten.Image) {
	for i, b := range btnPnl.buttons {
		b.DrawBtnPnl(btnPnl, btnPnl.btnStates[i])
	}
	var op ebiten.DrawImageOptions
	op.GeoM.Reset()
	scaleX, scaleY := 1.0, 1.0
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(btnPnl.panelLocX), float64(btnPnl.panelLocY))

	screen.DrawImage(btnPnl.bckGroundImg, &op)
	op.GeoM.Reset()
	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(float64(btnPnl.panelLocX), float64(btnPnl.panelLocY))
	op2.ColorScale.ScaleWithColor(color.RGBA{250, 250, 250, 255})
	op2.LineSpacing = fontSize0
	tempText := fmt.Sprintf("%s %t %t [%d]\n%t %d", btnPnl.panelName, btnPnl.active, btnPnl.Visible, len(btnPnl.buttons), btnPnl.isMouseOverPos(), btnPnl.btnStates[0])
	text.Draw(screen, tempText, textFaceMono, op2)
	// for _, b := range btnPnl.buttons {
	// 	b.Draw(screen)
	// }
}

func (btnPnl *ButtonPanel) Update() {

	//btnPnl.buttons[0].Update()
	// for i := 0; i < len(btnPnl.buttons); i++ {
	// 	btnPnl.buttons[i].Update()

	// }
	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		if btnPnl.Visible {
			btnPnl.Visible = false
		} else {
			btnPnl.Visible = true
		}

	}
	if btnPnl.active && btnPnl.Visible && btnPnl.isMouseOverPos() {
		for i, b := range btnPnl.buttons {
			//go b.Update()
			btnPnl.btnStates[i] = b.GetCurrState()
		}
	}
}

func (btnPnl *ButtonPanel) isMouseOverPos() bool {
	mX, mY := ebiten.CursorPosition()
	if ((mX < btnPnl.panelLocX+btnPnl.panelSizeX) && (mX > btnPnl.panelLocX)) && ((mY < btnPnl.panelLocY+btnPnl.panelSizeY) && (mY > btnPnl.panelLocY)) {
		return true
	}
	return false
}

// attempt to make a button appear out of thin air
func MakeImageOfButtons(color0 color.Color, color1 color.Color, w int, h int, brdr float32, antiA bool) *ebiten.Image {
	//the idea is to make a series of images;
	var tempFrame01 = *ebiten.NewImage(w, h)
	var tempFrame02 = *ebiten.NewImage(w, h)
	tempFrame02.Fill(color1)
	vector.StrokeLine(&tempFrame01, 0, 0, float32(w), 0, brdr, color0, antiA)
	vector.StrokeLine(&tempFrame01, float32(w), 0, float32(w), float32(h), brdr, color0, antiA)
	vector.StrokeLine(&tempFrame01, float32(w), float32(h), 0, float32(h), brdr, color0, antiA)
	vector.StrokeLine(&tempFrame01, 0, float32(h), 0, 0, brdr, color0, antiA)

	tempFrame02.DrawImage(&tempFrame01, nil)
	return &tempFrame02
}

func MakeImagesOfButtons(color0 color.Color, color1 color.Color, w int, h int, brdr float32, antiA bool) []ebiten.Image {
	a, b, c, d := color1.RGBA()
	var colorA2 = color.RGBA{uint8(a / 3), uint8(b / 3), uint8(c / 3), uint8(d)}
	var colorA1 = color.RGBA{uint8(a / 2), uint8(b / 2), uint8(c / 2), uint8(d)}
	var colorA0 = color.RGBA{uint8(a), uint8(b), uint8(c), uint8(d)}
	var tempFrames []ebiten.Image
	tempFrames = append(tempFrames, *MakeImageOfButtons(color0, colorA2, w, h, brdr, antiA))
	tempFrames = append(tempFrames, *MakeImageOfButtons(color0, colorA1, w, h, brdr, antiA))
	tempFrames = append(tempFrames, *MakeImageOfButtons(color0, colorA0, w, h, brdr, antiA))
	return tempFrames
}
