package pages

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func instructionsElements() []app.UI {
	return []app.UI{
		app.P().Text(`Click on the left blueprint book to select an image to convert. The three images will
now change to your image, a grayscale representation, and a tile preview image. To select a new image, click on 
the first image in the row again.`),
		app.P().Text(`Adjust the threshold slider for the desired preview image. The initial preview image may 
not be visible until the slider is adjusted. Black pixels in the preview image will represent where the tiles will 
be placed in the blueprints. The invert checkbox will change if tiles are placed above or below the threshold.
`),
		app.P().Text(`When you are satisfied with the preview, use the tile and size checkboxes 
to select the blueprints to generate. Click on the blueprint button and a blueprint book will 
be copied into the clipboard, ready to import as a blueprint string.
`),
	}
}
