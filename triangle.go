package main

import (
	"image"
	"math"
	"math/cmplx"

	"github.com/fogleman/gg"
)

type vec struct {
	X float64
	Y float64
}

type triangle struct {
	Red     bool
	A, B, C vec
}

func generateTriangleCircle(w, h int) []triangle {
	triangles := make([]triangle, 0)

	for i := 0; i < 10; i++ {
		b := cmplx.Rect(1, (2*float64(i)-1)*math.Pi/10)
		c := cmplx.Rect(1, (2*float64(i)+1)*math.Pi/10)
		if i%2 == 0 {
			b, c = c, b
		}
		triangles = append(triangles, triangle{Red: true, A: vec{0, 0}, B: vec{real(b), imag(b)}, C: vec{real(c), imag(c)}})
	}

	for i, v := range triangles {
		red := v.Red
		r := float64(h / 2)
		centerY := float64(0)
		centerX := (float64(w) - r*2) / 2
		if w < h {
			r = float64(w / 2)
			centerX = 0
			centerY = (float64(h) - r*2) / 2
		}
		a := vec_add(vec_mul(v.A, vec{r, r}), vec{r + centerX, r + centerY})
		b := vec_add(vec_mul(v.B, vec{r, r}), vec{r + centerX, r + centerY})
		c := vec_add(vec_mul(v.C, vec{r, r}), vec{r + centerX, r + centerY})
		triangles[i] = triangle{Red: red, A: a, B: b, C: c}
	}

	return triangles
}

const rad72 = 72 * (math.Pi / 180)

func generateRedTriangle(w, h int) triangle {
	A := vec{float64(w / 2), 0}

	hyp := float64(h) / math.Sin(rad72)
	adjacent := hyp * math.Cos(rad72)
	bottom := vec_add(A, vec{0, float64(h)})

	B := vec_sub(bottom, vec{adjacent, 0})
	C := vec_add(bottom, vec{adjacent, 0})

	return triangle{Red: true, A: A, B: B, C: C}
}

const rad36 = 36 * (math.Pi / 180)

func generateBlueTriangle(w, h int) triangle {

	adjacent := float64(w / 2)
	opp := math.Tan(rad36) * adjacent
	center := (float64(h) - opp) / 4

	A := vec{float64(w / 2), center}
	bottom := vec_add(A, vec{0, opp + center})

	B := vec_sub(bottom, vec{adjacent, 0})
	C := vec_add(bottom, vec{adjacent, 0})

	return triangle{Red: false, A: A, B: B, C: C}
}

var phi = vec{math.Phi, math.Phi}

func subdivide(n int, tri []triangle) []triangle {
	triangles := make([]triangle, len(tri))
	copy(triangles, tri)

	for i := 0; i < n; i++ {
		result := make([]triangle, 0)
		for _, t := range triangles {
			if t.Red {
				p := vec_add(t.A, vec_div(vec_sub(t.B, t.A), phi))
				result = append(result, triangle{Red: true, A: t.C, B: p, C: t.B}, triangle{Red: false, A: p, B: t.C, C: t.A})
			} else {
				q := vec_add(t.B, vec_div(vec_sub(t.A, t.B), phi))
				r := vec_add(t.B, vec_div(vec_sub(t.C, t.B), phi))
				result = append(result,
					triangle{Red: false, A: r, B: t.C, C: t.A},
					triangle{Red: false, A: q, B: r, C: t.B},
					triangle{Red: true, A: r, B: q, C: t.A},
				)
			}
		}
		triangles = result
	}
	return triangles
}

func drawTriangles(w, h int, triangles []triangle) image.Image {
	dc := gg.NewContext(w, h)
	for _, t := range triangles {
		drawTriangle(t, dc)
	}
	for _, t := range triangles {
		drawTriangleLines(t, dc)
	}
	return dc.Image()
}

func drawTriangle(t triangle, dc *gg.Context) {
	if t.Red {
		dc.SetRGB255(255, 0, 0)
	} else {
		dc.SetRGB255(0, 0, 255)
	}
	dc.MoveTo(t.A.X, t.A.Y)
	dc.LineTo(t.B.X, t.B.Y)
	dc.LineTo(t.C.X, t.C.Y)
	dc.LineTo(t.A.X, t.A.Y)
	dc.Fill()
	dc.SetLineWidth(1)
	dc.DrawLine(t.B.X, t.B.Y, t.C.X, t.C.Y)
	dc.Stroke()
}

func drawTriangleLines(t triangle, dc *gg.Context) {
	dc.SetRGB255(0, 0, 0)
	dc.SetLineWidth(2)
	dc.DrawLine(t.A.X, t.A.Y, t.B.X, t.B.Y)
	dc.Stroke()
	dc.DrawLine(t.A.X, t.A.Y, t.C.X, t.C.Y)
	dc.Stroke()
}
