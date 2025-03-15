package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func CreateBasicImg(clor color.Color, w int, h int) *ebiten.Image {
	var temp *ebiten.Image = ebiten.NewImage(w, h)
	temp.Fill(clor)
	return temp
}
