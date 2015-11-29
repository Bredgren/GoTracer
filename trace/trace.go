package trace

import (
	"image"
	"image/color"
	"math/rand"
	"sync"
)

// Trace renders an image according to the given options. Renders the image in chuncks of the given
// size in parallel.
func Trace(options *Options, gridSize int) *image.NRGBA {
	// Create scene
	// ...

	imgW := options.Resolution.W
	imgH := options.Resolution.H
	bounds := image.Rect(0, 0, imgW, imgH)
	img := image.NewNRGBA(bounds)

	gridGroup := sync.WaitGroup{}

	for y := bounds.Min.Y; y < bounds.Max.Y; y += gridSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += gridSize {
			xMax := x + gridSize
			if xMax > imgW {
				xMax = imgW
			}
			yMax := y + gridSize
			if yMax > imgH {
				yMax = imgH
			}
			x, y := x, y
			gridGroup.Add(1)
			go func() {
				for y := y; y < yMax; y++ {
					for x := x; x < xMax; x++ {
						// color := scene.TracePixel(x, y)
						color := color.NRGBA{
							R: uint8(rand.Intn(255)),
							G: uint8(rand.Intn(255)),
							B: uint8(rand.Intn(255)),
							A: 255,
						}
						img.SetNRGBA(x, y, color)
					}
				}
				gridGroup.Done()
			}()
		}
	}
	gridGroup.Wait()

	return img
}
