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
type AnimatedSprite struct {
	name       string
	imgs       []ebiten.Image
	state      int
	stateLimit int
	timing     []int //the timing per frame; in milisecs
	lastTime   time.Time
	active     bool
}

func (ani *AnimatedSprite) GetCurrFrame() *ebiten.Image {
	return &ani.imgs[ani.state]
}
func (ani *AnimatedSprite) init(title string, imgsIn []ebiten.Image, start int, end int, timingNums []int) *AnimatedSprite {
	temp := AnimatedSprite{
		name:   title,
		imgs:   nil,
		state:  0,
		timing: timingNums,
		active: false,
		//lastTime: nil,
	}
	temp.imgs = GetArrayOfImagesFromArray(imgsIn, start, end)
	return &temp
}
func (ani *AnimatedSprite) toString() string {
	return fmt.Sprintf("ANIMATION:%10s currentFrame/total: %3d/%3d\nSTATUS?:%t", ani.name, ani.state, len(ani.timing), ani.active)
}
func (ani *AnimatedSprite) GetCurrFrameSize() (int, int) {
	var a, b int = 0, 0
	a = ani.imgs[ani.state].Bounds().Dx()
	b = ani.imgs[ani.state].Bounds().Dy()
	return a, b
}

func (ani *AnimatedSprite) Update() {
	tems := time.Since(ani.lastTime)
	if int64(tems.Milliseconds()) > int64(ani.timing[ani.state]) {
		ani.lastTime = time.Now()
		if ani.state < ani.stateLimit {
			ani.state++
		} else {
			ani.state = 0
		}
	}
}

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
func (sprt *Sprite) drawOutline(screen *ebiten.Image) {

	// var path vector.Path
	//var
	var adval float64 = 17.0
	// var vexes []ebiten.Vertex
	// var indes []uint16

	point0x, point0y := float32(sprt.fpX+adval), float32(sprt.fpY+adval)
	point1x, point1y := float32(sprt.fpX+adval), float32(sprt.fpY-adval)
	point2x, point2y := float32(sprt.fpX-adval), float32(sprt.fpY-adval)
	point3x, point3y := float32(sprt.fpX-adval), float32(sprt.fpY+adval)
	//rotation;
	point0x, point0y = rotate(float64(point0x), float64(point0y), sprt.fpX, sprt.fpY, sprt.angle)
	point1x, point1y = rotate(float64(point1x), float64(point1y), sprt.fpX, sprt.fpY, sprt.angle)
	point2x, point2y = rotate(float64(point2x), float64(point2y), sprt.fpX, sprt.fpY, sprt.angle)
	point3x, point3y = rotate(float64(point3x), float64(point3y), sprt.fpX, sprt.fpY, sprt.angle)
	//used these for debugging
	// vector.DrawFilledCircle(screen, point0x, point0y, 5, color.RGBA{0, 0, 255, 255}, false)
	// vector.DrawFilledCircle(screen, point1x, point1y, 5, color.RGBA{0, 255, 255, 255}, false)
	// vector.DrawFilledCircle(screen, point2x, point2y, 5, color.RGBA{255, 0, 255, 255}, false)
	// vector.DrawFilledCircle(screen, point3x, point3y, 5, color.RGBA{0, 255, 0, 255}, false)
	sprt.msg += fmt.Sprintf("POINT 0: %5.2f,%5.2f: %5.2f %5.2f\n", point0x, point0y, point0x-float32(sprt.fpX), point0y-float32(sprt.fpY))
	sprt.msg += fmt.Sprintf("POINT 1: %5.2f,%5.2f: %5.2f %5.2f\n", point1x, point1y, point1x-float32(sprt.fpX), point1y-float32(sprt.fpY))
	vector.StrokeLine(screen, point0x, point0y, point1x, point1y, 1.5, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, point1x, point1y, point2x, point2y, 1.5, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, point2x, point2y, point3x, point3y, 1.5, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, point3x, point3y, point0x, point0y, 1.5, color.RGBA{255, 0, 0, 255}, true)
	// path.MoveTo(point0x, point0y)
	// path.LineTo(point0x, point0y)
	// path.LineTo(point1x, point1y)
	// path.LineTo(point2x, point2y)
	// path.LineTo(point3x, point3y)
	// path.LineTo(point0x, point0y)
	// //fmt.Printf("VS VS VS VS %d %4.5f,%4.5f\n", i, vs[i].DstX, vs[i].DstY)
	// op := &vector.StrokeOptions{}
	// op.LineCap = vector.LineCapSquare
	// op.LineJoin = vector.LineJoinMiter
	// op.MiterLimit = 4
	// op.Width = float32(1.5)
	// vs, is := path.AppendVerticesAndIndicesForStroke(vexes[:0], indes[:0], op)
	// for i := range vs {
	// 	vs[i].SrcX = 1
	// 	vs[i].SrcY = 1
	// 	//vs[i].DstX = vs[i].DstX * 2
	// 	//vs[i].DstY = vs[i].DstY * 2
	// 	vs[i].ColorR = 1.00
	// 	vs[i].ColorG = 0.12
	// 	vs[i].ColorB = 0.12
	// 	vs[i].ColorA = 1

	// }

	// screen.DrawTriangles(vs, is, wSubImage, &ebiten.DrawTrianglesOptions{AntiAlias: true})
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
func GetArrayOfImagesFromArray(imgs []ebiten.Image, start int, end int) []ebiten.Image {
	var temp []ebiten.Image
	for i := start; i < end; i++ {
		temp = append(temp, imgs[i])
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

/*
sprites are...

TODO:thoughts; there needs to be a way to control these things without neededing to necessarily control them manually; there might be some kind of solution to this problem;
but I'm unsure of how to go about it entirely;
Perhaps abstracting it to a layer of interfaces might do the trick;

either way I'm interested in working on the animation arc right now;
*/
type Sprite struct {
	Simg          []ebiten.Image //sprite image; to be replaced by an array
	backupImg     ebiten.Image
	animars       AnimatedSprite
	fpX, fpY      float64 //this is the position in x and y;
	vfX, vfY      float64
	imgHeight     int //
	imgWidth      int //
	showSimg      bool
	angle         float64 //the angle of the image
	imgArrCurrent int
	IsCentered    bool
	acceleration  float64
	maxSpeed      float64
	msg           string
	RotEnabled    bool
}

func (btn *Button) isMouseOverPos() bool {
	mX, mY := ebiten.CursorPosition()
	if ((mX < btn.bX+btn.bWidth) && (mX > btn.bX)) && ((mY < btn.bY+btn.bHeight) && (mY > btn.bY)) {
		return true
	}
	return false
}

/*
Idea: is that there needs to be a way to test if a sprite is in the boundries of another sprite regardless if there's an angle
so the first test will look for anything that is within a "big circle" taking the distance from the center of the rectangel to the corners;
the second test will be more in depth taking into account the specific orientation of the rectangle/collision box;
though perhaps it would be simpler to simply make the "collision boxes" more of a static thing;

then again then again I'm thinking chipmunk might need to come in handy here; but I'm going to wait until after I've implemented sound
*/
// func is_point_inBoxPos() {

// }

/*Sprite.Move:
*this will need a
 */
func (sprt *Sprite) Move(dir int) {
	var deltaX, deltaY float32 = 0.0, 0.0
	switch dir {
	case 0: //up/north
		deltaX, deltaY = rotate(0, -sprt.acceleration, 0, 0, sprt.angle)

		//break
	case 1: //right/east
		deltaX, deltaY = rotate(sprt.acceleration, 0, 0, 0, sprt.angle)

		//break
	case 2: //down/south
		deltaX, deltaY = rotate(0, sprt.acceleration, 0, 0, sprt.angle)

		//break
	case 3: //left/west
		deltaX, deltaY = rotate(-sprt.acceleration, 0, 0, 0, sprt.angle)

		//break
	default: //sloowing down might make this error checking
		// deltaX, deltaY = rotate(0, 0, 0, 0, sprt.angle)
	}
	if math.Abs(sprt.vfX) < sprt.maxSpeed {
		sprt.vfX += float64(deltaX)
	}
	if math.Abs(sprt.vfY) < sprt.maxSpeed {
		sprt.vfY += float64(deltaY)
	}

}
func (sprt *Sprite) Update() {
	//sprt.pX += sprt.vX
	//sprt.pY += sprt.vY
	//sprt.fpX += sprt.vfX
	if math.Abs(sprt.vfX) >= 0.01 {
		sprt.fpX += sprt.vfX
	} else {
		sprt.vfX = 0.0
	}
	if math.Abs(sprt.vfY) >= 0.01 {
		sprt.fpY += sprt.vfY
	} else {
		sprt.vfY = 0.0
	}
	sprt.MovementDrag(0.100000, 0.20000000)
	sprt.MovementDrag(0.02, 0.010) //for some reason this works really well with the speed set at
	var offsetX float64 = 0
	var offsetY float64 = 0
	if sprt.IsCentered {
		offsetX = float64(sprt.imgWidth) / 2.0
		offsetY = float64(sprt.imgHeight) / 2.0
	}
	if sprt.fpX < float64(0)+offsetX {
		// sprt.pX = -sprt.pX
		//sprt.pX = 0
		//sprt.vX = 0
		sprt.fpX = 0.0 + offsetX
		sprt.vfX = 0
		sprt.vfY = 0
	} else if mx := defScrnResX - sprt.imgWidth + int(offsetX); float64(mx) <= sprt.fpX {
		//sprt.pX = 2*mx - sprt.pX
		//sprt.vX = 0
		sprt.vfX = 0.0
		sprt.vfY = 0
		if sprt.IsCentered {
			sprt.fpX = 2*float64(mx) - sprt.fpX
		} else {
			sprt.fpX = 2*float64(mx) - sprt.fpX - offsetX
		}
		//sprt.fpX = 2*float64(mx) - sprt.fpX - offsetX
	}
	if sprt.fpY < float64(0)+offsetY {
		// sprt.pY = -sprt.pY
		//sprt.pY = 0

		//sprt.vY = 0
		sprt.fpY = 0.0 + offsetY
		sprt.vfY = 0
		sprt.vfX = 0
	} else if my := defScrnResY - sprt.imgHeight + int(offsetY); float64(my) <= sprt.fpY {
		//sprt.pY = 2*my - sprt.pY
		//sprt.vY = 0
		sprt.vfY = 0
		sprt.vfX = 0
		if sprt.IsCentered {
			sprt.fpY = 2*float64(my) - sprt.fpY
		} else {
			sprt.fpY = 2*float64(my) - sprt.fpY - offsetY
		}
		//fmt.Printf("MY: %d %f\n %f %f %f\n ", my, float64(my), sprt.fpY, offsetY, float64(2*float64(my)-sprt.fpY-offsetY))
	}
	if sprt.angle > maxAngle || sprt.angle < -(maxAngle) {
		sprt.angle = 0.0
	}
	sprt.animars.Update()
	sprt.msg = fmt.Sprintf("SPRITE:\n (pX,pY):%-06.2f,%-06.2f\n (vX,vY):%-06.2f,%-06.2f\n Angle: %-06.2f RotEn: %t\n", sprt.fpX, sprt.fpY, sprt.vfX, sprt.vfY, sprt.angle, sprt.RotEnabled)
	sprt.msg += fmt.Sprintf(" im(h,w) %3d,%3d STATE:%d %t %d\n ", sprt.imgHeight, sprt.imgWidth, sprt.animars.state, sprt.animars.active, time.Since(sprt.animars.lastTime).Milliseconds())
	sprt.MovementDrag(0.25, 0.5)
	// sprt.vfX = 0.0
	// sprt.vfY = 0.0
}

func (sprt *Sprite) MovementDrag(rateOfDecay float32, cuttoff float64) {
	if sprt.vfX < -cuttoff {
		sprt.vfX += float64(rateOfDecay)
	} else if sprt.vfX > cuttoff {
		sprt.vfX -= float64(rateOfDecay)
	} else {
		//sprt.vfX = 0.0
	}
	if sprt.vfY < -cuttoff {
		sprt.vfY += float64(rateOfDecay)
	} else if sprt.vfY > cuttoff {
		sprt.vfY -= float64(rateOfDecay)
	} else {
		//sprt.vfY = 0.0
	}
}
func (sprt *Sprite) Draw(screen *ebiten.Image, g *Game) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	var offsetX, offsetY float64 = 0.0, 0.0
	if sprt.IsCentered {
		offsetX = -float64(sprt.imgWidth) / 2
		offsetY = -float64(sprt.imgHeight) / 2
	}
	// w, h := 320, 240
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	g.op.GeoM.Rotate(2 * math.Pi * float64(sprt.angle) / float64(maxAngle))
	g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	g.op.GeoM.Translate(float64(sprt.fpX)+offsetX, float64(sprt.fpY)+offsetY)
	//g.op.GeoM.Translate(float64(sprt.pX), float64(sprt.pY))
	//screen.DrawImage(&sprt.backupImg, &g.op)
	if sprt.showSimg {
		screen.DrawImage(&sprt.Simg[sprt.imgArrCurrent], &g.op)
	}

	g.op.GeoM.Reset()
}
func (sprt *Sprite) DrawImageCentered(screen *ebiten.Image, g *Game, adjustX int, adjustY int) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	g.op.GeoM.Rotate(2 * math.Pi * float64(sprt.angle) / float64(maxAngle))
	// g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	if !sprt.IsCentered {
		g.op.GeoM.Translate(float64(w)/2, float64(h)/2)
	}
	g.op.GeoM.Translate(float64(sprt.fpX)+float64(adjustX), float64(sprt.fpY)+float64(adjustY))
	if sprt.showSimg {
		screen.DrawImage(sprt.animars.GetCurrFrame(), &g.op)
		//screen.DrawImage(&sprt.Simg[sprt.imgArrCurrent], &g.op)
	}
	//screen.DrawImage(&sprt.Simg[sprt.imgArrCurrent], &g.op)

	//vector.DrawFilledRect(screen, float32(sprt.fpX-16), float32(sprt.fpY-16), 32, 32, color.RGBA{0, 0, 200, 255}, false)
	vector.DrawFilledCircle(screen, float32(sprt.fpX), float32(sprt.fpY), 5, color.RGBA{250, 100, 100, 255}, true)
	g.op.GeoM.Reset()
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

// func (g *Game) FPSChanger() {
// 	tempTime := time.Since(g.lTime)
// 	if tempTime.Milliseconds() > 250 {
// 		g.fRate = (float64(g.frames) / tempTime.Seconds())
// 		g.frames = 0
// 		g.lTime = time.Now()
// 	} else {

// 		g.frames += 1
// 	}

// }

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
