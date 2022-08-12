package main

import (
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	subdivisor := binding.NewFloat()
	subdivisor.Set(1)

	numtri := binding.NewFloat()
	numtri.Set(1)

	var img image.Image
	raster := canvas.NewRaster(func(w, h int) image.Image {
		fmt.Println(w, h)
		sub, _ := subdivisor.Get()
		n, _ := numtri.Get()
		triangles := generateTriangles(int(n), w, h)
		triangles = subdivide(int(sub), triangles)
		img = drawTriangles(w, h, triangles)
		return img
	})
	raster.SetMinSize(fyne.NewSize(500, 500))

	subdivSlider := widget.NewSlider(1, 10)
	subdivSlider.Bind(subdivisor)

	numtriSlider := widget.NewSlider(1, 20)
	numtriSlider.Bind(numtri)

	bottom := container.NewAdaptiveGrid(3,
		newWithLabel(subdivSlider, "Subdivcount"),
		newWithLabel(numtriSlider, "Trianglecount"),
		widget.NewButton("draw", func() { raster.Refresh() }),
	)

	content := container.NewBorder(nil, bottom, nil, nil,
		container.NewMax(raster),
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func newWithLabel(w fyne.Widget, label string) *fyne.Container {
	return container.NewVBox(w, widget.NewLabel(label))
}
