package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	spriteX         int     = 64
	spriteY         int     = 64
	maxAngle        float64 = 360.00
	defScrnResX             = 320 //320
	defScrnResY             = 240 //240
	defWindowWidth          = 640
	defWindowHeight         = 480
	fontSize0               = 10
)

var (
	img *ebiten.Image

	imgs            []ebiten.Image
	btnImgs         []ebiten.Image
	btnImgs1        []ebiten.Image
	backgroundColor = ebiten.NewImage(defScrnResX, defScrnResY)
	foreground      = ebiten.NewImage(defScrnResX, defScrnResY)
	faceSrc         *text.GoTextFaceSource
	textFaceMono    text.Face
	textFaceReg20   text.Face
	defSpriteImg    = ebiten.NewImage(8, 8)
	wSubImage       = ebiten.NewImage(5, 5)
)

//type imArchive [][]ebiten.Image//the point of this is to store a lot of images; perhaps a record

// type ImArchive struct {
// 	imgs []ebiten.Image

// }

type VectorThing struct {
	vertices       []ebiten.Vertex
	indices        []uint16
	aa, showcenter bool
	drawn          bool
	angle          int //test angle Im thinking
}

func (vex *VectorThing) Update() {

}
func (vex *VectorThing) Draw(screen *ebiten.Image, g *Game) {
	target := screen
	joins := []vector.LineJoin{
		vector.LineJoinMiter,
		vector.LineJoinMiter,
		vector.LineJoinBevel,
		vector.LineJoinRound,
	}
	caps := []vector.LineCap{
		vector.LineCapButt,
		vector.LineCapRound,
		vector.LineCapSquare,
	}
	vex.drawLine(target, image.Rectangle{image.Point{200, 0}, image.Point{100, 240}}, caps[2], joins[3], 1.0)
}

func rotate(px float64, py float64, ox float64, oy float64, angle float64) (float32, float32) {
	// (float64, float64)
	//normalize angle
	angl := angle
	if angl < 0 {
		angl = angl + 360.00
	}

	if angl > 360 {
		angl = angl - 360
	}
	//normalize point;
	ppx := px - ox
	ppy := py - oy
	//convert to radians
	radians := (math.Pi * angl) / (180.00)
	//
	qx := (math.Cos(radians) * ppx) - (math.Sin(radians) * ppy)
	qy := (math.Sin(radians) * ppx) + (math.Cos(radians) * ppy)
	qx = ox + qx
	qy = oy + qy
	return float32(qx), float32(qy)
	// return qx, qy
}

// failed experimental attempt to import;
func (vex *VectorThing) drawLine(screen *ebiten.Image, region image.Rectangle, cap vector.LineCap, join vector.LineJoin, miterLimit float32) {
	c0x := float64(region.Min.X + region.Dx()/4)
	c0y := float64(region.Min.Y + region.Dy()/4)
	c1x := float64(region.Min.X + region.Dx()/4)
	c1y := float64(region.Max.Y + region.Dy()/4)
	c2x := float64(region.Max.X + region.Dx()/4)
	c2y := float64(region.Max.Y + region.Dy()/4)
	c3x := float64(region.Max.X + region.Dx()/4)
	c3y := float64(region.Min.Y + region.Dy()/4)

	var path vector.Path
	path.MoveTo(float32(c0x), float32(c0y))
	path.LineTo(float32(c1x), float32(c1y))
	path.LineTo(float32(c2x), float32(c2y))
	path.LineTo(float32(c3x), float32(c3y))
	path.LineTo(float32(c0x), float32(c0y))
	op := &vector.StrokeOptions{}
	op.LineCap = cap
	op.LineJoin = join
	op.MiterLimit = miterLimit
	op.Width = float32(5)
	vs, is := path.AppendVerticesAndIndicesForStroke(vex.vertices[:0], vex.indices[:0], op)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = 1
		vs[i].ColorG = 0.02
		vs[i].ColorB = 0.02
		vs[i].ColorA = 1

	}
	// vs[0].SrcX = 2
	// vs[0].SrcY =1
	screen.DrawTriangles(vs, is, wSubImage, &ebiten.DrawTrianglesOptions{AntiAlias: vex.aa})

}

