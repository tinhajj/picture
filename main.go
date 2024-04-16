package main

import (
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

var (
	winWidth, winHeight = 1000, 1000
	sizeEvent           size.Event
)

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
			case paint.Event:
				draw(pixBuffer)
				w.Upload(image.Pt(0, 0), screenBuffer, sizeEvent.Bounds())
			case lifecycle.Event:
				frameNumber++
				fmt.Printf("Event %d: From %s To %s\n", frameNumber, e.From, e.To)
				if e.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
	fmt.Println("end")
}

func draw(pixBuffer *image.RGBA) {
	dc := gg.NewContextForRGBA(pixBuffer)

	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(100.0/255.0, 200.0/255.0, 1)
	dc.Fill()

	dc.DrawCircle(500, 900, 400)
	dc.SetRGB(200.0/255.0, 200.0/255.0, 0.5)
	dc.Fill()
}
