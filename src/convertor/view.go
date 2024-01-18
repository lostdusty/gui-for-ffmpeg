package convertor

import (
	"errors"
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
	VideoFileInput       *File
	DirectoryForSave     string
	SocketPath           string
	OverwriteOutputFiles bool
}

type enableFormConversionStruct struct {
	fileVideoForConversion *widget.Button
	buttonForSelectedDir   *widget.Button
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

	fileVideoForConversion, fileVideoForConversionMessage, fileInput := v.getButtonFileVideoForConversion(form, progress, conversionMessage)
	buttonForSelectedDir, buttonForSelectedDirMessage, pathToSaveDirectory := v.getButtonForSelectingDirectoryForSaving()

	isOverwriteOutputFiles := false
	checkboxOverwriteOutputFiles := widget.NewCheck("Разрешить перезаписать файл", func(b bool) {
		isOverwriteOutputFiles = b
	})

	form.Items = []*widget.FormItem{
		{Text: "Файл для ковертации:", Widget: fileVideoForConversion},
		{Widget: fileVideoForConversionMessage},
		{Text: "Папка куда будет сохранятся:", Widget: buttonForSelectedDir},
		{Widget: buttonForSelectedDirMessage},
		{Widget: checkboxOverwriteOutputFiles},
	}
	form.SubmitText = "Конвертировать"

	enableFormConversionStruct := enableFormConversionStruct{
		fileVideoForConversion: fileVideoForConversion,
		buttonForSelectedDir:   buttonForSelectedDir,
		form:                   form,
	}

	form.OnSubmit = func() {
		if len(*pathToSaveDirectory) == 0 {
			showConversionMessage(conversionMessage, errors.New("Не выбрали папку для сохранения!"))
			enableFormConversion(enableFormConversionStruct)
			return
		}
		conversionMessage.Text = ""

		fileVideoForConversion.Disable()
		buttonForSelectedDir.Disable()
		form.Disable()

		socketPath, err := getSocketPath(fileInput, progress)

		if err != nil {
			showConversionMessage(conversionMessage, err)
			enableFormConversion(enableFormConversionStruct)
			return
		}

		setting := HandleConvertSetting{
			VideoFileInput:       fileInput,
			DirectoryForSave:     *pathToSaveDirectory,
			SocketPath:           socketPath,
			OverwriteOutputFiles: isOverwriteOutputFiles,
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

func (v View) getButtonForSelectingDirectoryForSaving() (button *widget.Button, buttonMessage *canvas.Text, dirPath *string) {
	buttonMessage = canvas.NewText("", color.RGBA{255, 0, 0, 255})
	buttonMessage.TextSize = 16
	buttonMessage.TextStyle = fyne.TextStyle{Bold: true}

	path := ""
	dirPath = &path

	button = widget.NewButton("выбрать", func() {
		fileDialog := dialog.NewFolderOpen(
			func(r fyne.ListableURI, err error) {
				if err != nil {
					buttonMessage.Text = err.Error()
					setStringErrorStyle(buttonMessage)
					return
				}
				if r == nil {
					return
				}

				path = r.Path()

				buttonMessage.Text = r.Path()
				setStringSuccessStyle(buttonMessage)
			}, v.w)
		fileDialog.Show()
	})

	return
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
	enableFormConversionStruct.buttonForSelectedDir.Enable()
	enableFormConversionStruct.form.Enable()
}
