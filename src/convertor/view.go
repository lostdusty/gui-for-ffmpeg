package convertor

import (
	"errors"
	"ffmpegGui/helper"
	"ffmpegGui/localizer"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"image/color"
)

type ViewContract interface {
	Main(
		runConvert func(setting HandleConvertSetting, progressbar *widget.ProgressBar) error,
	)
	SelectFFPath(
		ffmpegPath string,
		ffprobePath string,
		save func(ffmpegPath string, ffprobePath string) error,
		cancel func(),
	)
}

type View struct {
	w                fyne.Window
	localizerService localizer.ServiceContract
}

type HandleConvertSetting struct {
	VideoFileInput       *File
	DirectoryForSave     string
	OverwriteOutputFiles bool
}

type enableFormConversionStruct struct {
	fileVideoForConversion *widget.Button
	buttonForSelectedDir   *widget.Button
	form                   *widget.Form
}

func NewView(w fyne.Window, localizerService localizer.ServiceContract) *View {
	return &View{
		w:                w,
		localizerService: localizerService,
	}
}

func (v View) Main(
	runConvert func(setting HandleConvertSetting, progressbar *widget.ProgressBar) error,
) {
	form := &widget.Form{}

	conversionMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	conversionMessage.TextSize = 16
	conversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	progress := widget.NewProgressBar()

	fileVideoForConversion, fileVideoForConversionMessage, fileInput := v.getButtonFileVideoForConversion(form, progress, conversionMessage)
	buttonForSelectedDir, buttonForSelectedDirMessage, pathToSaveDirectory := v.getButtonForSelectingDirectoryForSaving()

	isOverwriteOutputFiles := false
	checkboxOverwriteOutputFilesTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "checkboxOverwriteOutputFilesTitle",
	})
	checkboxOverwriteOutputFiles := widget.NewCheck(checkboxOverwriteOutputFilesTitle, func(b bool) {
		isOverwriteOutputFiles = b
	})

	form.Items = []*widget.FormItem{
		{
			Text:   v.localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: "fileVideoForConversionTitle"}),
			Widget: fileVideoForConversion,
		},
		{
			Widget: fileVideoForConversionMessage,
		},
		{
			Text:   v.localizerService.GetMessage(&i18n.LocalizeConfig{MessageID: "buttonForSelectedDirTitle"}),
			Widget: buttonForSelectedDir,
		},
		{
			Widget: buttonForSelectedDirMessage,
		},
		{
			Widget: checkboxOverwriteOutputFiles,
		},
	}
	form.SubmitText = v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "converterVideoFilesSubmitTitle",
	})

	enableFormConversionStruct := enableFormConversionStruct{
		fileVideoForConversion: fileVideoForConversion,
		buttonForSelectedDir:   buttonForSelectedDir,
		form:                   form,
	}

	form.OnSubmit = func() {
		if len(*pathToSaveDirectory) == 0 {
			showConversionMessage(conversionMessage, errors.New(v.localizerService.GetMessage(&i18n.LocalizeConfig{
				MessageID: "errorSelectedFolderSave",
			})))
			enableFormConversion(enableFormConversionStruct)
			return
		}
		conversionMessage.Text = ""

		fileVideoForConversion.Disable()
		buttonForSelectedDir.Disable()
		form.Disable()

		setting := HandleConvertSetting{
			VideoFileInput:       fileInput,
			DirectoryForSave:     *pathToSaveDirectory,
			OverwriteOutputFiles: isOverwriteOutputFiles,
		}
		err := runConvert(setting, progress)
		if err != nil {
			showConversionMessage(conversionMessage, err)
			enableFormConversion(enableFormConversionStruct)
			return
		}
		enableFormConversion(enableFormConversionStruct)
	}

	converterVideoFilesTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "converterVideoFilesTitle",
	})
	v.w.SetContent(widget.NewCard(converterVideoFilesTitle, "", container.NewVBox(form, conversionMessage, progress)))
	form.Disable()
}

func (v View) getButtonFileVideoForConversion(form *widget.Form, progress *widget.ProgressBar, conversionMessage *canvas.Text) (*widget.Button, *canvas.Text, *File) {
	fileInput := &File{}

	fileVideoForConversionMessage := canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	fileVideoForConversionMessage.TextSize = 16
	fileVideoForConversionMessage.TextStyle = fyne.TextStyle{Bold: true}

	buttonTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	button := widget.NewButton(buttonTitle, func() {
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
				progress.Value = 0
				progress.Refresh()
				conversionMessage.Text = ""
			}, v.w)
		helper.FileDialogResize(fileDialog, v.w)
		fileDialog.Show()
	})

	return button, fileVideoForConversionMessage, fileInput
}

func (v View) getButtonForSelectingDirectoryForSaving() (button *widget.Button, buttonMessage *canvas.Text, dirPath *string) {
	buttonMessage = canvas.NewText("", color.RGBA{R: 255, G: 0, B: 0, A: 255})
	buttonMessage.TextSize = 16
	buttonMessage.TextStyle = fyne.TextStyle{Bold: true}

	path := ""
	dirPath = &path

	buttonTitle := v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "choose",
	})

	button = widget.NewButton(buttonTitle, func() {
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
		helper.FileDialogResize(fileDialog, v.w)
		fileDialog.Show()
	})

	return
}

func setStringErrorStyle(text *canvas.Text) {
	text.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	text.Refresh()
}

func setStringSuccessStyle(text *canvas.Text) {
	text.Color = color.RGBA{R: 49, G: 127, B: 114, A: 255}
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