func init() {
	backgroundColor.Fill(color.RGBA{150, 150, 150, 255})
	foreground.Fill(color.RGBA{255, 255, 255, 0})
	// defSpriteImg.Fill(color.RGBA{200, 15, 15, 255})
	//tempReader = bytes.NewReader()
	wSubImage.Fill(color.White)
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
	imgs = GetArrayOfImages(temp, 0, 0, 32, 0, 32, 0, 14)
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
	textFaceMono = &text.GoTextFace{
		Source: faceSrc,
		Size:   fontSize0,
	}
	faceSrc, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal("err: ", err)
	}
	textFaceReg20 = &text.GoTextFace{
		Source: faceSrc,
		Size:   20,
	}
}

// func is_point_inBoxPos() {

// }

type Game struct {
	//game variables and other such things go here;
	// fRate  float64
	// lTime  time.Time
	// frames int
	inited bool
	//enough of this
	sprt   Sprite
	op     ebiten.DrawImageOptions
	gMSG   string
	btn0   Button
	btn1   Button
	vectra VectorThing
	//TimerSys *ebiten.
}

func (g *Game) init() error {
	defer func() {
		g.inited = true
	}()

	g.sprt = Sprite{
		Simg:          imgs,
		backupImg:     *defSpriteImg,
		animars:       *g.sprt.animars.init("idle", imgs, 4, 12, []int{500, 500, 700, 500, 300, 500, 200, 500}),
		fpX:           64.0,
		fpY:           64.0,
		vfY:           0.0,
		vfX:           0.0,
		imgHeight:     32,
		imgWidth:      32,
		angle:         0,
		imgArrCurrent: 0,
		showSimg:      true,
		IsCentered:    true,
		acceleration:  0.12,
		maxSpeed:      4.0,
		msg:           "",
		RotEnabled:    true,
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
	g.vectra = VectorThing{
		vertices:   nil,
		indices:    nil,
		aa:         true,
		showcenter: false,
	}
	g.vectra.vertices = append(g.vectra.vertices, ebiten.Vertex{SrcX: 0.0, SrcY: 0.0, DstX: -50.0, DstY: 50.0})
	g.vectra.vertices = append(g.vectra.vertices, ebiten.Vertex{SrcX: 0.0, SrcY: 0.0, DstX: 50.0, DstY: 50.0})
	g.vectra.vertices = append(g.vectra.vertices, ebiten.Vertex{SrcX: 0.0, SrcY: 0.0, DstX: 50.0, DstY: -50.0})
	g.vectra.vertices = append(g.vectra.vertices, ebiten.Vertex{SrcX: 0.0, SrcY: 0.0, DstX: -50.0, DstY: -50.0})
	g.vectra.drawn = false

	g.sprt.animars.active = true
	g.sprt.animars.stateLimit = 6
	//g.sprt.animars.lastTime = time.Now()
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
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		//g.sprt.vX += 1
		//g.sprt.vfX += 0.5
		// if g.sprt.vfX < (g.sprt.maxSpeed) {
		// 	g.sprt.vfX += 0.5
		// }
		g.sprt.Move(1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		//g.sprt.vX -= 1
		//g.sprt.vfX -= 0.5
		// if g.sprt.vfX > -(g.sprt.maxSpeed) {
		// 	g.sprt.vfX -= 0.5
		// }
		g.sprt.Move(3)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		//g.sprt.vY += 1
		//g.sprt.vfY += 0.5
		// if g.sprt.vfY > -(g.sprt.maxSpeed) {
		// 	g.sprt.vfY -= 0.5
		// }
		g.sprt.Move(0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		//g.sprt.vY -= 1
		// if g.sprt.vfY < (g.sprt.maxSpeed) {
		// 	g.sprt.vfY += 0.5
		// }
		g.sprt.Move(2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if g.sprt.RotEnabled {
			g.sprt.angle += 1.0
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		if g.sprt.RotEnabled {
			g.sprt.angle -= 1.0
		}
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyZ) {
		if g.sprt.imgArrCurrent < (len(g.sprt.Simg) - 1) {
			g.sprt.imgArrCurrent += 1
		} else {
			g.sprt.imgArrCurrent = 0
		}
		g.sprt.imgHeight = g.sprt.Simg[g.sprt.imgArrCurrent].Bounds().Max.Y
		g.sprt.imgWidth = g.sprt.Simg[g.sprt.imgArrCurrent].Bounds().Max.X
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyX) {
		if g.sprt.showSimg {
			g.sprt.showSimg = false
		} else {
			g.sprt.showSimg = true
		}

	}
	if inpututil.IsKeyJustReleased(ebiten.KeyC) {
		if g.sprt.IsCentered {
			g.sprt.IsCentered = false
		} else {
			g.sprt.IsCentered = true
		}

	}

	g.btn0.Update(g)
	g.btn1.Update(g)
	if g.btn0.buttonState == 2 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.sprt.angle = 0.0
		//backgroundColor.Fill(color.RGBA{150, 150, 200, 255})
	}
	if g.btn1.buttonState == 2 {
		g.sprt.RotEnabled = !g.sprt.RotEnabled
	}

	return nil
}

/*
Moving all of this to here is a good move IMO;
Keeps it simple
*/
func writeRegText(screen *ebiten.Image, inStrng string, xx int, yy int, textSize int, colored color.RGBA) {
	op := &text.DrawOptions{}
	scaler := 2.0
	// op.GeoM.Translate(float64(xx+(textSize/2)), float64(yy+(textSize/2)))
	op.GeoM.Translate(float64(xx)*scaler, float64(yy)*scaler)

	op.GeoM.Scale(1/scaler, 1/scaler)
	// op.GeoM.Scale(2, 2)

	op.ColorScale.ScaleWithColor(colored)
	//fontSize0
	op.LineSpacing = float64(20)
	outStrng := inStrng + fmt.Sprintf(" %d ", textSize)

	text.Draw(screen, outStrng, textFaceReg20, op)
}
func (g *Game) PreDraw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 0})
	g.gMSG = fmt.Sprintf("FPS:%3.1f %3.1f\n", ebiten.ActualFPS(), ebiten.ActualTPS())
	g.sprt.drawOutline(screen)
	g.sprt.DrawImageCentered(screen, g, 0, 0)

	g.btn0.Draw(screen, g)
	g.btn1.Draw(screen, g)

	// g.vectra.Draw(screen, g)
	// if !g.vectra.drawn {
	// 	g.vectra.drawMyShape(screen)
	// 	g.vectra.drawn = true
	// }

	//g.gMSG = fmt.Sprintf("FPS:%3.1f %3.1f\n", ebiten.ActualFPS(), ebiten.ActualTPS())
	g.gMSG += g.sprt.msg
	g.gMSG += fmt.Sprintf("%t %3.2f %3.2f angl %3d\n", g.vectra.drawn, g.vectra.vertices[0].DstX, g.vectra.vertices[0].DstY, g.vectra.angle)
	ebitenutil.DebugPrint(screen, g.gMSG)
	writeRegText(screen, "Hello\nThis Is A Fish\nNO NOT REALLY", 40, 135, 10, color.RGBA{0, 0, 0, 255})
	//writeRegText(screen, "Hello", 40, 145, 10, color.RGBA{0, 0, 0, 255})
	//writeRegText(screen, "Hello", 40, 155, 10, color.RGBA{0, 0, 0, 255})
	g.gMSG = ""
	g.sprt.msg = ""
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.PreDraw(foreground)
	screen.DrawImage(backgroundColor, nil)
	screen.DrawImage(foreground, nil)
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
