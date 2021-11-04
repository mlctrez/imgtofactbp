package components

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"strconv"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/edgeefy"
	"github.com/mlctrez/factbp/blueprint"
	"github.com/mlctrez/imgtofactbp/conversions"
	"github.com/nfnt/resize"
)

type Index struct {
	app.Compo
	picker    *FilePicker
	clipboard *Clipboard
	original  image.Image
	scaled    image.Image
	grayscale image.Image
	threshold uint32
	inverted  bool
	tileType  string
}

func (i *Index) OnMount(ctx app.Context) {
	i.picker.Handle(ctx, i.imageChanged)
	i.clipboard.HandlePaste(ctx, "image/", i.imagePaste)
}

var _ app.Mounter = (*Index)(nil)

const ImageRenderWidth = 300

func (i *Index) imageChanged(file *FilePickerResponse) {
	img, _, err := conversions.Base64ToImage(file.Data)
	if err != nil {
		fmt.Println("error decoding image", err)
		return
	}
	i.original = img
	i.renderImages()
}

func (i *Index) renderImages() {
	// normalize image width to 400px
	i.scaled = resize.Resize(ImageRenderWidth, 0, i.original, resize.Lanczos3)
	setImageSrc("uploadedImage", conversions.ImageToBase64(i.scaled))
	grayScale, _ := conversions.GrayScale(i.scaled)
	i.grayscale = grayScale
	setImageSrc("grayScaleImage", conversions.ImageToBase64(grayScale))

}

func (i *Index) imagePaste(data *ClipboardPasteData) {
	pastedImage, _, err := conversions.Base64ToImage(data.Data)
	if err != nil {
		fmt.Println(err)
		return
	}
	i.original = pastedImage
	i.renderImages()
}

func setImageSrc(id string, src string) {
	app.Window().GetElementByID(id).Set("src", src)
}

func (i *Index) resizeClicked(ctx app.Context, e app.Event) {
	if i.original == nil {
		fmt.Println("blueprint width or image not set")
		return
	}
	widthString := app.Window().GetElementByID("blueprintWidth").Get("value").String()
	width, err := strconv.Atoi(widthString)
	if err != nil {
		fmt.Println("width is not an int")
		return
	}

	pixels, err := edgeefy.GrayPixelsFrommImage(i.original)
	if err != nil {
		panic(err)
	}
	//pixels = edgeefy.CannyEdgeDetect(pixels, false, .5, .1)
	grey := edgeefy.GrayImageFromGrayPixels(pixels)

	img := resize.Resize(uint(width), 0, grey, resize.Lanczos3)
	w := &bytes.Buffer{}
	w.WriteString("data:image/png;base64,")
	encoder := base64.NewEncoder(base64.StdEncoding, w)
	err = png.Encode(encoder, img)
	_ = encoder.Close()
	if err != nil {
		fmt.Println("error resizing image")
		return
	}
	app.Window().GetElementByID("resizedImage").Set("src", w.String())
}

func tileTypes() []string {
	return []string{"landfill", "stone-path", "concrete", "refined-concrete",
		"hazard-concrete-left", "refined-hazard-concrete-left"}
}

func (i *Index) Render() app.UI {
	fmt.Println("render of index page")
	i.picker = (&FilePicker{ID: "hiddenFilePicker", Multiple: false}).Accept("image/*")
	i.clipboard = &Clipboard{ID: "clipboard"}
	if i.threshold == 0 {
		fmt.Println("setting i.threshold = 40000")
		i.threshold = 40000
		i.inverted = false
		i.tileType = tileTypes()[0]
	}

	body := app.Div().Class("container").Body(
		navbar(),
		i.imagesRow(),
		app.Div().Class("row row-cols-auto").Body(
			app.Div().Class("col").Body(
				app.Hr(), app.P().Text(instructions),
			),
		),
		app.Div().Class("row row-cols-auto").Body(
			app.Div().Class("col-4").Body(
				app.Label().For("threshold").Class("form-label").Text("threshold"),
				app.Input().ID("threshold").Type("range").Class("form-range").
					Min("0").Max(80000).Step(1000).Value(fmt.Sprintf("%d", i.threshold)).
					OnChange(func(ctx app.Context, e app.Event) {
						parseInt, err := strconv.ParseInt(ctx.JSSrc().Get("value").String(), 10, 32)
						if err == nil {
							i.threshold = uint32(parseInt)
							i.renderPreview(uint32(parseInt))
						}
					}),
				app.Select().Class("form-select").
					ID("tileType").Body(func() []app.UI {
					var opts []app.UI
					for i, s := range tileTypes() {
						opts = append(opts, app.Option().ID(fmt.Sprintf("tile_option_%d", i)).Value(s).Text(s))
					}
					return opts
				}()...).OnChange(func(ctx app.Context, e app.Event) {
					i.tileType = ctx.JSSrc().Get("value").String()
					fmt.Println("set tile type to", i.tileType)
				}),
				app.Div().Class("form-check").Body(
					app.Input().Class("form-check-input").Type("checkbox").ID("inverted").
						OnChange(i.invertChange).Checked(i.inverted),
					app.Label().Class("form-check-label").For("inverted").Text("Invert"),
				),
				app.Hr(),
				app.Div().Class("row row-cols-auto").Body(
					app.Div().Class("col").Body(sizeCb(30), sizeCb(60), sizeCb(80)),
					app.Div().Class("col").Body(sizeCb(100), sizeCb(120), sizeCb(150)),
					app.Div().Class("col").Body(
						app.Button().Class("btn btn-success").Text("Copy Blueprint Book").
							OnClick(i.copyBlueprintBook),
					),
				),
				/*
					<label for="exampleFormControlTextarea1" class="form-label">Example textarea</label>
					  <textarea class="form-control" id="exampleFormControlTextarea1" rows="3"></textarea>
				*/
				app.Div().Class("row row-cols-auto").Body(
					app.Label().For("blueprintText").Class("form-label").Text("blueprint string"),
					app.Textarea().Class("form-control").ID("blueprintText").Rows(20),
				),
			),
			app.Div().Class("col-4"),
		),
	)

	return body
}

