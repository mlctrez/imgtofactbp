package pages

import "github.com/maxence-charriere/go-app/v9/pkg/app"

func instructionsElements() []app.UI {
	return []app.UI{
		app.P().Style("font-size","14px").Text(`Click on the left blueprint book to select an image to 
convert. The three images will now change to your image, a grayscale representation, and a tile preview image. 
To select a new image, click on the first image in the row again.`),
		app.P().Style("font-size","14px").Text(`Adjust the threshold slider for the desired preview image. 
The preview image may only show as one color until the slider is adjusted. Black pixels in the preview image 
will represent where tiles will be placed in the blueprint.  Use the invert checkbox if the preview 
is opposite of what you want.`),
		app.P().Style("font-size","14px").Text(`When you are satisfied with the preview, use the 
tile and size checkboxes to select the blueprints to generate. Clicking on the blueprint button will copy a 
blueprint book into the clipboard, ready to import as a blueprint string.
`),
	}
}
