package filepicker

import (
	"errors"
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/imgtofactbp/components"
)

// FilePicker hides an input of type file and allows Actions for
type FilePicker struct {
	app.Compo
	ID       string
	Multiple bool
	accept   string
}

type Response struct {
	Name string
	Size int
	Type string
	Data string
}

// Accept sets the accept attribute.
//
// The open file dialog still allows picking any type so
// the response should be verified.
func (p *FilePicker) Accept(types string) *FilePicker {
	p.accept = types
	return p
}

// Open triggers the choose file dialog.
func (p *FilePicker) Open() {
	app.Window().GetElementByID(p.ID).Call("click")
}

func (p *FilePicker) readFile(file app.Value) (
	response *Response, err error) {

	done := make(chan bool)
	reader := app.Window().Get("FileReader").New()
	reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		done <- true
		return nil
	}))
	reader.Call("readAsDataURL", file)
	<-done

	if reader.Get("error").Truthy() {
		// TODO: extract real error from reader.Get("error")
		return nil, errors.New("error reading file")
	}

	response = &Response{
		Name: file.Get("name").String(),
		Size: file.Get("size").Int(),
		Type: file.Get("type").String(),
		Data: reader.Get("result").String(),
	}
	return
}

func (p *FilePicker) onChange(ctx app.Context, e app.Event) {

	list, ok := components.ValueHelper{Root: e}.List("target", "files")
	if !ok {
		return
	}
	for _, fileValue := range list {
		if file, err := p.readFile(fileValue); err == nil {
			ctx.NewActionWithValue("FilePicker:file", file, app.Tag{Name: "id", Value: p.ID})
		} else {
			fmt.Println("onChange got error", err)
		}
	}
}

// Handle registers the function callback for this file picker
func (p *FilePicker) Handle(ctx app.Context, handler func(file *Response)) {
	ctx.Handle("FilePicker:file", func(context app.Context, action app.Action) {
		if response, ok := action.Value.(*Response); ok && action.Tags.Get("id") == p.ID {
			handler(response)
		}
	})
}

func (p *FilePicker) Render() app.UI {
	if p.ID == "" {
		panic("filePicker must have an ID attribute")
	}
	fileInput := app.Input().Type("file").ID(p.ID)
	fileInput.Style("display", "none").Multiple(p.Multiple)
	if p.accept != "" {
		fileInput.Accept(p.accept)
	}
	fileInput.OnChange(p.onChange)
	return fileInput
}
