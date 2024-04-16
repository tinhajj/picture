package main

import (
	"fmt"
	"image"

	"github.com/fogleman/gg"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
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
			case key.Event:
				draw(pixBuffer, e.String())
				w.Upload(image.Pt(0, 0), screenBuffer, sizeEvent.Bounds())
			case paint.Event:
				draw(pixBuffer, "nothing")
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

func draw(pixBuffer *image.RGBA, words string) {
	dc := gg.NewContextForRGBA(pixBuffer)
	dc.Clear()

	start := 500.0
	color := 100.0
	for i := 0; i < 10; i++ {
		dc.SetColor(image.White)
		dc.DrawString(words, 10, 10)
		dc.DrawCircle(500, start, 400)
		dc.SetRGB(color/255.0, 100.0/255.0, 1)
		dc.Fill()

		start += 10
		color += 10
	}
}
