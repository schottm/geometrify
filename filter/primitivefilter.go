package filter

import (
	"fmt"
	"math"
)

type result struct {

	primitive Primitive
	impact int
}

func GenerateImage(source *Drawing, generator PrimitiveGenerator, iterations, parallel, samples int) *Drawing {

	parallel = min(samples, parallel)
	//calculate needed the samples to calculate per thread
	var target = NewDrawing(source.Width, source.Height, source.Opaque())
	var shifted = NewDrawing(source.Width, source.Height, false)

	fmt.Printf("Simplifing image(%dx%d): %d iterations, %d samples, %d threads.\n", source.Width, source.Height, iterations, samples, parallel)

	if source.Opaque() {

		for y := 0; y < shifted.Height; y++ {

			for x := 0; x < shifted.Width; x++ {

				shifted.Set(x, y, 0, 0, 0, (math.MaxUint8 >> 1) + 1)
			}
		}
	}

	//create channel for thread (equivalent to wait in java, but 1000-times better xD)
	var channel chan *result = make(chan *result, parallel)
	//var generators = make([]PrimitiveGenerator, parallel)

	//the random func in go is not concurrent -> we have to create a random generator for each thread
	/*
	for i := range generators {

		generators[i] = create(time.Now().UnixNano() * int64(i + 1))
	}
	*/

	for i := 0; i < iterations; i++ {

		var process chan Primitive = make(chan Primitive, min(samples, parallel * 4))

		if i % 10 == 0 {
			fmt.Println("Iteration: ", i)
		}

		//calculate the samples in different threads

		for j := 0; j < parallel; j++ {

			go func() {
				r := findPrimitive(source, target, shifted, process)

				channel <- r
			}()
		}

		for j := 0; j < samples; j++ {
			process <- generator.generate(source.Width, source.Height)
		}

		close(process)

		//find the primitive with least weight
		var overallRes *result = nil
		for j := 0; j < parallel; j++ {

			primitive := <-channel

			if overallRes == nil {

				if primitive.impact > 0 {

					overallRes = primitive
				}

			} else if overallRes.impact < primitive.impact {

				overallRes = primitive
			}
		}

		if overallRes != nil {

			addToImage(target, shifted, overallRes.primitive)
		}
	}

	close(channel)

	return target
}

func findPrimitive(source, target, shifted *Drawing, queue <-chan Primitive) *result {

	var impact int = 0
	var match Primitive = nil

	for {

		primitive, ok := <-queue
		if !ok {
			return &result{match, impact}
		}

		if primitive != nil {
			primitive.SetColor(generateColor(source, primitive))

			var currentImpact = generateImpact(source, target, shifted, primitive)

			if currentImpact >= impact {

				impact = currentImpact
				match = primitive
			}
		}
	}
}

func generateColor(image *Drawing, primitive Primitive) Color {

	var red, green, blue, alpha int
	var count int = 0

	for y := primitive.Bounds().Min.Y; y <= primitive.Bounds().Max.Y; y++ {

		for x := primitive.Bounds().Min.X; x <= primitive.Bounds().Max.X; x++ {

			if primitive.IsInside(x, y) {

				var i = image.indexOf(x, y)
				red += int(image.data[i])
				green += int(image.data[i + 1])
				blue += int(image.data[i + 2])
				alpha += int(image.data[i + 3])

				count++
			}
		}
	}

	if count == 0 {

		return Color{0, 0, 0, 0}
	}

	return Color{R: uint8(red / count), G: uint8(green / count), B: uint8(blue / count), A: uint8(alpha / count)}
}

func generateImpact(source, target, shifted *Drawing, primitive Primitive) int {

	var impact int
	var r = primitive.GetColor().R >> 1
	var g = primitive.GetColor().G >> 1
	var b = primitive.GetColor().B >> 1
	var a = primitive.GetColor().A >> 1

	//var sr, sg, sb, sa uint16

	for y := primitive.Bounds().Min.Y; y <= primitive.Bounds().Max.Y; y++ {

		for x := primitive.Bounds().Min.X; x <= primitive.Bounds().Max.X; x++ {

			if primitive.IsInside(x, y) {

				var i = source.indexOf(x, y)

				impact += int(diff(source.data[i], target.data[i]))
				impact += int(diff(source.data[i + 1], target.data[i + 1]))
				impact += int(diff(source.data[i + 2], target.data[i + 2]))
				impact += int(diff(source.data[i + 3], target.data[i + 3]))

				impact -= int(diff(source.data[i], r + shifted.data[i]))
				impact -= int(diff(source.data[i + 1], g + shifted.data[i + 1]))
				impact -= int(diff(source.data[i + 2], b + shifted.data[i + 2]))
				impact -= int(diff(source.data[i + 3], a + shifted.data[i + 3]))

				/*
				sr = source.data[i]
				sg = source.data[i + 1]
				sb = source.data[i + 2]
				sa = source.data[i + 3]
				impact += int(diff(sr, target.data[i]))
				impact += int(diff(sg, target.data[i + 1]))
				impact += int(diff(sb, target.data[i + 2]))
				impact += int(diff(sa, target.data[i + 3]))
				impact -= int(diff(sr, r + shifted.data[i]))
				impact -= int(diff(sg, g + shifted.data[i + 1]))
				impact -= int(diff(sb, b + shifted.data[i + 2]))
				impact -= int(diff(sa, a + shifted.data[i + 3]))
				*/
			}
		}
	}


	return impact
}

func addToImage(target, shifted *Drawing, primitive Primitive) {

	var r = primitive.GetColor().R >> 1
	var g = primitive.GetColor().G >> 1
	var b = primitive.GetColor().B >> 1
	var a = primitive.GetColor().A >> 1

	for y := primitive.Bounds().Min.Y; y <= primitive.Bounds().Max.Y; y++ {

		for x := primitive.Bounds().Min.X; x <= primitive.Bounds().Max.X; x++ {

			if primitive.IsInside(x, y) {

				var i = target.indexOf(x, y)

				target.data[i] = r + shifted.data[i]
				target.data[i + 1] = g + shifted.data[i + 1]
				target.data[i + 2] = b + shifted.data[i + 2]
				target.data[i + 3] = a + shifted.data[i + 3]

				shifted.data[i] = (target.data[i] >> 1) + 1
				shifted.data[i + 1] = (target.data[i + 1] >> 1) + 1
				shifted.data[i + 2] = (target.data[i + 2] >> 1) + 1
				shifted.data[i + 3] = (target.data[i + 3] >> 1) + 1
			}
		}
	}
}

func diff(a, b uint8) uint8 {

	if a > b {
		return a - b
	}
	return b - a
}



