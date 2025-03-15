package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

/*
Idea: is that there needs to be a way to test if a sprite is in the boundries of another sprite regardless if there's an angle
so the first test will look for anything that is within a "big circle" taking the distance from the center of the rectangel to the corners;
the second test will be more in depth taking into account the specific orientation of the rectangle/collision box;
though perhaps it would be simpler to simply make the "collision boxes" more of a static thing;

then again then again I'm thinking chipmunk might need to come in handy here; but I'm going to wait until after I've implemented sound
*/
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
	screen.DrawImage(&sprt.backupImg, &g.op)
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
	screen.DrawImage(&sprt.backupImg, &g.op)
	//vector.DrawFilledRect(screen, float32(sprt.fpX-16), float32(sprt.fpY-16), 32, 32, color.RGBA{0, 0, 200, 255}, false)
	vector.DrawFilledCircle(screen, float32(sprt.fpX), float32(sprt.fpY), 5, color.RGBA{250, 100, 100, 255}, true)
	g.op.GeoM.Reset()
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
