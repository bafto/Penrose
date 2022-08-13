package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/lusingander/colorpicker"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetMaster()

	subdivisor := binding.NewFloat()
	subdivisor.Set(0)

	zoom := binding.NewFloat()
	zoom.Set(1)

	triType := "Thin-Triangle (1)"

	redrawOnly := false
	triangles := []triangle{generateThinTriangle(500, 500)}
	raster := canvas.NewRaster(func(w, h int) image.Image {
		sub, _ := subdivisor.Get()
		z, _ := zoom.Get()

		if redrawOnly {
			redrawOnly = false
			return drawTriangles(w, h, triangles, z)
		}

		switch triType {
		case "Thin-Triangle (1)":
			triangles = []triangle{generateThinTriangle(w, h)}
		case "Thin-Triangle (10)":
			triangles = generateTriangleCircle(w, h)
		case "Thick-Triangle (1)":
			triangles = []triangle{generateThickTriangle(w, h)}
		}

		triangles = subdivideTriangles(int(sub), triangles)
		return drawTriangles(w, h, triangles, z)
	})
	raster.SetMinSize(fyne.NewSize(500, 500))

	subdivSlider := widget.NewSliderWithData(0, 20, subdivisor)

	zoomSlider := widget.NewSliderWithData(0, 20, zoom)
	zoomSlider.Step = 0.1

	triTypeSelect := widget.NewSelect([]string{"Thin-Triangle (1)", "Thin-Triangle (10)", "Thick-Triangle (1)"}, func(s string) {
		triType = s
		raster.Refresh()
	})
	triTypeSelect.SetSelected(triType)

	thinColorPicker := colorpicker.New(200, colorpicker.StyleHue)
	thinColorPicker.SetColor(thinColor)
	thinColorPicker.SetOnChanged(func(c color.Color) { thinColor = c })

	thickColorPicker := colorpicker.New(200, colorpicker.StyleHue)
	thickColorPicker.SetColor(thickColor)
	thickColorPicker.SetOnChanged(func(c color.Color) { thickColor = c })

	bottom := container.NewGridWithColumns(4,
		newWithLabel(subdivSlider, "Density"),
		newWithLabel(triTypeSelect, "Base Shape"),
		newWithLabel(zoomSlider, "Zoom"),
		widget.NewButton("Thin-Rhombus color", func() {
			dialog.ShowCustomConfirm("Select Color", "Confirm", "Cancel", fyne.NewContainer(thinColorPicker), func(b bool) {}, w)
		}),
		widget.NewButton("Thick-Rhombus color", func() {
			dialog.ShowCustomConfirm("Select Color", "Confirm", "Cancel", fyne.NewContainer(thickColorPicker), func(b bool) {
				if b {
					raster.Refresh()
				}
			}, w)
		}),
		widget.NewButton("Random Colors", func() {
			thinColor, thickColor = rand_rgba(), rand_rgba()
			thinColorPicker.SetColor(thinColor)
			thickColorPicker.SetColor(thickColor)
			raster.Refresh()
		}),
		widget.NewButton("Draw", func() { raster.Refresh() }),
		widget.NewButton("Save As", func() {
			dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
				if uc == nil {
					return
				}
				if err != nil {
					dialog.NewError(err, w).Show()
					return
				}
				size := raster.Size()
				z, _ := zoom.Get()
				img := drawTriangles(int(size.Width), int(size.Height), triangles, z)
				switch ext := uc.URI().Extension(); ext {
				case ".png":
					err = png.Encode(uc, img)
				case ".jpeg", ".jpg":
					err = jpeg.Encode(uc, img, nil)
				case ".gif":
					err = gif.Encode(uc, img, nil)
				default:
					err = fmt.Errorf("Invalid file extenstion %s\nValid extensions are .png, .jpg, .jpeg or .gif", ext)
				}
				if err != nil {
					dialog.NewError(err, w).Show()
				}
			}, w).Show()
		}),
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

func rgba(r, g, b, a uint8) color.Color {
	return color.NRGBA{R: r, G: g, B: b, A: a}
}

func init() {
	rand.Seed(time.Now().Unix())
}

func rand_rgba() color.Color {
	r, g, b := rand.Intn(255), rand.Intn(255), rand.Intn(255)
	return color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}
