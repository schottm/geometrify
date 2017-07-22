package main

import "C"

import (
	"os"
	"fmt"
	"image/png"
	"geometrify/filter"
	"time"
	"flag"
	"runtime"
	"strings"
	"io/ioutil"
)

func main() {

	inputPtr := flag.String("i", "","the input file")
	outputPtr := flag.String("o", "", "the output file")

	iterationsPtr := flag.Int("n", 200, "the iterations")
	samplesPtr := flag.Int("s", 30, "the samples (should be a multiple of your thread count)")
	core := flag.Int("t", runtime.NumCPU(), "the thread count (default is cpu count)")

	flag.Parse()

	if *inputPtr == "" {
		panic("You have to select an input image!")
	}
	if *outputPtr == "" {
		panic( "You have to select an output image!")
	}

	if !strings.HasSuffix(*inputPtr, ".png") || !strings.HasSuffix(*outputPtr, ".png") {
		panic("You files have to be .png files!")
	}

	if !IsValid(*outputPtr) {
		panic(*outputPtr + " is not a valid path!")
	}

	var file, _ = os.Open(*inputPtr)
	defer file.Close()
	var img, err = png.Decode(file)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Trianglifining image: %s, ouput: %s.\n", file.Name(), *outputPtr)
	fmt.Printf("Parameters: %d iterations, %d samples, %d cores used.\n", *iterationsPtr, *samplesPtr, *core)

	var source = filter.DrawingFromImage(img)

	var start = time.Now().UnixNano()
	var drawing = filter.GenerateImage(source, filter.NewTriangleGenerator(1000000), *iterationsPtr, *core,  *samplesPtr)

	var diff = (time.Now().UnixNano() - start) / 1000000
	var min = diff / (60 * 1000)
	diff = diff % (60 * 1000)
	var sec = diff / 1000
	diff = diff % 1000

	fmt.Printf("Calculation time: %d min, %d sec, %d ms\n", min, sec, diff)

	var out *os.File = nil
	if _, err := os.Stat(*outputPtr); os.IsNotExist(err) {
		out, _ = os.Create(*outputPtr)
	} else {

		out, _ = os.OpenFile(*outputPtr, os.O_WRONLY, 0)
	}
	defer out.Close()

	err = png.Encode(out, filter.DrawingToImage(drawing))
	if err != nil {
		panic(err)
	}
}

//export Calculate
func Calculate(width, height int, data []int32, iterations, samples, threads int, seed int64) []int32 {

	source := filter.NewDrawing(width, height, false)
	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			pixel := data[y * width + x]
			alpha := uint16(pixel >> 24)
			red := uint16(pixel >> 16)
			green := uint16(pixel >> 8)
			blue := uint16(pixel)

			source.Set(x, y, red, green, blue, alpha)
		}
	}

	if threads == -1 {
		threads = runtime.NumCPU()
	}

	result := filter.GenerateImage(source, filter.NewTriangleGenerator(seed), iterations, threads, samples)

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			r, g, b, a := result.Get(x, y)
			pixel := int32(a) << 24 | int32(r) << 16 | int32(g) << 8 | int32(b)

			data[y * width + x] = pixel
		}
	}

	return data
}

func IsValid(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}