package main

import (
	"github.com/jdxyw/generativeart"
	"github.com/jdxyw/generativeart/arts"
	"github.com/jdxyw/generativeart/common"
	"image/color"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	c := generativeart.NewCanva(600, 400)
	/*
		c.SetColorSchema([]color.RGBA{
			{0xCF, 0x2B, 0x34, 0xFF},
			{0xF0, 0x8F, 0x46, 0xFF},
			{0xF0, 0xC1, 0x29, 0xFF},
			{0x19, 0x6E, 0x94, 0xFF},
			{0x35, 0x3A, 0x57, 0xFF},
		})
	*/
	c.SetColorSchema([]color.RGBA{
		common.White,
		common.Tomato,
		common.Azure,
		common.Mintcream,
	})
	c.SetBackground(common.NavajoWhite)
	c.FillBackground()
	c.SetLineWidth(1.0)
	c.SetLineColor(common.Orange)
	c.Draw(arts.NewColorCircle2(30))
	c.ToPNG("circle.png")
}
