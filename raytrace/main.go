/*
The raytracer application takes a scene as a json file and produces a rendered version of that scene as
a png or jpeg image. See lib.Options for details on all the options.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/Bredgren/gotracer/lib"
	"github.com/Bredgren/gotracer/trace"
)

const (
	usageStr = `Usage: raytracer [OPTIONS]
Reads a json-formatted scene description from stdin and writes the image to stdout.
`
)

var (
	imgFormats = [...]string{"jpg", "jpeg", "png"}
)

// cmdline options
var (
	inFile      = ""
	outFile     = ""
	gridSize    = 50
	format      = "jpg"
	jpegQuality = 95
	// verbose     = false
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usageStr)
		flag.PrintDefaults()
	}
	flag.StringVar(&inFile, "in", inFile, "Read scene from this json file instead of stdin")
	flag.StringVar(&outFile, "out", outFile, "Save image here instead of sending to stdout. If the format option is unspecified then it will attempt to determine the format from the file name")
	flag.IntVar(&gridSize, "gridSize", gridSize, "The image is divided into a grid an each section rendered in parallel. This is the size in pixels of each square")
	flag.StringVar(&format, "format", format, fmt.Sprintf("Image format, one of %s", imgFormats))
	flag.IntVar(&jpegQuality, "jpegQuality", jpegQuality, "JPEG quality")
	// flag.BoolVar(&verbose, "v", verbose, "Print progress reports")
}

// type grid struct {
// 	scene                  *raytracer.Scene
// 	xMin, yMin, xMax, yMax int // Min inclusive, Max exclusive
// 	img                    *image.NRGBA
// 	ch                     chan float64 // signal done
// }

// func traceGrid(g grid) {
// 	for y := g.yMin; y < g.yMax; y++ {
// 		for x := g.xMin; x < g.xMax; x++ {
// 			g.img.SetNRGBA(x, y, g.scene.TracePixel(x, y))
// 		}
// 	}
// 	g.ch <- 1
// }

func main() {
	flag.Parse()

	// Verify requested format
	imgFormat := ""
	for _, f := range imgFormats {
		if format == f {
			imgFormat = format
			break
		}
	}
	if imgFormat == "" {
		fmt.Fprintf(os.Stderr, "Invalid image format '%s'\n", format)
		flag.Usage()
		return
	}

	// Choose stdin or user-specified file
	var inReader io.Reader = os.Stdin
	if inFile != "" {
		file, e := os.Open(inFile)
		if e != nil {
			log.Fatalf("opening file to read: %s\n", e)
		}
		defer file.Close()
		inReader = io.Reader(file)
	}

	// Read and decode options
	var options lib.Options
	jsonOpts, e := ioutil.ReadAll(inReader)
	if e != nil {
		log.Fatalf("reading file: %s\n", e)
	}
	e = json.Unmarshal(jsonOpts, &options)
	if e != nil {
		log.Fatalf("decoding json: %s\n", e)
	}

	img := trace.Trace(options, gridSize)

	// Choose stdout or user-specified file
	var outWriter io.Writer = os.Stdout
	if outFile != "" {
		var e error
		outWriter, e = os.Create(outFile)
		if e != nil {
			log.Fatalf("creating file to write: %s\n", e)
		}
	}

	// Send/save image
	switch format {
	case "jpg", "jpeg":
		jpeg.Encode(outWriter, img, &jpeg.Options{Quality: jpegQuality})
	case "png":
		png.Encode(outWriter, img)
	}

	// 	sceneFilePath := sceneDir + "/" + sceneFile
	// 	scene := raytracer.Parse(sceneFilePath)

	// 	imgW = scene.Camera.ImageWidth
	// 	imgH = scene.Camera.ImageHeight
	// 	// rayCounts = make([][]map[int]int, imgW)
	// 	// for x := 0; x < imgW; x++ {
	// 	// 	rayCounts[x] = make([]map[int]int, imgH)
	// 	// 	for y := 0; y < imgH; y++ {
	// 	// 		rayCounts[x][y] = make(map[int]int)
	// 	// 	}
	// 	// }
	// 	bounds := image.Rect(0, 0, imgW, imgH)
	// 	img := image.NewNRGBA(bounds)

	// 	gridChan := make(chan float64, 100)
	// 	gridCount := 0.0

	// 	log.Println("Begin tracing")
	// 	begin := time.Now()
	// 	for y := bounds.Min.Y; y < bounds.Max.Y; y += gridSize {
	// 		for x := bounds.Min.X; x < bounds.Max.X; x += gridSize {
	// 			xMax := x + gridSize
	// 			if xMax > imgW {
	// 				xMax = imgW
	// 			}
	// 			yMax := y + gridSize
	// 			if yMax > imgH {
	// 				yMax = imgH
	// 			}
	// 			go traceGrid(grid{scene, x, y, xMax, yMax, img, gridChan})
	// 			gridCount++
	// 		}
	// 	}
	// 	for g := 0.0; g < gridCount; g += <-gridChan {
	// 		log.Printf("%.2f%%", g/gridCount*100)
	// 	}
	// 	end := time.Now()
	// 	log.Printf("Done tracing, took %v", end.Sub(begin))

	// 	if noImg {
	// 		return
	// 	}

	// 	files, err := ioutil.ReadDir(renderDir)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	count := 0
	// 	for _, file := range files {
	// 		name := file.Name()
	// 		if strings.HasPrefix(name, "render") && len(name) > 10 &&
	// 			!strings.Contains(name, "Rays") && !strings.Contains(name, "scene") {
	// 			number, err := strconv.Atoi(name[6 : len(name)-4])
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			if number > count {
	// 				count = number
	// 			}
	// 		}
	// 	}

	// 	baseFileName := fmt.Sprintf("%s/render%d", renderDir, count+1)

	// 	outFile := fmt.Sprintf("%s.%s", baseFileName, format)

	// 	file, err := os.Create(outFile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer file.Close()

	// 	if format == "jpg" {
	// 		jpeg.Encode(file, img, &jpeg.Options{jpegQuality})
	// 	} else if format == "png" {
	// 		png.Encode(file, img)
	// 	}
	// 	log.Printf("Saved %s", outFile)

	// 	copyFileContents(sceneFilePath, baseFileName+"scene.json")

	// 	// saveRayCounts(baseFileName)
	// }

	// func copyFileContents(src, dst string) (err error) {
	// 	in, err := os.Open(src)
	// 	if err != nil {
	// 		return
	// 	}
	// 	defer in.Close()
	// 	out, err := os.Create(dst)
	// 	if err != nil {
	// 		return
	// 	}
	// 	defer func() {
	// 		cerr := out.Close()
	// 		if err == nil {
	// 			err = cerr
	// 		}
	// 	}()
	// 	if _, err = io.Copy(out, in); err != nil {
	// 		return
	// 	}
	// 	err = out.Sync()
	// 	log.Printf("Saved %s", dst)
	// 	return
}
