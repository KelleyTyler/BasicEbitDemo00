package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
)

// type FrameRateTracking struct {
// 	fRate  float64
// 	lTime  time.Time
// 	frames int
// }

// type IntVector2D struct {
// 	x int
// 	y int
// }

// func (vec0 *IntVector2D) VectorAdditionI(vec1 IntVector2D) {
// 	vec0.x = vec0.x + vec1.x
// 	vec0.y = vec0.x + vec1.y
// }

const (
	spriteX         int = 64
	spriteY         int = 64
	maxAngle        int = 360
	defScrnResX         = 320 //320
	defScrnResY         = 240 //240
	defWindowWidth      = 640
	defWindowHeight     = 480
	fontSize0           = 10
)

var img *ebiten.Image

var imgs []ebiten.Image
var btnImgs []ebiten.Image
var btnImgs1 []ebiten.Image
var backgroundColor = ebiten.NewImage(defScrnResX, defScrnResY)
var faceSrc *text.GoTextFaceSource
var textface text.Face

func init() {
	backgroundColor.Fill(color.RGBA{200, 200, 200, 255})
	var err error
	img, _, err = ebitenutil.NewImageFromFile("assets/Square_32x32Texture.png")
	if err != nil {
		log.Fatal(err)
	}
	var temp *ebiten.Image = nil

	temp, _, err = ebitenutil.NewImageFromFile("assets/256x64Texture.png")
	if err != nil {
		log.Fatal(err)
	}
	imgs = GetArrayOfImages(temp, 0, 0, 32, 0, 32, 0, 8)
	temp, _, err = ebitenutil.NewImageFromFile("assets/96x32Buttons.png")
	if err != nil {
		log.Fatal(err)
	}
	btnImgs = GetArrayOfImages(temp, 0, 0, 32, 0, 16, 0, 3)
	btnImgs1 = GetArrayOfImages(temp, 0, 1, 32, 0, 16, 0, 3)
	faceSrc, err = text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
	if err != nil {
		log.Fatal("err: ", err)
	}
	textface = &text.GoTextFace{
		Source: faceSrc,
		Size:   fontSize0,
	}

}

/*
quite proud of this function; was an improvement on the help I'd seen online;
this is a function that can do a lot;
ISSUES/TODO: error checking before I move it to a more modular location;
*/
func GetArrayOfImages(source *ebiten.Image, skipTilesX int, skipTilesY int, subImageX int, xBuf int, subImageY int, yBuf int, numImages int) []ebiten.Image {
	var temp []ebiten.Image
	//the number we skip to;
	a, b := 0, 0

	if (subImageX * skipTilesX) > (source.Bounds().Max.X) {
		//find out by how much..
		e := source.Bounds().Max.X / subImageX
		f := skipTilesX - e
		//fmt.Printf("OVERFLOW %d %d\n", e, f)
		b++
		a = f
	} else {
		a = skipTilesX
	}
	b = skipTilesY
	for i := 0; i < numImages; i++ {
		if (a * subImageX) >= source.Bounds().Max.X {
			b++
			a = 0
		}
		//fmt.Printf("| SBounds: MIN: %3d %3d MAX: %3d %3d", source.Bounds().Min.X, source.Bounds().Min.Y, source.Bounds().Max.X, source.Bounds().Max.Y)
		cropsize := image.Rect(0, 0, subImageX, subImageY)
		cropsize = cropsize.Add(image.Point{(subImageX * a) + xBuf, (subImageY * b) + yBuf})
		temp2 := source.SubImage(cropsize)
		temp3 := ebiten.NewImageFromImage(temp2)
		//fmt.Printf(" TEMP%d:Dx/Dy: %d %d MAX: %d,%d\n", i, temp2.Bounds().Dx(), temp2.Bounds().Dy(), temp2.Bounds().Max.X, temp2.Bounds().Max.Y)
		temp = append(temp, *temp3)
		a++
	}
	return temp
}

type Button struct {
	Simg            []ebiten.Image
	buttonState     int
	bX, bY          int
	bHeight, bWidth int
	label           string
}

type Sprite struct {
	Simg      []ebiten.Image //sprite image; to be replaced by an array
	pX, pY    int            //this is the position in x and y;
	fpX, fpY  float64
	vfX, vfY  float64
	vX, vY    int //this is the velocity in x and y;
	imgHeight int //
	imgWidth  int //

	angle         int //the angle of the image
	imgArrCurrent int
}

