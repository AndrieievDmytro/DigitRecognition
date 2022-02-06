package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"github.com/fogleman/gg"
)

func setCanvas() (gui.Env, *gg.Context, *gui.Mux, *image.RGBA) {
	originalEnv, err := win.New(win.Title("Paint"), win.Size(800, 600))
	if err != nil {
		panic(err)
	}
	mux, env := gui.NewMux(originalEnv)
	r := image.Rect(0, 0, 800, 600)
	canvas := image.NewRGBA(r)
	draw.Draw(canvas, r, image.White, r.Min, draw.Src)

	env.Draw() <- func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, r, canvas, image.Point{}, draw.Src)
		return r
	}
	// Prepares a context for rendering onto the specified image
	drawContext := gg.NewContextForRGBA(canvas)
	return env, drawContext, mux, canvas
}

func handleDrawEvent(originalEnv, env gui.Env, drawContext *gg.Context, canvas *image.RGBA) {
	var (
		clr           = color.Color(color.Black)
		mousePressed  = false
		buttonPressed = false
		prev          image.Point
		bold          = 2
	)
	for event := range env.Events() {
		switch event := event.(type) {

		case win.WiClose:
			close(originalEnv.Draw())

		case win.MoDown:
			if event.Point.In(image.Rect(0, 0, 800, 600)) {
				mousePressed = true
				prev = event.Point
			}

		case win.MoUp:
			mousePressed = false
			gg.SavePNG("number.png", canvas)

		case win.MoMove:
			if mousePressed && !buttonPressed {
				x0, y0, x1, y1 := prev.X, prev.Y, event.X, event.Y
				prev = event.Point

				env.Draw() <- func(drw draw.Image) image.Rectangle {
					clr = color.Black
					drawContext.SetColor(clr)
					drawContext.SetLineCapRound()
					drawContext.SetLineWidth(float64(bold))
					drawContext.DrawLine(float64(x0), float64(y0), float64(x1), float64(y1))
					drawContext.Stroke()
					rect := image.Rect(x0, y0, x1, y1)
					rect.Min.X -= bold
					rect.Min.Y -= bold
					rect.Max.X += bold
					rect.Max.Y += bold
					draw.Draw(drw, rect, canvas, rect.Min, draw.Src)
					return rect
				}
			} else if mousePressed && buttonPressed {
				x0, y0, x1, y1 := prev.X, prev.Y, event.X, event.Y
				prev = event.Point

				env.Draw() <- func(drw draw.Image) image.Rectangle {
					clr = color.White
					drawContext.SetColor(clr)
					drawContext.SetLineCapRound()
					drawContext.SetLineWidth(float64(20))
					drawContext.DrawLine(float64(x0), float64(y0), float64(x1), float64(y1))
					drawContext.Stroke()
					rect := image.Rect(x0, y0, x1, y1)
					rect.Min.X -= 20
					rect.Min.Y -= 20
					rect.Max.X += 20
					rect.Max.Y += 20
					draw.Draw(drw, rect, canvas, rect.Min, draw.Src)
					return rect
				}
			}
		case win.KbDown:
			buttonPressed = true

		case win.KbUp:
			buttonPressed = false
		}
	}
}

func createGUI() {
	originalEnv, drawContext, mux, canvas := setCanvas()
	handleDrawEvent(originalEnv, mux.MakeEnv(), drawContext, canvas)

}
