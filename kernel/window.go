package kernel

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/helper"
)

type WindowContract interface {
	SetContent(content fyne.CanvasObject)
	SetMainMenu(menu *fyne.MainMenu)
	NewFileOpen(callback func(fyne.URIReadCloser, error), location fyne.ListableURI) *dialog.FileDialog
	NewFolderOpen(callback func(fyne.ListableURI, error), location fyne.ListableURI) *dialog.FileDialog
	ShowAndRun()
	GetLayout() LayoutContract
}

type Window struct {
	windowFyne fyne.Window
	layout     LayoutContract
}

func newWindow(w fyne.Window, layout LayoutContract) Window {
	w.Resize(fyne.Size{Width: 799, Height: 599})
	w.CenterOnScreen()

	return Window{
		windowFyne: w,
		layout:     layout,
	}
}

func (w Window) SetContent(content fyne.CanvasObject) {
	w.windowFyne.SetContent(w.layout.SetContent(content))
}

func (w Window) NewFileOpen(callback func(fyne.URIReadCloser, error), location fyne.ListableURI) *dialog.FileDialog {
	fileDialog := dialog.NewFileOpen(callback, w.windowFyne)
	helper.FileDialogResize(fileDialog, w.windowFyne)
	fileDialog.Show()
	if location != nil {
		fileDialog.SetLocation(location)
	}
	return fileDialog
}

func (w Window) NewFolderOpen(callback func(fyne.ListableURI, error), location fyne.ListableURI) *dialog.FileDialog {
	fileDialog := dialog.NewFolderOpen(callback, w.windowFyne)
	helper.FileDialogResize(fileDialog, w.windowFyne)
	fileDialog.Show()
	if location != nil {
		fileDialog.SetLocation(location)
	}
	return fileDialog
}

func (w Window) SetMainMenu(menu *fyne.MainMenu) {
	w.windowFyne.SetMainMenu(menu)
}

func (w Window) ShowAndRun() {
	w.windowFyne.ShowAndRun()
}

func (w Window) GetLayout() LayoutContract {
	return w.layout
}
