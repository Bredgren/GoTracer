/*
TODO:
 - Torus
 - Disc
 - Texture offset
 - Beers law (again)
 - Skybox
 - BSPTree
 - Trimesh
 - CSG
 - Caustics
 - Distance estimators
 - Fractals
*/
package main

import (
	"flag"
	"fmt"
	"image"
	// "image/color"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Bredgren/gotracer"
)

const (
	sceneDir = "scene"
	renderDir = "renderedscene"
	usageStr = "Usage: raytracer scenefile"
)

var (
	// rayCounts [][]map[int]int
	imgW int
	imgH int
)

var (
	sceneFile = ""
	noImg = false
	gridSize = 50
	format = "jpg"
	jpegQuality = 95
	countRays = false
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usageStr)
		flag.PrintDefaults()
	}
	flag.BoolVar(&noImg, "NoImg", noImg, "Don't create an image if present.")
	flag.IntVar(&gridSize, "gridSize", gridSize, "Size of simultaneous trace grids.")
	flag.StringVar(&format, "format", format, "Image format [jpg | png].")
	flag.IntVar(&jpegQuality, "jpegQuality", jpegQuality, "JPEG quality.")
	flag.BoolVar(&countRays, "CountRays", countRays, "Make images showing ray counts.")
	flag.Parse()

	if flag.NArg() < 1 {
		panic(usageStr)
	}

	if format != "jpg" && format != "png" {
		panic("Unknown format " + format)
	}

	sceneFile = flag.Arg(0)

	runtime.GOMAXPROCS(runtime.NumCPU())
}

type grid struct {
	scene *raytracer.Scene
	xMin, yMin, xMax, yMax int // Min inclusive, Max exclusive
	img *image.NRGBA
	ch chan float64 // signal done
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
	sceneFilePath := sceneDir + "/" + sceneFile
	scene := raytracer.Parse(sceneFilePath)

	imgW = scene.Camera.ImageWidth
	imgH = scene.Camera.ImageHeight
	// rayCounts = make([][]map[int]int, imgW)
	// for x := 0; x < imgW; x++ {
	// 	rayCounts[x] = make([]map[int]int, imgH)
	// 	for y := 0; y < imgH; y++ {
	// 		rayCounts[x][y] = make(map[int]int)
	// 	}
	// }
	bounds := image.Rect(0, 0, imgW, imgH)
	img := image.NewNRGBA(bounds)

	gridChan := make(chan float64, 100)
	gridCount := 0.0

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
	for g := 0.0; g < gridCount; g += <-gridChan {
		log.Printf("%.2f%%", g / gridCount * 100)
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
		if strings.HasPrefix(name, "render") && len(name) > 10 &&
			!strings.Contains(name, "Rays") && !strings.Contains(name, "scene") {
			number, err := strconv.Atoi(name[6:len(name) - 4])
			if err != nil {
				log.Fatal(err)
			}
			if number > count {
				count = number
			}
		}
	}

	baseFileName := fmt.Sprintf("%s/render%d", renderDir, count + 1)

	outFile := fmt.Sprintf("%s.%s", baseFileName, format)

 	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if format == "jpg" {
		jpeg.Encode(file, img, &jpeg.Options{jpegQuality})
	} else if format == "png" {
		png.Encode(file, img)
	}
	log.Printf("Saved %s", outFile)

	copyFileContents(sceneFilePath, baseFileName + "scene.json")

	// saveRayCounts(baseFileName)
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	log.Printf("Saved %s", dst)
	return
}

// func saveRayCounts(baseFileName string) {
// 	imgs := new([raytracer.NumRayTypes]*image.NRGBA)
// 	bounds := image.Rect(0, 0, imgW, imgH)
// 	for i := 0; i < len(imgs); i++ {
// 		imgs[i] = image.NewNRGBA(bounds)
// 	}

// 	maxCount := new([raytracer.NumRayTypes]int)
// 	for x := 0; x < len(rayCounts); x++ {
// 		for y := 0; y < len(rayCounts[x]); y++ {
// 			for rayType := 0; rayType < raytracer.NumRayTypes; rayType++ {
// 				if maxCount[rayType] < rayCounts[x][y][rayType] {
// 					maxCount[rayType] = rayCounts[x][y][rayType]
// 				}
// 			}
// 		}
// 	}

// 	for x := 0; x < len(rayCounts); x++ {
// 		for y := 0; y < len(rayCounts[x]); y++ {
// 			for rayType := 0; rayType < raytracer.NumRayTypes; rayType++ {
// 				ratio := float64(rayCounts[x][y][rayType]) / float64(maxCount[rayType])
// 				shade := uint8(ratio * 255)
// 				color := color.NRGBA{shade, shade, shade, 255}
// 				imgs[rayType].SetNRGBA(x, y, color)
// 			}
// 		}
// 	}

// 	fileName := fmt.Sprintf("%sPrimaryRays.png", baseFileName)
// 	saveRayFile(fileName, imgs[raytracer.PrimaryRay])
// 	fileName = fmt.Sprintf("%sShadowRays.png", baseFileName)
// 	saveRayFile(fileName, imgs[raytracer.ShadowRay])
// }

// func saveRayFile(fileName string, img *image.NRGBA) {
//  	file, err := os.Create(fileName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	png.Encode(file, img)
// 	log.Printf("Saved %s", fileName)
// }
