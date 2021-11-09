package pages

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func instructionsElements() []app.UI {
	return []app.UI{
		app.P().Text(`Click on the blueprint book to select an image file. 
The book image will be replaced and a grayscale of the image will appear to the right.
To select a new image, click again on the first image in the row.`),
		app.P().Text(`Adjust the threshold slider and a preview image will appear 
to the right of the grayscale image. Black pixels will represent where the selected tile 
type(s) will appear in the blueprints. Use the invert checkbox if needed based on your image.
`),
		app.P().Text(`When you are satisfied with the preview image, use the size checkboxes and tile checkboxes 
to select the blueprints to generate. 
Click on the blueprint button and the blueprint book will be copied the clipboard.
`),
		app.Hr(),
	}
}
