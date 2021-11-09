package pages

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/factbp/blueprint"
	"github.com/mlctrez/goapp-mdc/pkg/bar"
	"github.com/mlctrez/goapp-mdc/pkg/base"
	"github.com/mlctrez/goapp-mdc/pkg/button"
	"github.com/mlctrez/goapp-mdc/pkg/checkbox"
	"github.com/mlctrez/goapp-mdc/pkg/icon"
	"github.com/mlctrez/goapp-mdc/pkg/slider"
	"github.com/mlctrez/imgtofactbp/components/clipboard"
	"github.com/mlctrez/imgtofactbp/components/filepicker"
	"github.com/mlctrez/imgtofactbp/conversions"
	"github.com/nfnt/resize"
)

const ImageRenderWidth = 300

type Index struct {
	app.Compo
	base.JsUtil
	picker         *filepicker.FilePicker
	clipboard      *clipboard.Clipboard
	original       image.Image
	scaled         image.Image
	grayscale      image.Image
	threshold      *slider.Continuous
	thresholdValue uint32
	inverted       bool
}

func (i *Index) OnNav(context app.Context) {
	p := context.Page()
	p.SetTitle("Image to Factorio Blueprint")
	p.SetAuthor("mlctrez")
	p.SetKeywords("factorio, blueprint, image, convert")
	p.SetDescription("A progressive web application for converting images to factorio tile blueprints.")
}

var _ app.Navigator = (*Index)(nil)

func (i *Index) OnMount(ctx app.Context) {
	ctx.Handle(string(slider.MDCSliderChange), func(context app.Context, action app.Action) {
		if action.Name == string(slider.MDCSliderChange) && action.Value == i.threshold {
			value, err := strconv.ParseFloat(action.Tags.Get("value"), 64)
			if err != nil {
				log.Println(err)
				return
			}
			i.thresholdValue = uint32(value)
			i.renderPreview()
		}
	})
	app.Window().GetElementByID("blueprint").Call("addEventListener", "click",
		app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			i.copyBlueprintBook(ctx, app.Event{})
			return nil
		}))
	i.picker.Handle(ctx, i.imageChanged)
	i.clipboard.HandlePaste(ctx, "image/", i.imagePaste)
}

func (i *Index) imageChanged(file *filepicker.Response) {
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
	i.renderPreview()
}

func (i *Index) imagePaste(data *clipboard.PasteData) {
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

func (i *Index) Render() app.UI {
	if i.picker == nil {
		i.picker = (&filepicker.FilePicker{ID: "hiddenFilePicker", Multiple: false}).Accept("image/*")
		i.clipboard = &clipboard.Clipboard{ID: "clipboard"}
		inputRange := &slider.InputRange{Id: "threshold", Name: "threshold", Label: "threshold", Min: 0, Max: 80000, Value: 40000, Step: 500}
		i.threshold = &slider.Continuous{Discrete: true, Id: "thresholdSlider", Range: inputRange}
	}
	topBar := &bar.TopAppBar{Title: "Image to Factorio Blueprint", Fixed: false}
	//topBar.Navigation = []app.HTMLButton{navButton()}
	topBar.Actions = []app.HTMLButton{githubButton()}
	body := app.Div().Body(&AppUpdateBanner{}, topBar, topBar.Main().Body(i.body()...))
	return body
}

func (i *Index) body() (body []app.UI) {
	body = append(body, i.picker, i.clipboard)
	body = append(body, instructionsElements()...)
	body = append(body, i.imagesRow())

	blueprintButton := &button.Button{Id: "blueprint", Icon: string(icon.MIConstruction),
		TrailingIcon: true, Raised: true, Label: "blueprint"}
	invertedCheckbox := &checkbox.Checkbox{Id: "invert", Label: "invert image", Callback: func(input app.HTMLInput) {
		input.OnChange(func(ctx app.Context, e app.Event) {
			i.inverted = ctx.JSSrc().Get("checked").Bool()
			i.renderPreview()
		})
	}}
	body = append(body, app.Div().Style("display", "inline").Body(
		i.threshold,
		app.Br(),
		app.Div().Style("display", "flex").Body(
			app.Div().Style("display", "inline-block").Body(i.tileCheckBoxes()...),
			app.Div().Style("display", "inline-block").Body(i.widthCheckBoxes()...),
		),
		invertedCheckbox, app.Raw("<span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</span>"), blueprintButton,
	))

	return body
}

func (i *Index) tileCheckBoxes() []app.UI {
	var content []app.UI
	for _, tile := range tileTypes() {
		content = append(content, &checkbox.Checkbox{Id: "tile-" + tile, Label: tile})
	}
	return content
}

func (i *Index) widthCheckBoxes() []app.UI {
	var content []app.UI
	for _, width := range blueprintWidths() {
		content = append(content, &checkbox.Checkbox{Id: fmt.Sprintf("width-%d", width), Label: fmt.Sprintf("%d", width)})
	}
	return content
}

func (i *Index) imagesRow() app.HTMLDiv {
	return app.Div().Style("display", "flex").Body(
		app.Img().ID("uploadedImage").Src("/web/logo-512.png").Width(ImageRenderWidth).
			Style("cursor", "pointer").OnClick(func(ctx app.Context, e app.Event) { i.picker.Open() }),
		app.Img().ID("grayScaleImage").Src("/web/logobw-512.png").Width(ImageRenderWidth).
			Style("cursor", "not-allowed"),
		app.Div().Class("col").Body(app.Img().ID("blueprintRender").
			Style("cursor", "not-allowed").Src("/web/logobw-512.png").Width(ImageRenderWidth)),
	)
}

func (i *Index) renderPreview() {
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
			if r > i.thresholdValue {
				preview.Set(x, y, onColor)
			} else {
				preview.Set(x, y, offColor)
			}
		}
	}
	setImageSrc("blueprintRender", conversions.ImageToBase64(preview))

}

