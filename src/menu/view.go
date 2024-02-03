package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/localizer"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/image/colornames"
	"net/url"
)

type ViewContract interface {
	About(ffmpegVersion string, ffprobeVersion string)
}

type View struct {
	w                fyne.Window
	app              fyne.App
	appVersion       string
	localizerService localizer.ServiceContract
}

func NewView(w fyne.Window, app fyne.App, appVersion string, localizerService localizer.ServiceContract) *View {
	return &View{
		w:                w,
		app:              app,
		appVersion:       appVersion,
		localizerService: localizerService,
	}
}

func (v View) About(ffmpegVersion string, ffprobeVersion string) {
	view := v.app.NewWindow(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "about",
	}))
	view.Resize(fyne.Size{Width: 793, Height: 550})
	view.SetFixedSize(true)

	programmName := canvas.NewText(" GUI for FFmpeg", colornames.Darkgreen)
	programmName.TextStyle = fyne.TextStyle{Bold: true}
	programmName.TextSize = 20

	programmLink := widget.NewHyperlink(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "git.kor-elf.net",
		Path:   "kor-elf/gui-for-ffmpeg/releases",
	})

	licenseLink := widget.NewHyperlink(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "licenseLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "git.kor-elf.net",
		Path:   "kor-elf/gui-for-ffmpeg/src/branch/main/LICENSE",
	})

	programmVersion := widget.NewRichTextFromMarkdown(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmVersion",
		TemplateData: map[string]string{
			"Version": v.appVersion,
		},
	}))

	aboutText := widget.NewRichText(
		&widget.TextSegment{
			Text: v.localizerService.GetMessage(&i18n.LocalizeConfig{
				MessageID: "aboutText",
			}),
		},
	)
	image := canvas.NewImageFromFile("icon.png")
	image.SetMinSize(fyne.Size{Width: 100, Height: 100})
	image.FillMode = canvas.ImageFillContain

	ffmpegTrademark := widget.NewRichTextFromMarkdown(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "ffmpegTrademark",
	}))
	ffmpegLGPL := widget.NewRichTextFromMarkdown(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "ffmpegLGPL",
	}))

	view.SetContent(
		container.NewScroll(container.NewVBox(
			container.NewBorder(nil, nil, container.NewVBox(image), nil, container.NewVBox(
				programmName,
				programmVersion,
				aboutText,
				ffmpegTrademark,
				ffmpegLGPL,
				v.getCopyright(),
				container.NewHBox(programmLink, licenseLink),
			)),
			v.getAboutFfmpeg(ffmpegVersion),
			v.getAboutFfprobe(ffprobeVersion),
		)),
	)
	view.CenterOnScreen()
	view.Show()
}

func (v View) getCopyright() *widget.RichText {
	return widget.NewRichTextFromMarkdown("Copyright (c) 2024 **[Leonid Nikitin (kor-elf)](https://git.kor-elf.net/kor-elf/)**.")
}

func (v View) getAboutFfmpeg(version string) *fyne.Container {
	programmName := canvas.NewText(" FFmpeg", colornames.Darkgreen)
	programmName.TextStyle = fyne.TextStyle{Bold: true}
	programmName.TextSize = 20

	programmLink := widget.NewHyperlink(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "",
	})

	licenseLink := widget.NewHyperlink(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "licenseLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "legal.html",
	})

	return container.NewVBox(
		programmName,
		widget.NewLabel(version),
		widget.NewRichTextFromMarkdown("**FFmpeg** is a trademark of **[Fabrice Bellard](http://bellard.org/)**, originator of the **[FFmpeg](https://ffmpeg.org/about.html)** project."),
		widget.NewRichTextFromMarkdown("This software uses libraries from the **FFmpeg** project under the **[LGPLv2.1](http://www.gnu.org/licenses/old-licenses/lgpl-2.1.html)**."),
		container.NewHBox(programmLink, licenseLink),
	)
}

func (v View) getAboutFfprobe(version string) *fyne.Container {
	programmName := canvas.NewText(" FFprobe", colornames.Darkgreen)
	programmName.TextStyle = fyne.TextStyle{Bold: true}
	programmName.TextSize = 20

	programmLink := widget.NewHyperlink(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "ffprobe.html",
	})

	licenseLink := widget.NewHyperlink(v.localizerService.GetMessage(&i18n.LocalizeConfig{
		MessageID: "licenseLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "legal.html",
	})

	return container.NewVBox(
		programmName,
		widget.NewLabel(version),
		widget.NewRichTextFromMarkdown("**FFmpeg** is a trademark of **[Fabrice Bellard](http://bellard.org/)**, originator of the **[FFmpeg](https://ffmpeg.org/about.html)** project."),
		widget.NewRichTextFromMarkdown("This software uses libraries from the **FFmpeg** project under the **[LGPLv2.1](http://www.gnu.org/licenses/old-licenses/lgpl-2.1.html)**."),
		container.NewHBox(programmLink, licenseLink),
	)
}
