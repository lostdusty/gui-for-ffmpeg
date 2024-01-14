package convertor

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type ViewContract interface {
	Main(
		runConvert func(setting HandleConvertSetting) error,
		getSocketPath func(File, *widget.ProgressBar) (string, error),
	)
}

type View struct {
	w fyne.Window
}

type HandleConvertSetting struct {
	VideoFileInput File
	SocketPath     string
}

func NewView(w fyne.Window) *View {
	return &View{w}
}

func (v View) Main(
	runConvert func(setting HandleConvertSetting) error,
	getSocketPath func(File, *widget.ProgressBar) (string, error),
) {
	var fileInput File
	var form *widget.Form

	fileVideoForConversionMessage := canvas.NewText("", color.RGBA{255, 0, 0, 255})
	fileVideoForConversionMessage.TextSize = 16
	fileVideoForConversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	conversionMessage := canvas.NewText("", color.RGBA{255, 0, 0, 255})
	conversionMessage.TextSize = 16
	conversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	progress := widget.NewProgressBar()
	progress.Hide()

	fileVideoForConversion := widget.NewButton("выбрать", func() {
		fileDialog := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if err != nil {
					fileVideoForConversionMessage.Text = err.Error()
					setStringErrorStyle(fileVideoForConversionMessage)
					return
				}

				if r == nil {
					return
				}

				fileInput = File{
					Path: r.URI().Path(),
					Name: r.URI().Name(),
					Ext:  r.URI().Extension(),
				}
				fileVideoForConversionMessage.Text = r.URI().Path()
				setStringSuccessStyle(fileVideoForConversionMessage)

				form.Enable()
			}, v.w)
		fileDialog.Show()
	})

	form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Файл для ковертации:", Widget: fileVideoForConversion},
			{Widget: fileVideoForConversionMessage},
		},
		SubmitText: "Конвертировать",
		OnSubmit: func() {
			fileVideoForConversion.Disable()
			form.Disable()

			socketPath, err := getSocketPath(fileInput, progress)

			if err != nil {
				conversionMessage.Text = err.Error()
				setStringErrorStyle(conversionMessage)
				fileVideoForConversion.Enable()
				form.Enable()
			}

			setting := HandleConvertSetting{
				VideoFileInput: fileInput,
				SocketPath:     socketPath,
			}
			err = runConvert(setting)
			if err != nil {
				conversionMessage.Text = err.Error()
				setStringErrorStyle(conversionMessage)
			}
			fileVideoForConversion.Enable()
			form.Enable()
		},
	}

	v.w.SetContent(widget.NewCard("Конвертор видео файлов", "", container.NewVBox(form, conversionMessage, progress)))
	form.Disable()
}

func setStringErrorStyle(text *canvas.Text) {
	text.Color = color.RGBA{255, 0, 0, 255}
	text.Refresh()
}

func setStringSuccessStyle(text *canvas.Text) {
	text.Color = color.RGBA{49, 127, 114, 255}
	text.Refresh()
}