func selectedWidths() (result []int) {
	widthsMap := make(map[int]bool)
	for _, width := range blueprintWidths() {
		id := fmt.Sprintf("width-%d-input", width)
		cb := app.Window().GetElementByID(id)
		if cb.Truthy() && cb.Get("checked").Truthy() {
			widthsMap[width] = true
		}
	}
	result = []int{}
	for w := range widthsMap {
		result = append(result, w)
	}
	sort.Ints(result)
	return
}

func selectedTiles() (result []string) {
	tilesMap := make(map[string]bool)
	tiles := tileTypes()
	for _, tile := range tiles {
		id := fmt.Sprintf("tile-%s-input", tile)
		cb := app.Window().GetElementByID(id)
		if cb.Truthy() && cb.Get("checked").Truthy() {
			tilesMap[tile] = true
		}
	}
	result = []string{}
	for t := range tilesMap {
		result = append(result, t)
	}
	sort.Strings(result)
	return
}

func sp(in string) *string {
	return &in
}

func (i *Index) copyBlueprintBook(ctx app.Context, e app.Event) {

	if i.grayscale == nil {
		return
	}

	widths := selectedWidths()
	tiles := selectedTiles()

	if len(widths) == 0 || len(tiles) == 0 {
		return
	}

	fmt.Println("widths", widths)
	fmt.Println("tiles", tiles)

	t := i.thresholdValue

	container := &blueprint.Container{}
	container.Book = &blueprint.Book{}
	container.Book.Label = sp("Image to Blueprint")

	var description string
	description += "Image to Blueprint book generated from https://mlctrez.github.io/imgtofactbp/\n\n"
	description += fmt.Sprintf("tiles : %s\n", strings.Join(tiles, ", "))
	var widthsStr []string
	for _, width := range widths {
		widthsStr = append(widthsStr, fmt.Sprintf("%d", width))
	}
	description += fmt.Sprintf("widths : %s\n", strings.Join(widthsStr, ", "))

	container.Book.Description = sp(description)

	for _, tileLoop := range tiles {
		tile := tileLoop

		tileBook := &blueprint.Book{}
		tileBook.Label = &tile
		fmt.Println("book for", tile)

		for _, widthLoop := range widths {
			width := widthLoop
			resized := conversions.ResizeWidth(i.grayscale, width)
			bp := buildBlueprint(fmt.Sprintf("width-%d", width), resized, func(r, g, b, a uint32) string {
				if i.inverted {
					if r > t {
						return tile
					}
					return ""
				} else {
					if r > t {
						return ""
					}
					return tile
				}
			})
			tileBook.AddBlueprint(bp)
		}
		container.Book.AddBook(tileBook)
	}

	output := &bytes.Buffer{}
	container.Write(output)
	i.clipboard.WriteText(output.String())
}

func buildBlueprint(label string, img image.Image, tileAt func(r, g, b, a uint32) string) *blueprint.Blueprint {
	blue := &blueprint.Blueprint{}
	blue.Label = &label
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

func navButton() app.HTMLButton {
	return app.Button().Class(icon.MaterialIconsClass, icon.MaterialIconButton).
		Body(app.Img().Src("/web/logobw-512.png"))
}

func githubButton() app.HTMLButton {
	return app.Button().Title("show me the code!").
		Class(icon.MaterialIconsClass, icon.MaterialIconButton).
		Body(app.Raw(GitHubSvg)).OnClick(func(ctx app.Context, e app.Event) {
		app.Window().Call("open", "https://github.com/mlctrez/imgtofactbp/")
	})
}

const GitHubSvg = `<svg class="mdc-button__icon" width="32" height="32" aria-hidden="true" viewBox="0 0 16 16">
    <path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path>
</svg>`

func tileTypes() []string {
	return []string{
		"landfill", "stone-path", "concrete", "refined-concrete",
		"hazard-concrete-left", "hazard-concrete-right",
		"refined-hazard-concrete-left", "refined-hazard-concrete-right"}
}

func blueprintWidths() (result []int) {
	for i := 20; i <= 400; i = i + 20 {
		result = append(result, i)
	}
	return
}
