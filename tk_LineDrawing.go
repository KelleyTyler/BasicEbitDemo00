package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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
