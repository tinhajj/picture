package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	winWidth, winHeight = 1000, 1000
	sizeEvent           size.Event
	uistate             state
)

type state struct {
	mousex     float64
	mousey     float64
	activeitem int
	mousedown  bool
}

func main() {
	fmt.Println("start")

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Simple Window for Graphics",
			Width:  winWidth,
			Height: winHeight,
		})
		if err != nil {
			panic(err)
		}
		defer w.Release()

		screenSize := image.Point{
			winWidth, winHeight,
		}
		screenBuffer, err := s.NewBuffer(screenSize)
		if err != nil {
			panic(err)
		}
		defer screenBuffer.Release()
		pixBuffer := screenBuffer.RGBA()
		_ = pixBuffer

		var frameNumber int
		for {
			gtx := newContext(pixBuffer)
			publish := false
			switch e := w.NextEvent().(type) {
			case size.Event:
				sizeEvent = e
				// we need to create a new screen buffer, there's no way to resize the old one
				screenBuffer.Release()
				screenBuffer, err = s.NewBuffer(image.Point{e.WidthPx, e.HeightPx})
				if err != nil {
					panic(err)
				}
				pixBuffer = screenBuffer.RGBA()
			case mouse.Event:
				fmt.Println(e)
				uistate.mousex = float64(e.X)
				uistate.mousey = float64(e.Y)
			case key.Event:
				//words = words + string(e.Rune)
				//draw(pixBuffer, words)
				//w.Upload(image.Pt(0, 0), screenBuffer, sizeEvent.Bounds())
				//publish = true
			case paint.Event:
				//button(gtx, 10, 10, 100, 50)
				//button(gtx, 300, 100, 100, 50)
				//publish = true
			case lifecycle.Event:
				frameNumber++
				fmt.Printf("Event %d: From %s To %s\n", frameNumber, e.From, e.To)
				if e.To == lifecycle.StageDead {
					return
				}
			}
			button(gtx, 10, 10, 100, 50)
			button(gtx, 300, 100, 100, 50)
			publish = true
			if publish {
				w.Upload(image.Pt(0, 0), screenBuffer, sizeEvent.Bounds())
				w.Publish()
			}
		}
	})
	fmt.Println("end")
}

func regionhit(x, y, w, h float64) bool {
	if uistate.mousex < x || uistate.mousey < y || uistate.mousex >= x+w || uistate.mousey >= y+h {
		return false
	}
	return true
}

func button(gtx *gg.Context, x, y, w, h float64) {
	gtx.SetColor(color.White)
	if regionhit(x, y, w, h) {
		gtx.SetColor(color.RGBA{0, 255, 0, 255})
	}
	gtx.DrawRectangle(x, y, w, h)
	gtx.Fill()
}

func newContext(pixBuffer *image.RGBA) *gg.Context {
	dc := gg.NewContextForRGBA(pixBuffer)
	dc.Clear()
	return dc
}

func draw(pixBuffer *image.RGBA, words string) {
	dc := gg.NewContextForRGBA(pixBuffer)
	dc.Clear()
	f, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(f, &truetype.Options{Size: 40})
	dc.SetFontFace(face)

	start := 500.0
	color := 100.0
	for i := 0; i < 10; i++ {
		dc.SetColor(image.White)
		dc.DrawString(words, 10, 50)
		//dc.DrawStringWrapped(words, 10, 50, 100, 50, 50, 10.0, gg.AlignLeft)
		dc.DrawCircle(500, start, 400)
		dc.SetRGB(color/255.0, 100.0/255.0, 1)
		dc.Fill()

		start += 10
		color += 10
	}
}
