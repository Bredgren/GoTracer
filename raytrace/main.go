/*
The raytracer application takes a scene as a json file and produces a rendered version of that scene as
a png or jpeg image. See trace.Options for details on all the options.
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
	"path/filepath"

	"github.com/Bredgren/gotracer/trace"
	"github.com/Bredgren/gotracer/trace/options"
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
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usageStr)
		flag.PrintDefaults()
	}
	flag.StringVar(&inFile, "in", inFile, "Read scene from this json file instead of stdin")
	flag.StringVar(&outFile, "out", outFile, "Save image here instead of sending to stdout. The file extension will override the format option")
	flag.IntVar(&gridSize, "gridSize", gridSize, "The image is divided into a grid an each section rendered in parallel. This is the size in pixels of each square")
	flag.StringVar(&format, "format", format, fmt.Sprintf("Image format, one of %s", imgFormats))
	flag.IntVar(&jpegQuality, "jpegQuality", jpegQuality, "JPEG quality")
}

func main() {
	flag.Parse()
	log.SetPrefix("raytrace: ")
	log.SetFlags(0)

	// Use file extension
	if outFile != "" {
		format = filepath.Ext(outFile)[1:]
	}

	// Verify requested format
	imgFormat := ""
	for _, f := range imgFormats {
		if format == f {
			imgFormat = format
			break
		}
	}
	if imgFormat == "" {
		fmt.Fprintf(os.Stderr, "Invalid image format: %s\n", format)
		flag.Usage()
		return
	}

	// Choose stdin or user-specified file
	var inReader io.Reader = os.Stdin
	if inFile != "" {
		f, e := os.Open(inFile)
		if e != nil {
			log.Fatalf("opening file to read: %s\n", e)
		}
		inReader = f
		defer func() {
			if e := f.Close(); e != nil {
				log.Fatalf("closing input file: %v\n", e)
			}
		}()
	}

	// Read and decode options
	var options options.Options
	jsonOpts, e := ioutil.ReadAll(inReader)
	if e != nil {
		log.Fatalf("reading options: %v\n", e)
	}
	if e := json.Unmarshal(jsonOpts, &options); e != nil {
		log.Fatalf("decoding options: %v\n", e)
	}

	img := trace.Trace(&options, gridSize)

	// Choose stdout/user-specified file
	var outWriter io.Writer = os.Stdout
	if outFile != "" {
		f, e := os.Create(outFile)
		if e != nil {
			log.Fatalf("creating output file: %v\n", e)
		}
		outWriter = f
		defer func() {
			if e := f.Close(); e != nil {
				log.Fatalf("closing output file: %v\n", e)
			}
		}()
	}

	// Send/save image
	switch format {
	case "jpg", "jpeg":
		jpeg.Encode(outWriter, img, &jpeg.Options{Quality: jpegQuality})
	case "png":
		png.Encode(outWriter, img)
	default:
		panic(fmt.Sprintf("Saving image as format %s is not implemented", format))
	}
}
