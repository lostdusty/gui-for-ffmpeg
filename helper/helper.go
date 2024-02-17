package helper

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func FileDialogResize(fileDialog *dialog.FileDialog, w fyne.Window) {
	contentSize := w.Content().Size()
	fileDialog.Resize(fyne.Size{Width: contentSize.Width - 50, Height: contentSize.Height - 50})
}
