package filter

type Triangle struct {

	A Point
	B Point
	C Point

	/*
	y23 int
	x32 int
	y31 int
	x13 int
	det int
	minD int
	maxD int
	*/

	//Bound BoundingBox

	Color
}

/*
func (triangle *Triangle) IsInside(x, y int) bool {


	dx := x - triangle.C.X
	dy := y - triangle.C.Y
	a := triangle.y23 * dx + triangle.x32 * dy

	if a < triangle.minD || a > triangle.maxD {
		return false
	}
	b := triangle.y31 * dx + triangle.x13 * dy
	if b < triangle.minD || b > triangle.maxD {
		return false
	}
	c := triangle.det - a - b
	if c < triangle.minD || c > triangle.maxD {
		return false
	}
	return true

	/*
	var asX int = x - triangle.A.X
	var asY int = y - triangle.A.Y

	var ab bool = (triangle.B.X - triangle.A.X) * asY - (triangle.B.Y - triangle.A.Y) * asX > 0

	if (triangle.C.X - triangle.A.X) * asY - (triangle.C.Y - triangle.A.Y) * asX > 0 == ab {
		return false
	}

	return (triangle.C.X - triangle.B.X) * (y - triangle.B.Y) - (triangle.C.Y - triangle.B.Y) * (x - triangle.B.X) > 0 == ab

}
*/

func (triangle *Triangle) ForEach(do func(int, int)) {

	if triangle.B.Y == triangle.C.Y {
		fillBottomFlatTriangle(triangle.A, triangle.B, triangle.C, do)
	} else if triangle.A.Y == triangle.B.Y {
		fillTopFlatTriangle(triangle.A, triangle.B, triangle.C, do)
	} else {
		x := triangle.A.X + int(float32(triangle.B.Y - triangle.A.Y) / float32(triangle.C.Y - triangle.A.Y) * float32(triangle.C.X - triangle.A.X))
		d := Point{x, triangle.B.Y}
		fillBottomFlatTriangle(triangle.A, triangle.B, d, do)
		fillTopFlatTriangle(triangle.B, d, triangle.C, do)
	}
}

func fillBottomFlatTriangle(p1, p2, p3 Point, do func(int, int)) {
	invslope1 := float32(p2.X - p1.X) / float32(p2.Y - p1.Y)
	invslope2 := float32(p3.X - p1.X) / float32(p3.Y - p1.Y)

	curx1 := float32(p1.X)
	curx2 := float32(p1.X)

	for scanlineY := p1.Y; scanlineY <= p2.Y; scanlineY++ {
		drawLine(int(curx1), int(curx2), scanlineY, do)
		curx1 += invslope1
		curx2 += invslope2
	}
}

func fillTopFlatTriangle(p1, p2, p3 Point, do func(int, int)) {
	invslope1 := float32(p3.X - p1.X) / float32(p3.Y - p1.Y)
	invslope2 := float32(p3.X - p2.X) / float32(p3.Y - p2.Y)

	curx1 := float32(p3.X)
	curx2 := float32(p3.X)

	for scanlineY := p3.Y; scanlineY > p1.Y; scanlineY-- {
		drawLine(int(curx1), int(curx2), scanlineY, do)
		curx1 -= invslope1
		curx2 -= invslope2
	}
}

func drawLine(x1, x2, y int, do func(int, int)) {

	if x1 > x2 {
		temp := x1
		x1 = x2
		x2 = temp
	}

	for i := x1; i <= x2; i++ {
		do(i, y)
	}
}

/*
func (triangle *Triangle) Bounds() *BoundingBox {

	return &triangle.Bound
}
*/

func (triangle *Triangle) GetColor() *Color {

	return &triangle.Color
}

func (triangle *Triangle) SetColor(color Color) {

	triangle.Color = color
}

func NewTriangle(a Point, b Point, c Point) *Triangle {

	p1 := a
	p2 := b
	p3 := c
	if b.Y < p1.Y {
		p2 = p1
		p1 = b
	}
	if c.Y < p1.Y {
		p3 = p1
		p1 = c
	}
	if p3.Y < p2.Y {
		temp := p2
		p2 = p3
		p3 = temp
	}
	/*
	y23 := b.Y - c.Y
	x32 := c.X - b.X
	y31 := c.Y - a.Y
	x13 := a.X - c.X
	det := y23 * x13 - x32 * y31
	minD := min(det, 0)
	maxD := max(det, 0)

	var triangle = Triangle{a, b, c, y23, x32, y31, x13, det, minD, maxD, BoundingBox{}, Color{}}

	*/
	var triangle = Triangle{A: p1, B: p2, C: p3, Color: Color{}}

	/*
	triangle.Bound.Min.X = min(a.X, min(b.X, c.X))
	triangle.Bound.Min.Y = min(a.Y, min(b.Y, c.Y))
	triangle.Bound.Max.X = max(a.X, max(b.X, c.X))
	triangle.Bound.Max.Y = max(a.Y, max(b.Y, c.Y))
	*/

	return &triangle
}

type TriangleGenerator struct {

	random JavaRandom
}

type JavaRandom struct {

	Seed int64
}

func (jr *JavaRandom) Next(bits int32) int32 {

	jr.Seed = (jr.Seed * 0x5DEECE66D + 0xB) & ((1 << 48) - 1)

	return int32(uint64(jr.Seed) >> uint64(48 - bits))
}

func (jr *JavaRandom) NextInt(bound int32) int32 {

	if (bound & -bound) == bound { // i.e., n is a power of 2
		return int32((int64(bound) * int64(jr.Next(31))) >> 31)
	}
	var bits int32
	var val = bound
	for bits - val + (bound - 1) < 0 {
		bits = jr.Next(31)
		val = bits % bound
	}
	return val
}

func (jr *JavaRandom) Intn(bound int) int {

	return int(jr.NextInt(int32(bound)))
}

func NewJavaRandom(seed int64) *JavaRandom {

	return &JavaRandom{seed ^ 0x5DEECE66D & ((1 << 48) - 1)}
}

func NewTriangleGenerator(seed int64) PrimitiveGenerator {

	return &TriangleGenerator{JavaRandom{seed}}
}

func (generator *TriangleGenerator) generate(width, height int) Primitive {

	var a = Point{generator.random.Intn(width), generator.random.Intn(height)}
	var b = Point{generator.random.Intn(width), generator.random.Intn(height)}
	var c = Point{generator.random.Intn(width), generator.random.Intn(height)}

	return NewTriangle(a, b, c)
}

func max(x int, y int) int {

	if x > y {
		return x
	}
	return y
}

func min(x int, y int) int {

	if x < y {
		return x
	}
	return y
}




