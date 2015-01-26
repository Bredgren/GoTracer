/*
TODO:
 - Soft shadows
 - Beer's law
 - Fresnel term
 - Texture mapping
 - Bump mapping
 - Antialiasing
 - BSPTree
 - Cylinder
 - Cone
 - Torus
 - Disc
 - Skybox
 - Trimesh
 - CSG
 - Caustics
*/
package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Bredgren/raytracer"
)

const (
	sceneDir = "scene"
	renderDir = "renderedscene"
	usageStr = "Usage: raytracer scenefile"
)

var (
	sceneFile = ""
	noImg = false
	gridSize = 50
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usageStr)
		flag.PrintDefaults()
	}
	flag.BoolVar(&noImg, "NoImg", false, "Don't create an image if present.")
	flag.IntVar(&gridSize, "gridSize", gridSize, "Size of simultaneous trace grids.")
	flag.Parse()

	if flag.NArg() < 1 {
		panic(usageStr)
	}

	sceneFile = flag.Arg(0)

	runtime.GOMAXPROCS(runtime.NumCPU())
}

type grid struct {
	scene *raytracer.Scene
	xMin, yMin, xMax, yMax int // Min inclusive, Max exclusive
	img *image.NRGBA
	ch chan int // signal done
}

func traceGrid(g grid) {
	for y := g.yMin; y < g.yMax; y++ {
		for x := g.xMin; x < g.xMax; x++ {
			g.img.SetNRGBA(x, y, g.scene.TracePixel(x, y))
		}
	}
	g.ch <- 1
}

func main() {
	scene := raytracer.Parse(sceneDir + "/" + sceneFile)

	imgW := scene.Camera.ImageWidth
	imgH := scene.Camera.ImageHeight
	bounds := image.Rect(0, 0, imgW, imgH)
	img := image.NewNRGBA(bounds)

	gridChan := make(chan int, 100)
	gridCount := 0

	log.Println("Begin tracing")
	begin := time.Now()
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
			go traceGrid(grid{scene, x, y, xMax, yMax, img, gridChan})
			gridCount++
		}
	}
	for g := 0; g < gridCount; g += <-gridChan {}
	end := time.Now()
	log.Printf("Done tracing, took %v", end.Sub(begin))

	if noImg {
		return
	}

	files, err := ioutil.ReadDir(renderDir)
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for _, file := range files {
		name :=file.Name()
		if strings.HasPrefix(name, "render") && len(name) > 10 {
			number, err := strconv.Atoi(name[6:len(name) - 4])
			if err != nil {
				log.Fatal(err)
			}
			if number > count {
				count = number
			}
		}
	}

	outFile := fmt.Sprintf("%s/render%d.png", renderDir, count + 1)

	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	png.Encode(file, img)
}
