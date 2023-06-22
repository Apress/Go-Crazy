package main

import (
	"fmt"
	"github.com/jdxyw/generativeart"
	"github.com/jdxyw/generativeart/arts"
	"github.com/jdxyw/generativeart/common"
	"image/color"
	"math/rand"
	"time"
)

var DRAWINGS = map[string]generativeart.Engine{
	"maze":      arts.NewMaze(10),
	"julia":     arts.NewJulia(func(z complex128) complex128 { return z*z + complex(-0.1, 0.651) }, 40, 1.5, 1.5),
	"randcicle": arts.NewRandCicle(30, 80, 0.2, 2, 10, 30, true),
	"blackhole": arts.NewBlackHole(200, 400, 0.01),
	"janus":     arts.NewJanus(5, 10),
	"random":    arts.NewRandomShape(150),
	"silksky":   arts.NewSilkSky(15, 5),
	"circles":   arts.NewColorCircle2(30),
}

func main() {
	drawMany(DRAWINGS)
}

func drawMany(drawings map[string]generativeart.Engine) {

	for k, _ := range drawings {
		drawOne(k)
	}

}

func drawOne(art string) string {
	rand.Seed(time.Now().Unix())
	c := generativeart.NewCanva(600, 400)
	c.SetColorSchema([]color.RGBA{
		{0xCF, 0x2B, 0x34, 0xFF},
		{0xF0, 0x8F, 0x46, 0xFF},
		{0xF0, 0xC1, 0x29, 0xFF},
		{0x19, 0x6E, 0x94, 0xFF},
		{0x35, 0x3A, 0x57, 0xFF},
	})
	c.SetBackground(common.NavajoWhite)
	c.FillBackground()
	c.SetLineWidth(1.0)
	c.SetLineColor(common.Orange)
	c.Draw(DRAWINGS[art])

	fileName := fmt.Sprintf("/tmp/%s_%d.png", art, rand.Float64())
	c.ToPNG(fileName)
	return fileName
}
