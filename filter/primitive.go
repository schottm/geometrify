package filter

type Primitive interface {

	//IsInside(x, y int) bool

	//Bounds() *BoundingBox
	ForEach(do func(int, int))

	GetColor() *Color

	SetColor(color Color)
}

type Point struct {

	X, Y int
}

type Color struct {

	R, G, B, A uint16
}


type BoundingBox struct {

	Min, Max Point
}

type PrimitiveGenerator interface {

	generate(width, height int) Primitive
}
