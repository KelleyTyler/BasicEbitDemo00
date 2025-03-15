/*
The Purpose of this file is to offload many of the functions and structures that were scattered throughout the original "main.go" file;
*/
package main

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	//fmt
)

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
