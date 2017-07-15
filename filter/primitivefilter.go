package filter

import (
	"fmt"
	"sync"
)

type result struct {

	primitive Primitive
	impact int
}

func GenerateImage(source *Drawing, generator PrimitiveGenerator, iterations, parallel, samples int) *Drawing {

	parallel = min(samples, parallel)
	//calculate needed the samples to calculate per thread
	var target = NewDrawing(source.Width, source.Height, source.Opaque())
	fmt.Printf("Simplifing image(%dx%d): %d iterations, %d samples, %d threads.\n", source.Width, source.Height, iterations, samples, parallel)

	//create channel for thread (equivalent to wait in java, but 1000-times better xD)
	//var channel chan *result = make(chan *result, parallel)
	//var generators = make([]PrimitiveGenerator, parallel)

	//the random func in go is not concurrent -> we have to create a random generator for each thread
	/*
	for i := range generators {

		generators[i] = create(time.Now().UnixNano() * int64(i + 1))
	}
	*/
	var wg = &sync.WaitGroup{}
	var mutex = &sync.Mutex{}

	for i := 0; i < iterations; i++ {

		var process chan Primitive = make(chan Primitive, samples)

		if i % 10 == 0 {
			fmt.Println("Iteration: ", i)
		}

		//calculate the samples in different threads
		var impact = 0
		var result Primitive = nil
		//var overallRes *result = nil

		for j := 0; j < parallel; j++ {

			go func() {

				/*
				for primitive := range process {

					current := findImpact(source, target, primitive)

					mutex.Lock()
					if current > impact {
						impact = current
						result = primitive
					}
					mutex.Unlock()
				}
				*/

				primitive, current := findPrimitive(source, target, process)

				//replace current primitive if this one is better
				mutex.Lock()
				if current > impact {
					impact = current
					result = primitive
				}
				/*
				if overallRes == nil {

					if primitive.impact > 0 {

						overallRes = primitive
					}

				} else if overallRes.impact < primitive.impact {

					overallRes = primitive
				}
				*/
				mutex.Unlock()

				wg.Done()
			}()
		}

		wg.Add(parallel)
		for j := 0; j < samples; j++ {
			process <- generator.generate(source.Width, source.Height)
		}
		close(process)
		//wait for calculation to finish
		wg.Wait()

		if result != nil {

			addToImage(target, result)
		}
	}

	//close(channel)

	return target
}

func findImpact(source, target *Drawing, primitive Primitive) int {

	primitive.SetColor(generateColor(source, primitive))

	return generateImpact(source, target, primitive)
}

func findPrimitive(source, target *Drawing, queue <-chan Primitive) (Primitive, int) {

	var impact int = 0
	var match Primitive = nil

	for primitive := range queue {

		if primitive != nil {
			primitive.SetColor(generateColor(source, primitive))

			var currentImpact = generateImpact(source, target, primitive)

			if currentImpact >= impact {

				impact = currentImpact
				match = primitive
			}
		}
	}
	return match, impact
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

	return Color{R: uint16(red / count), G: uint16(green / count), B: uint16(blue / count), A: uint16(alpha / count)}
}

func generateImpact(source, target *Drawing, primitive Primitive) int {

	var impact int
	var r = primitive.GetColor().R
	var g = primitive.GetColor().G
	var b = primitive.GetColor().B
	var a = primitive.GetColor().A

	//var sr, sg, sb, sa uint16

	for y := primitive.Bounds().Min.Y; y <= primitive.Bounds().Max.Y; y++ {

		for x := primitive.Bounds().Min.X; x <= primitive.Bounds().Max.X; x++ {

			if primitive.IsInside(x, y) {

				var i = source.indexOf(x, y)

				impact += avgdiff(source.data[i], target.data[i], r)
				impact += avgdiff(source.data[i + 1], target.data[i + 1], g)
				impact += avgdiff(source.data[i + 2], target.data[i + 2], b)
				impact += avgdiff(source.data[i + 3], target.data[i + 3], a)

				/*
				impact += int(diff(source.data[i], target.data[i]))
				impact += int(diff(source.data[i + 1], target.data[i + 1]))
				impact += int(diff(source.data[i + 2], target.data[i + 2]))
				impact += int(diff(source.data[i + 3], target.data[i + 3]))

				impact -= int(diff(source.data[i], (r + target.data[i]) >> 1))
				impact -= int(diff(source.data[i + 1], (g + target.data[i + 1]) >> 1))
				impact -= int(diff(source.data[i + 2], (b + target.data[i + 2]) >> 1))
				impact -= int(diff(source.data[i + 3], (a + target.data[i + 3]) >> 1))
				*/

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

func addToImage(target *Drawing, primitive Primitive) {

	var r = primitive.GetColor().R
	var g = primitive.GetColor().G
	var b = primitive.GetColor().B
	var a = primitive.GetColor().A

	for y := primitive.Bounds().Min.Y; y <= primitive.Bounds().Max.Y; y++ {

		for x := primitive.Bounds().Min.X; x <= primitive.Bounds().Max.X; x++ {

			if primitive.IsInside(x, y) {

				var i = target.indexOf(x, y)

				target.data[i] = (r + target.data[i]) >> 1
				target.data[i + 1] = (g + target.data[i + 1]) >> 1
				target.data[i + 2] = (b + target.data[i + 2]) >> 1
				target.data[i + 3] = (a + target.data[i + 3]) >> 1
			}
		}
	}
}

func avgdiff(source, target, primitive uint16) int {

	tc := int(target + primitive) >> 1
	s := int(source)
	t := int(target)
	if s > t {

		if s > tc {

			return tc - t
		} else {

			return (s << 1) - t - tc
		}
	} else {

		if s > tc {

			return t + tc - (s << 1)
		} else {

			return t - tc
		}
	}
}