func sizeCb(size int) app.HTMLDiv {
	id := fmt.Sprintf("size-%d", size)
	text := fmt.Sprintf("Width %d", size)
	checked := size < 80
	return app.Div().Class("form-check").Body(
		app.Input().Class("form-check-input").Type("checkbox").ID(id).Checked(checked),
		app.Label().Class("form-check-label").For(id).Text(text),
	)
}

func (i *Index) invertChange(ctx app.Context, e app.Event) {
	i.inverted = ctx.JSSrc().Get("checked").Bool()
	i.renderPreview(i.threshold)
}

func navbar() app.HTMLDiv {
	return app.Div().Class("container").Body(
		app.Nav().Class("navbar navbar-light bg-light").Body(
			app.Div().Class("container-fluid").Body(
				app.Span().Class("navbar-text").Text("Image to blueprint converter v0.1"),
			),
		),
	)
}

func (i *Index) imagesRow() app.HTMLDiv {
	return app.Div().Class("row row-cols-auto").Body(
		app.Div().Class("col").Body(
			i.picker,
			i.clipboard,
			app.Img().ID("uploadedImage").Src("/web/logo-192.png").
				Width(ImageRenderWidth).Style("cursor", "pointer").
				OnClick(func(ctx app.Context, e app.Event) { i.picker.Click() }),
		),
		app.Div().Class("col").Body(
			app.Img().ID("grayScaleImage").Src("").
				Width(ImageRenderWidth).Style("cursor", "not-allowed"),
		),
		app.Div().Class("col").Body(
			app.Img().ID("blueprintRender").Src("").Width(ImageRenderWidth),
		),
	)
}

func (i *Index) renderPreview(t uint32) {
	if i.grayscale == nil {
		return
	}

	img := i.grayscale
	preview := image.NewGray(i.grayscale.Bounds())
	onColor := color.White
	offColor := color.Black
	if i.inverted {
		onColor = color.Black
		offColor = color.White
	}
	for x := 0; x < img.Bounds().Max.X; x = x + 1 {
		for y := 0; y < img.Bounds().Max.Y; y = y + 1 {
			r, _, _, _ := img.At(x, y).RGBA()
			if r > t {
				preview.Set(x, y, onColor)
			} else {
				preview.Set(x, y, offColor)
			}
		}
	}
	setImageSrc("blueprintRender", conversions.ImageToBase64(preview))

}

func (i *Index) copyBlueprintBook(ctx app.Context, e app.Event) {

	sizeMap := make(map[int]bool)
	sizes := []int{30, 60, 80, 100, 120, 150}
	for _, size := range sizes {
		checkbox := app.Window().GetElementByID(fmt.Sprintf("size-%d", size))
		if checkbox.Truthy() {
			sizeMap[size] = checkbox.Get("checked").Bool()
		}
	}
	fmt.Println("copyBlueprintBook i.tileType=", i.tileType)
	book := &blueprint.GenericForm{}
	for _, size := range sizes {
		if !sizeMap[size] {
			continue
		}
		t := i.threshold
		width := conversions.ResizeWidth(i.grayscale, size)
		book.AddBlueprint(buildBlueprint(fmt.Sprintf("size-%d", size), width, func(r, g, b, a uint32) string {
			if i.inverted {
				if r > t {
					return i.tileType
				}
				return ""
			} else {
				if r > t {
					return ""
				}
				return i.tileType
			}
		}))
	}
	output := &bytes.Buffer{}
	book.Write(output)
	app.Window().GetElementByID("blueprintText").Set("value", output.String())
}

func buildBlueprint(label string, img image.Image, tileAt func(r, g, b, a uint32) string) *blueprint.Blueprint {
	blue := &blueprint.Blueprint{Label: &label}
	for x := 0; x < img.Bounds().Max.X; x = x + 1 {
		for y := 0; y < img.Bounds().Max.Y; y = y + 1 {
			name := tileAt(img.At(x, y).RGBA())
			if name != "" {
				blue.AddTile(blueprint.TileWithPosition(name, float64(x), float64(y)))
			}
		}
	}
	return blue
}

var instructions = `

Click on the blueprint book above and select a image file to convert to a tile blueprint.
This image will replace the blueprint book and a grayscale version will appear immediately to the right.

Adjust the threshold slider and a preview image will appear to the right of the grayscale image. Black pixels will
represent where the selected tile type will appear in the blueprint. Use the invert checkbox to flip black and
white tiles in the resulting blueprint.

When you're satisfied with the preview image, click on the copy button and your clipboard will now contain
the blueprint book string to paste into Factorio. Use the size checkboxes to determine which sizes of the
resulting blueprint are placed in the book.

`
