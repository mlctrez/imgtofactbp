package clipboard

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/imgtofactbp/components"
)

// Clipboard allows copy and paste operations as go-app Actions
//
// ID is required
type Clipboard struct {
	app.Compo
	ID string
}

// PasteData is the type and data pasted into the clipboard
type PasteData struct {
	Type string
	Data string
}

func (pd *PasteData) String() string {
	sd := pd.Data
	if len(pd.Data) > 201 {
		sd = sd[:200] + " ..truncated.."
	}
	return fmt.Sprintf("&PasteData {Type:%q, Data:%q}", pd.Type, sd)
}

func (c *Clipboard) Render() app.UI {
	if c.ID == "" {
		panic("Clipboard must have an ID set")
	}
	return app.Div().ID(c.ID).Style("display", "none").Text("clipboard component")
}

func (c *Clipboard) WriteText(value string) {
	if cb, ok := (components.ValueHelper{Root: app.Window()}.Get("navigator", "clipboard")); ok {
		cb.Call("writeText", value)
	}
}

func (c *Clipboard) OnMount(ctx app.Context) {
	app.Window().AddEventListener("paste", c.pasteEventListener)
}

func (c *Clipboard) pasteEventListener(ctx app.Context, e app.Event) {
	log.Println("pasteEventListener")
	list, ok := components.ValueHelper{Root: e}.List("clipboardData", "items")
	if !ok {
		fmt.Println("paste with no list")
		return
	}
	for i, item := range list {
		log.Println("pasteEventListener item",i,item)
		clipboardItem, err := c.readPasteData(item)
		if err == ProtectedData {
			log.Println("pasteEventListener protected")
			continue
		}
		if err != nil {
			fmt.Println("error reading clipboard item", i, err)
			continue
		}
		ctx.NewActionWithValue("Clipboard:paste", clipboardItem, app.Tag{Name: "id", Value: c.ID})
	}
}

var ProtectedData = errors.New("protected data")

func (c *Clipboard) readPasteData(item app.Value) (result *PasteData, err error) {
	result = &PasteData{Type: item.Get("type").String()}
	switch item.Get("kind").String() {
	case "string":
		done := make(chan bool)
		item.Call("getAsString", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			result.Data = args[0].String()
			done <- true
			return nil
		}))
		<-done
	case "file":
		done := make(chan bool)
		reader := app.Window().Get("FileReader").New()
		reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
			done <- true
			return nil
		}))
		reader.Call("readAsDataURL", item.Call("getAsFile"))
		<-done
		if reader.Get("error").Truthy() {
			// TODO: extract real error from reader.Get("error")
			return nil, errors.New("error reading clipboard data")
		}
		result.Data = reader.Get("result").String()
	default:
		return nil, ProtectedData
	}
	return result, nil
}

// HandlePaste registers a callback for this clipboard when the paste event occurs
func (c *Clipboard) HandlePaste(ctx app.Context, prefix string, handler func(pasteData *PasteData)) {
	ctx.Handle("Clipboard:paste", func(context app.Context, action app.Action) {
		if data, ok := action.Value.(*PasteData); ok {
			if action.Tags.Get("id") != c.ID {
				return
			}
			if strings.HasPrefix(data.Type, prefix) {
				handler(data)
			}
		}
	})
}
