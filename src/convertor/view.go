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
		getSocketPath func(*File, *widget.ProgressBar) (string, error),
	)
}

type View struct {
	w fyne.Window
}

type HandleConvertSetting struct {
	VideoFileInput *File
	SocketPath     string
}

type enableFormConversionStruct struct {
	fileVideoForConversion *widget.Button
	form                   *widget.Form
}

func NewView(w fyne.Window) *View {
	return &View{w}
}

func (v View) Main(
	runConvert func(setting HandleConvertSetting) error,
	getSocketPath func(*File, *widget.ProgressBar) (string, error),
) {
	form := &widget.Form{}

	conversionMessage := canvas.NewText("", color.RGBA{255, 0, 0, 255})
	conversionMessage.TextSize = 16
	conversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	progress := widget.NewProgressBar()
	progress.Hide()

	fileVideoForConversion, fileVideoForConversionMessage, fileInput := v.getButtonFileVideoForConversion(form, progress, conversionMessage)

	form.Items = []*widget.FormItem{
		{Text: "Файл для ковертации:", Widget: fileVideoForConversion},
		{Widget: fileVideoForConversionMessage},
	}
	form.SubmitText = "Конвертировать"

	enableFormConversionStruct := enableFormConversionStruct{
		fileVideoForConversion: fileVideoForConversion,
		form:                   form,
	}

	form.OnSubmit = func() {
		fileVideoForConversion.Disable()
		form.Disable()

		socketPath, err := getSocketPath(fileInput, progress)

		if err != nil {
			showConversionMessage(conversionMessage, err)
			enableFormConversion(enableFormConversionStruct)
			return
		}

		setting := HandleConvertSetting{
			VideoFileInput: fileInput,
			SocketPath:     socketPath,
		}
		err = runConvert(setting)
		if err != nil {
			showConversionMessage(conversionMessage, err)
			enableFormConversion(enableFormConversionStruct)
			return
		}
		enableFormConversion(enableFormConversionStruct)
	}

	v.w.SetContent(widget.NewCard("Конвертор видео файлов", "", container.NewVBox(form, conversionMessage, progress)))
	form.Disable()
}

func (v View) getButtonFileVideoForConversion(form *widget.Form, progress *widget.ProgressBar, conversionMessage *canvas.Text) (*widget.Button, *canvas.Text, *File) {
	fileInput := &File{}

	fileVideoForConversionMessage := canvas.NewText("", color.RGBA{255, 0, 0, 255})
	fileVideoForConversionMessage.TextSize = 16
	fileVideoForConversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	button := widget.NewButton("выбрать", func() {
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

				fileInput.Path = r.URI().Path()
				fileInput.Name = r.URI().Name()
				fileInput.Ext = r.URI().Extension()

				fileVideoForConversionMessage.Text = r.URI().Path()
				setStringSuccessStyle(fileVideoForConversionMessage)

				form.Enable()
				progress.Hide()
				conversionMessage.Text = ""
			}, v.w)
		fileDialog.Show()
	})

	return button, fileVideoForConversionMessage, fileInput
}

func setStringErrorStyle(text *canvas.Text) {
	text.Color = color.RGBA{255, 0, 0, 255}
	text.Refresh()
}

func setStringSuccessStyle(text *canvas.Text) {
	text.Color = color.RGBA{49, 127, 114, 255}
	text.Refresh()
}

func showConversionMessage(conversionMessage *canvas.Text, err error) {
	conversionMessage.Text = err.Error()
	setStringErrorStyle(conversionMessage)
}

func enableFormConversion(enableFormConversionStruct enableFormConversionStruct) {
	enableFormConversionStruct.fileVideoForConversion.Enable()
	enableFormConversionStruct.form.Enable()
}
