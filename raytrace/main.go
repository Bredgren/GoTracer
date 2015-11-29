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
	var options trace.Options
	jsonOpts, e := ioutil.ReadAll(inReader)
	if e != nil {
		log.Fatalf("reading file: %s\n", e)
	}
	e = json.Unmarshal(jsonOpts, &options)
	if e != nil {
		log.Fatalf("decoding json: %s\n", e)
	}

	img := trace.Trace(&options, gridSize)

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
}