func (btn *Button) isMouseOverPos() bool {
	mX, mY := ebiten.CursorPosition()
	if ((mX < btn.bX+btn.bWidth) && (mX > btn.bX)) && ((mY < btn.bY+btn.bHeight) && (mY > btn.bY)) {
		return true
	}
	return false
}

/*Sprite.Move:
*this will need a
 */
func (sprt *Sprite) Update() {
	//sprt.pX += sprt.vX
	//sprt.pY += sprt.vY
	sprt.fpX += sprt.vfX
	sprt.fpY += sprt.vfY
	if sprt.fpX < float64(0) {
		// sprt.pX = -sprt.pX
		//sprt.pX = 0
		//sprt.vX = 0
		sprt.fpX = 0.0
		sprt.vfX = 0
	} else if mx := defScrnResX - sprt.imgWidth; float64(mx) <= sprt.fpX {
		//sprt.pX = 2*mx - sprt.pX
		//sprt.vX = 0
		sprt.vfX = 0.0
		sprt.fpX = 2*float64(mx) - sprt.fpX
	}
	if sprt.fpY < float64(0) {
		// sprt.pY = -sprt.pY
		//sprt.pY = 0

		//sprt.vY = 0
		sprt.fpY = 0.0
		sprt.vfY = 0
	} else if my := defScrnResY - sprt.imgHeight; float64(my) <= sprt.fpY {
		//sprt.pY = 2*my - sprt.pY
		//sprt.vY = 0
		sprt.vfY = 0
		sprt.fpY = 2*float64(my) - sprt.fpY
	}
	//sprt.angle++
	if sprt.angle == maxAngle {
		sprt.angle = 0
	}

}

func (sprt *Sprite) Draw(screen *ebiten.Image, g *Game) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	// w, h := 320, 240
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	g.op.GeoM.Rotate(2 * math.Pi * float64(sprt.angle) / float64(maxAngle))
	g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	g.op.GeoM.Translate(float64(sprt.fpX), float64(sprt.fpY))
	//g.op.GeoM.Translate(float64(sprt.pX), float64(sprt.pY))
	screen.DrawImage(&sprt.Simg[sprt.imgArrCurrent], &g.op)
	g.op.GeoM.Reset()
}

func (btn *Button) Draw(screen *ebiten.Image, g *Game) {
	//w, h := btn.Simg[btn.buttonState].Bounds().Dx(), btn.Simg[btn.buttonState].Bounds().Dy()
	g.op.GeoM.Reset()
	//g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	//g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	g.op.GeoM.Scale(2.0, 2.0)
	g.op.GeoM.Translate(float64(btn.bX), float64(btn.bY))
	screen.DrawImage(&btn.Simg[btn.buttonState], &g.op)
	g.op.GeoM.Reset()
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(btn.bX+(fontSize0/2)), float64(btn.bY+(fontSize0/2)))
	op.ColorScale.ScaleWithColor(color.RGBA{250, 250, 250, 255})
	op.LineSpacing = fontSize0
	text.Draw(screen, btn.label, textface, op)
}
func (btn *Button) Update(g *Game) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
		btn.buttonState = 2
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && btn.isMouseOverPos() {
		btn.buttonState = 1
	} else {
		btn.buttonState = 0
	}
}

type Game struct {
	//game variables and other such things go here;
	fRate  float64
	lTime  time.Time
	frames int
	inited bool
	//enough of this
	sprt Sprite
	op   ebiten.DrawImageOptions
	gMSG string
	btn0 Button
	btn1 Button
}

