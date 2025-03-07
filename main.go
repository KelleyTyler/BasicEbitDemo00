package main

import (
	"fmt"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type FrameRateTracking struct {
	fRate  float64
	lTime  time.Time
	frames int
}

type Game struct {
	//game variables and other such things go here;
	fRate  float64
	lTime  time.Time
	frames int
}

func (g *Game) Update() error {
	//game logic goes here;
	//this might be the basic CPU-type logic only though... not sure;
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
	//this might be graphic layout?? a means perhaps to control layers?
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%3.1f\nA:%3.1f\nFRAMES:%d", g.fRate, g.fRateAvg, g.frames))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS:%3.1f", g.fRate))
	g.FPSChanger()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	//initializing
	var g = Game{
		fRate:  0.0,
		lTime:  time.Now(),
		frames: 0,
	}
	ebiten.SetWindowSize(640, 480)
	// presentTime := time.Now()
	ebiten.SetWindowTitle("Hello, World!")
	//this is where the game logic runs
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
