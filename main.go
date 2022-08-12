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

	triType := "Red (1)"

	var img image.Image
	raster := canvas.NewRaster(func(w, h int) image.Image {
		fmt.Println(w, h)
		sub, _ := subdivisor.Get()

		var triangles []triangle
		switch triType {
		case "Red (1)":
			triangles = []triangle{generateRedTriangle(w, h)}
			fmt.Println(triangles)
		case "Red (10)":
			triangles = generateTriangleCircle(w, h)
		case "Blue (1)":
			triangles = []triangle{generateBlueTriangle(w, h)}
		}

		triangles = subdivide(int(sub), triangles)
		img = drawTriangles(w, h, triangles)
		return img
	})
	raster.SetMinSize(fyne.NewSize(500, 500))

	subdivSlider := widget.NewSlider(1, 10)
	subdivSlider.Bind(subdivisor)

	triTypeSelect := widget.NewSelect([]string{"Red (1)", "Red (10)", "Blue (1)"}, func(s string) {
		triType = s
		raster.Refresh()
	})

	bottom := container.NewAdaptiveGrid(3,
		newWithLabel(subdivSlider, "Subdivcount"),
		newWithLabel(triTypeSelect, "Triangletype"),
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
