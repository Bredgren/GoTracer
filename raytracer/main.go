/*
TODO:
 - Reflection
 - Refraction
 - Texture mapping
 - Bump mapping
 - Antialiasing
 - Skybox
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
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usageStr)
		flag.PrintDefaults()
	}
	flag.BoolVar(&noImg, "NoImg", false, "Don't create an image if present.")
	flag.Parse()

	if flag.NArg() < 1 {
		panic(usageStr)
	}

	sceneFile = flag.Arg(0)
}

func main() {
	scene := raytracer.Parse(sceneDir + "/" + sceneFile)

	bounds := image.Rect(0, 0, scene.Camera.ImageWidth, scene.Camera.ImageHeight)
	img := image.NewNRGBA(bounds)

	log.Println("Begin tracing")
	begin := time.Now()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.SetNRGBA(x, y, scene.TracePixel(x, y))
		}
	}
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
