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

func generateTriangles(w, h int) []triangle {
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
		if w < h {
			r = float64(w / 2)
		}
		a := vec_add(vec_mul(v.A, vec{r, r}), vec{r, r})
		b := vec_add(vec_mul(v.B, vec{r, r}), vec{r, r})
		c := vec_add(vec_mul(v.C, vec{r, r}), vec{r, r})
		triangles[i] = triangle{Red: red, A: a, B: b, C: c}
	}

	return triangles
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
