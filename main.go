package main

import (
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
)

type FrameRateTracking struct {
	fRate  float64
	lTime  time.Time
	frames int
}

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
)

var img *ebiten.Image

var imgs []ebiten.Image

var whiteImage = ebiten.NewImage(defScrnResX, defScrnResY)

func init() {
	whiteImage.Fill(color.RGBA{150, 90, 90, 255})
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
	// var tempRect = image.Rect(0, 0, 32, 32)
	// imgs = append(imgs, *RelativeCrop(temp, tempRect))
	// tempRect = image.Rect(16, 0, 32, 32)
	// imgs = append(imgs, *RelativeCrop(temp, tempRect))
	// tempRect = image.Rect(32, 0, 32, 32)
	// imgs = append(imgs, *RelativeCrop(temp, tempRect))
	// tempRect = image.Rect(64, 0, 32, 32)
	// imgs = append(imgs, *RelativeCrop(temp, tempRect))
	imgs = GetArrayOfImages(temp, 32, 32, 7)
}

func GetArrayOfImages(source *ebiten.Image, subImageX int, subImageY int, numImages int) []ebiten.Image {
	var temp []ebiten.Image
	//var tempRect = image.Rect(0, 0, 32, 32)
	for i := 0; i < numImages; i++ {
		t := (i * subImageX)
		tempRect := image.Rect(t, 0, 32, 32)
		rx, ry := source.Bounds().Min.X+tempRect.Min.X, source.Bounds().Min.Y+tempRect.Min.Y
		fmt.Printf("%2d|%2d|MIN %2d,%2d|MAX %2d,%2d| R: %3d %3d ", i, 0+t, tempRect.Min.X, tempRect.Min.Y, tempRect.Max.X, tempRect.Max.Y, rx, ry)
		fmt.Printf("| SBounds: MIN: %3d %3d MAX: %3d %3d", source.Bounds().Min.X, source.Bounds().Min.Y, source.Bounds().Max.X, source.Bounds().Max.Y)
		//temp2 := RelativeCrop(source, tempRect)

		//bounds := source.Bounds()
		//width := bounds.Dx()
		cropsize := image.Rect(0, 0, subImageY, subImageY)
		cropsize = cropsize.Add(image.Point{(subImageX * i), 0})

		temp2 := source.SubImage(cropsize)
		temp3 := ebiten.NewImageFromImage(temp2)
		// temp2 := source.SubImage(image.Rect(rx, ry, rx+tempRect.Max.X, ry+tempRect.Max.Y)).(*ebiten.Image)
		//temp3.Fill(color.RGBA{uint8(15), uint8(60), uint8(25), uint8(100)})
		fmt.Printf(" TEMP%d:Dx/Dy: %d %d MAX: %d,%d\n", i, temp2.Bounds().Dx(), temp2.Bounds().Dy(), temp2.Bounds().Max.X, temp2.Bounds().Max.Y)
		temp = append(temp, *temp3)
	}
	return temp
}

func RelativeCrop(source *ebiten.Image, r image.Rectangle) *ebiten.Image {
	rx, ry := source.Bounds().Min.X+r.Min.X, source.Bounds().Min.Y+r.Min.Y
	return source.SubImage(image.Rect(rx, ry, rx+r.Max.X, ry+r.Max.Y)).(*ebiten.Image)
}

type Sprite struct {
	Simg      []ebiten.Image //sprite image; to be replaced by an array
	pX, pY    int            //this is the position in x and y;
	vX, vY    int            //this is the velocity in x and y;
	imgHeight int            //
	imgWidth  int            //

	angle         int //the angle of the image
	imgArrCurrent int
	imgArrDown    bool
}

// func spriteInit(fioPath string) *Sprite {
// 	temp := Sprite{Simg: nil, pX: spriteX, pY: spriteY, vX: 0, vY: 0, imgHeight: 0, imgWidth: 0}
// 	temp.Simg = img // going to want to make img an array in the future; or a map or something;
// 	// so that way there's multiple codes;
// 	return &temp
// }

// func (sprt *Sprite) Init(fIoPath string) error {

//		return nil
//	}
/*Sprite.Move:
*this will need a
 */
func (sprt *Sprite) Update() {
	sprt.pX += sprt.vX
	sprt.pY += sprt.vY
	if sprt.pX < 0 {
		// sprt.pX = -sprt.pX
		sprt.pX = 0
		sprt.vX = 0
	} else if mx := defScrnResX - sprt.imgWidth; mx <= sprt.pX {
		sprt.pX = 2*mx - sprt.pX
		sprt.vX = 0
	}
	if sprt.pY < 0 {
		// sprt.pY = -sprt.pY
		sprt.pY = 0

		sprt.vY = 0
	} else if my := defScrnResY - sprt.imgHeight; my <= sprt.pY {
		sprt.pY = 2*my - sprt.pY
		sprt.vY = 0
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
	g.op.GeoM.Translate(float64(sprt.pX), float64(sprt.pY))
	screen.DrawImage(&sprt.Simg[sprt.imgArrCurrent], &g.op)
	g.op.GeoM.Reset()
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
}

func (g *Game) init() error {
	defer func() {
		g.inited = true
	}()

	g.sprt = Sprite{
		Simg:          imgs,
		pX:            64,
		pY:            64,
		vX:            0,
		vY:            0,
		imgHeight:     32,
		imgWidth:      32,
		angle:         0,
		imgArrCurrent: 0,
		//imgArrDown:    false,
	}
	//g.sprt.Simg = append(g.sprt.Simg, *img)
	return nil
}
func (g *Game) Update() error {
	if !g.inited {
		g.init()
		fmt.Printf("INITIATED \n")
	}
	//game logic goes here;
	//this might be the basic CPU-type logic only though... not sure;
	g.sprt.Update()
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.sprt.vX += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.sprt.vX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.sprt.vY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.sprt.vY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && !g.sprt.imgArrDown {

		//g.sprt.imgArrDown = true

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
	screen.DrawImage(whiteImage, nil)
	g.FPSChanger()
	//screen.DrawImage(img, nil)
	// screen.DrawImage(g.sprt.Simg, &g.op)
	g.sprt.Draw(screen, g)

	//this might be graphic layout?? a means perhaps to control layers?
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%3.1f\nA:%3.1f\nFRAMES:%d", g.fRate, g.fRateAvg, g.frames))
	g.gMSG = fmt.Sprintf("FPS:%3.1f\nSPRITE:\n(pX,pY):%3d,%3d\n(vX,Vy):%3d,%3d\nImg(W,H):%3d,%3d\nAngle:%3d\nIMG:%3d", g.fRate, g.sprt.pX, g.sprt.pY, g.sprt.vX, g.sprt.vY, g.sprt.imgWidth, g.sprt.imgHeight, g.sprt.angle, g.sprt.imgArrCurrent)
	ebitenutil.DebugPrint(screen, g.gMSG)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return defScrnResX, defScrnResY
}

func main() {
	//initializing
	// var g = Game{
	// 	fRate:  0.0,
	// 	lTime:  time.Now(),
	// 	frames: 0,
	// }
	ebiten.SetWindowSize(defWindowWidth, defWindowHeight)
	// presentTime := time.Now()
	ebiten.SetWindowTitle("EBITEN TEST!")
	//this is where the game logic runs
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
