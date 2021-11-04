package components

import (
	"errors"
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

// Clipboard allows copy and paste operations as go-app Actions
type Clipboard struct {
	app.Compo
	ID string
}

func (c *Clipboard) Render() app.UI {
	if c.ID == "" {
		panic("Clipboard must have an ID set")
	}
	return app.Div().ID(c.ID).Style("display", "none").Text("clipboard component")
}

func (c *Clipboard) OnDismount() {
}

func (c *Clipboard) pasteEventListener(ctx app.Context, e app.Event) {
	clipboardData := e.Get("clipboardData")
	if !clipboardData.Truthy() {
		return
	}
	items := clipboardData.Get("items")
	if !items.Truthy() {
		return
	}
	for i := 0; i < items.Length(); i++ {
		item := items.Index(i)
		itemType := item.Get("type").String()
		switch itemType {
		case "image/png", "image/jpg", "image/gif", "text/plain", "text/html":
			clipboardItem, err := c.readClipboardItem(itemType, item)
			if err != nil {
				fmt.Println("error reading clipboard item", i, itemType, err)
				continue
			}
			ctx.NewActionWithValue("Clipboard:paste", clipboardItem, app.Tag{Name: "id", Value: c.ID},
				app.Tag{Name: "type", Value: itemType})
		}
	}
}

type ClipboardPasteData struct {
	Type string
	Data string
}

// HandlePaste registers a callback for this clipboard when the paste event occurs
func (c *Clipboard) HandlePaste(ctx app.Context, prefix string, handler func(pasteData *ClipboardPasteData)) {
	ctx.Handle("Clipboard:paste", func(context app.Context, action app.Action) {
		itemType := action.Tags.Get("type")
		if data, ok := action.Value.(string); ok &&
			action.Tags.Get("id") == c.ID &&
			strings.HasPrefix(itemType, prefix) {
			handler(&ClipboardPasteData{Type: itemType, Data: data})
		}
	})
}

func (c *Clipboard) readClipboardItem(itemType string, item app.Value) (data string, err error) {

	if strings.HasPrefix(itemType, "text") {
		// TODO: figure out how to get text items
		//data = item.Call("getAsString", itemType).String()
		//fmt.Println(itemType, data)
		return
	}

	done := make(chan bool)
	reader := app.Window().Get("FileReader").New()
	reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		done <- true
		return nil
	}))
	if strings.HasPrefix(itemType, "image") {
		reader.Call("readAsDataURL", item.Call("getAsFile"))
	}
	<-done

	if reader.Get("error").Truthy() {
		// TODO: extract real error from reader.Get("error")
		return "", errors.New("error reading clipboard data")
	}

	return reader.Get("result").String(), nil
}

func (c *Clipboard) OnMount(ctx app.Context) {
	app.Window().AddEventListener("paste", c.pasteEventListener)
}

var _ app.Mounter = (*Clipboard)(nil)
var _ app.Dismounter = (*Clipboard)(nil)

/*

   document.addEventListener('paste', (event) => {
           let items = event.clipboardData.items;
           for (const itemKey in items) {
               const item = items[itemKey]
               if (item.type && item.type.indexOf("image") === 0) {
                   const reader = new FileReader();
                   reader.onload = function (event) {
                       fetch('/clips', {
                           method: "POST", headers: {'Content-Type': 'application/json'},
                           body: JSON.stringify({"clip": event.target.result.toString()})
                       }).then(showClips)
                           .catch((error) => {
                               console.error('Error:', error);
                           });
                   };
                   reader.readAsDataURL(item.getAsFile());
               }
           }
       }
   );

*/