func (g *Game) init() error {
	defer func() {
		g.inited = true
	}()

	g.sprt = Sprite{
		Simg:          imgs,
		pX:            64,
		pY:            64,
		fpX:           64.0,
		fpY:           64.0,
		vX:            0,
		vY:            0,
		vfY:           0.0,
		vfX:           0.0,
		imgHeight:     32,
		imgWidth:      32,
		angle:         0,
		imgArrCurrent: 0,
		//imgArrDown:    false,
	}
	//g.sprt.Simg = append(g.sprt.Simg, *img)
	g.btn0 = Button{
		Simg:        btnImgs,
		bX:          defScrnResX - 64,
		bY:          32,
		bHeight:     32,
		bWidth:      64,
		buttonState: 0,
		label:       "Btn 0",
	}
	g.btn1 = Button{
		Simg:        btnImgs1,
		bX:          defScrnResX - 64,
		bY:          64,
		bHeight:     32,
		bWidth:      64,
		buttonState: 0,
		label:       "Btn 1",
	}
	return nil
}
func (g *Game) Update() error {
	if !g.inited {
		g.init()
		//fmt.Printf("INITIATED \n")
	}
	//game logic goes here;
	//this might be the basic CPU-type logic only though... not sure;
	g.sprt.Update()
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		//g.sprt.vX += 1
		g.sprt.vfX += 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		//g.sprt.vX -= 1
		g.sprt.vfX -= 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		//g.sprt.vY += 1
		g.sprt.vfY += 0.5
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		//g.sprt.vY -= 1
		g.sprt.vfY -= 0.5
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		if g.sprt.imgArrCurrent < (len(g.sprt.Simg) - 1) {
			g.sprt.imgArrCurrent += 1
		} else {
			g.sprt.imgArrCurrent = 0
		}
		g.sprt.imgHeight = g.sprt.Simg[g.sprt.imgArrCurrent].Bounds().Max.Y
		g.sprt.imgWidth = g.sprt.Simg[g.sprt.imgArrCurrent].Bounds().Max.X
	}
	// if inpututil.IsKeyJustReleased(ebiten.KeyD) {
	// 	if g.btn0.buttonState < (len(g.btn0.Simg) - 1) {
	// 		g.btn0.buttonState += 1
	// 	} else {
	// 		g.btn0.buttonState = 0
	// 	}
	// 	//g.sprt.imgHeight = g.sprt.Simg[g.sprt.imgArrCurrent].Bounds().Max.Y
	// 	//g.sprt.imgWidth = g.sprt.Simg[g.sprt.imgArrCurrent].Bounds().Max.X
	// }

	// if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && g.btn0.isMouseOverPos() {
	// 	g.btn0.buttonState = 2
	// 	// if g.btn0.buttonState < (len(g.btn0.Simg) - 1) {

	// 	// } else {

	// 	// }
	// } else if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && g.btn0.isMouseOverPos() {
	// 	g.btn0.buttonState = 1
	// } else {
	// 	g.btn0.buttonState = 0
	// }
	g.btn0.Update(g)
	g.btn1.Update(g)
	if g.btn0.buttonState == 2 {
		backgroundColor.Fill(color.RGBA{150, 150, 200, 255})
	} else if g.btn1.buttonState == 2 {
		backgroundColor.Fill(color.RGBA{200, 150, 150, 255})
	} else {
		backgroundColor.Fill(color.RGBA{200, 200, 200, 255})
	}

	return nil
}
func (g *Game) FPSChanger() {
	tempTime := time.Since(g.lTime)
	if tempTime.Milliseconds() > 250 {
		g.fRate = (float64(g.frames) / tempTime.Seconds())
		g.frames = 0
		g.lTime = time.Now()
	} else {

		g.frames += 1
	}

}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundColor, nil)
	g.FPSChanger()
	g.sprt.Draw(screen, g)
	g.btn0.Draw(screen, g)
	g.btn1.Draw(screen, g)
	g.gMSG = fmt.Sprintf("FPS:%3.1f\nSPRITE:\n(pX,pY):%3f,%3f\n", g.fRate, g.sprt.fpX, g.sprt.fpY)
	g.gMSG += fmt.Sprintf("(vfX,vfY):%3f,%3f\nImg(W,H):%3d,%3d\nAngle:%3d\nIMG:%3d", g.sprt.vfX, g.sprt.vfY, g.sprt.imgWidth, g.sprt.imgHeight, g.sprt.angle, g.sprt.imgArrCurrent)
	ebitenutil.DebugPrint(screen, g.gMSG)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return defScrnResX, defScrnResY
}

func main() {

	ebiten.SetWindowSize(defWindowWidth, defWindowHeight)
	// presentTime := time.Now()
	ebiten.SetWindowTitle("EBITEN TEST!")
	//this is where the game logic runs
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
