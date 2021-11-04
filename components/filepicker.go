package components

import (
	"errors"
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type FilePicker struct {
	app.Compo
	ID       string
	Multiple bool
	accept   string
}

type FilePickerResponse struct {
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

// Click calls click on the hidden file input using the element.
func (p *FilePicker) Click() {
	app.Window().GetElementByID(p.ID).Call("click")
}

func (p *FilePicker) readFile(file app.Value) (
	response *FilePickerResponse, err error) {

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

	response = &FilePickerResponse{
		Name: file.Get("name").String(),
		Size: file.Get("size").Int(),
		Type: file.Get("type").String(),
		Data: reader.Get("result").String(),
	}
	return
}

func (p *FilePicker) onChange(ctx app.Context, e app.Event) {
	target := e.Get("target")
	if !target.Truthy() {
		return
	}
	files := target.Get("files")
	if files.Truthy() && files.Length() > 0 {
		for i := 0; i < files.Length(); i++ {
			if file, err := p.readFile(files.Index(i)); err == nil {
				ctx.NewActionWithValue("FilePicker:file", file, app.Tag{Name: "id", Value: p.ID})
			} else {
				fmt.Println("onChange got error", err)
			}
		}
	}
}

// Handle registers the function callback for this file picker
func (p *FilePicker) Handle(ctx app.Context, handler func(file *FilePickerResponse)) {
	ctx.Handle("FilePicker:file", func(context app.Context, action app.Action) {
		if response, ok := action.Value.(*FilePickerResponse); ok && action.Tags.Get("id") == p.ID {
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
