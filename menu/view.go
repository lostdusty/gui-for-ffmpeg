package menu

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"git.kor-elf.net/kor-elf/gui-for-ffmpeg/kernel"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/image/colornames"
	"net/url"
)

type ViewContract interface {
	About(ffmpegVersion string, ffprobeVersion string)
}

type View struct {
	app kernel.AppContract
}

func NewView(app kernel.AppContract) *View {
	return &View{
		app: app,
	}
}

func (v View) About(ffmpegVersion string, ffprobeVersion string) {
	view := v.app.GetAppFyne().NewWindow(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "about",
	}))
	view.Resize(fyne.Size{Width: 793, Height: 550})
	view.SetFixedSize(true)

	programmName := canvas.NewText(" GUI for FFmpeg", colornames.Darkgreen)
	programmName.TextStyle = fyne.TextStyle{Bold: true}
	programmName.TextSize = 20

	programmLink := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "git.kor-elf.net",
		Path:   "kor-elf/gui-for-ffmpeg/releases",
	})

	licenseLink := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "licenseLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "git.kor-elf.net",
		Path:   "kor-elf/gui-for-ffmpeg/src/branch/main/LICENSE",
	})

	licenseLinkOther := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "licenseLinkOther",
	}), &url.URL{
		Scheme: "https",
		Host:   "git.kor-elf.net",
		Path:   "kor-elf/gui-for-ffmpeg/src/branch/main/LICENSE-3RD-PARTY.txt",
	})

	programmVersion := widget.NewRichTextFromMarkdown(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmVersion",
		TemplateData: map[string]string{
			"Version": v.app.GetAppFyne().Metadata().Version,
		},
	}))

	aboutText := widget.NewRichText(
		&widget.TextSegment{
			Text: v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
				MessageID: "aboutText",
			}),
		},
	)
	image := canvas.NewImageFromFile("icon.png")
	image.SetMinSize(fyne.Size{Width: 100, Height: 100})
	image.FillMode = canvas.ImageFillContain

	ffmpegTrademark := widget.NewRichTextFromMarkdown(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "ffmpegTrademark",
	}))
	ffmpegLGPL := widget.NewRichTextFromMarkdown(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
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
				container.NewHBox(licenseLinkOther),
			)),
			v.getAboutFfmpeg(ffmpegVersion),
			v.getAboutFfprobe(ffprobeVersion),
			widget.NewCard(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
				MessageID: "AlsoUsedProgram",
			}), "", v.getOther()),
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

	programmLink := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "",
	})

	licenseLink := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
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

	programmLink := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
		MessageID: "programmLink",
	}), &url.URL{
		Scheme: "https",
		Host:   "ffmpeg.org",
		Path:   "ffprobe.html",
	})

	licenseLink := widget.NewHyperlink(v.app.GetLocalizerService().GetMessage(&i18n.LocalizeConfig{
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

func (v View) getOther() *fyne.Container {
	return container.NewVBox(
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("fyne.io/fyne/v2", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/fyne",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 3-Clause License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/fyne/blob/master/LICENSE",
		})),
		widget.NewRichTextFromMarkdown("Copyright (C) 2018 Fyne.io developers (see [AUTHORS](https://github.com/fyne-io/fyne/blob/master/AUTHORS))"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("fyne.io/systray", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/systray",
		})),
		container.NewHBox(widget.NewHyperlink("Apache License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/systray/blob/master/LICENSE",
		})),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/BurntSushi/toml", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "BurntSushi/toml",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "BurntSushi/toml/blob/master/COPYING",
		})),
		widget.NewLabel("Copyright (c) 2013 TOML authors"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/davecgh/go-spew", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "davecgh/go-spew",
		})),
		container.NewHBox(widget.NewHyperlink("ISC License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "davecgh/go-spew/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2012-2016 Dave Collins <dave@davec.name>"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/fredbi/uri", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fredbi/uri",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fredbi/uri/blob/master/LICENSE.md",
		})),
		widget.NewLabel("Copyright (c) 2018 Frederic Bidon"),
		widget.NewLabel("Copyright (c) 2015 Trey Tacon"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/fsnotify/fsnotify", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fsnotify/fsnotify",
		})),
		container.NewHBox(widget.NewHyperlink("BSD-3-Clause license", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fsnotify/fsnotify/blob/main/LICENSE",
		})),
		widget.NewLabel("Copyright © 2012 The Go Authors. All rights reserved."),
		widget.NewLabel("Copyright © fsnotify Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/fyne-io/gl-js", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/gl-js",
		})),
		container.NewHBox(widget.NewHyperlink("BSD-3-Clause license", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/gl-js/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/fyne-io/glfw-js", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/glfw-js",
		})),
		container.NewHBox(widget.NewHyperlink("MIT License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/glfw-js/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2014 Dmitri Shuralyov"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/fyne-io/image", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/image",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 3-Clause License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "fyne-io/image/blob/main/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2022, Fyne.io"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/go-gl/gl", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-gl/gl",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-gl/gl/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2014 Eric Woroshow"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/go-gl/glfw/v3.3/glfw", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-gl/glfw/",
		})),
		container.NewHBox(widget.NewHyperlink("BSD-3-Clause license", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-gl/glfw/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2012 The glfw3-go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/go-text/render", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-text/render",
		})),
		container.NewHBox(widget.NewHyperlink("Unlicense OR BSD-3-Clause", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-text/render/blob/main/LICENSE",
		})),
		widget.NewLabel("Copyright 2021 The go-text authors"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/go-text/typesetting", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-text/typesetting",
		})),
		container.NewHBox(widget.NewHyperlink("Unlicense OR BSD-3-Clause", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-text/typesetting/blob/main/LICENSE",
		})),
		widget.NewLabel("Copyright 2021 The go-text authors"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/godbus/dbus/v5", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "godbus/dbus",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 2-Clause \"Simplified\" License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "godbus/dbus/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2013, Georg Reinke (<guelfey at gmail dot com>), Google"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/gopherjs/gopherjs", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "gopherjs/gopherjs",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 2-Clause \"Simplified\" License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "gopherjs/gopherjs/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2013 Richard Musiol. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/jinzhu/inflection", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "jinzhu/inflection",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "jinzhu/inflection/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2015 - Jinzhu"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/jsummers/gobmp", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "jsummers/gobmp",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "jsummers/gobmp/blob/master/COPYING.txt",
		})),
		widget.NewLabel("Copyright (c) 2012-2015 Jason Summers"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/mattn/go-sqlite3", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "mattn/go-sqlite3",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "mattn/go-sqlite3/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2014 Yasuhiro Matsumoto"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/nicksnyder/go-i18n/v2", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "nicksnyder/go-i18n",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "nicksnyder/go-i18n/blob/main/LICENSE",
		})),
		widget.NewRichTextFromMarkdown("Copyright (c) 2014 Nick Snyder [https://github.com/nicksnyder](https://github.com/nicksnyder)"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/pmezard/go-difflib", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "pmezard/go-difflib",
		})),
		container.NewHBox(widget.NewHyperlink("License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "pmezard/go-difflib/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2013, Patrick Mezard"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/srwiley/oksvg", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "srwiley/oksvg",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 3-Clause License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "srwiley/oksvg/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2018, Steven R Wiley"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/srwiley/rasterx", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "srwiley/rasterx",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 3-Clause License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "srwiley/rasterx/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2018, Steven R Wiley"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/stretchr/testify", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "stretchr/testify",
		})),
		container.NewHBox(widget.NewHyperlink("MIT License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "stretchr/testify/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/tevino/abool", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "tevino/abool",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "tevino/abool/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2016 Tevin Zhang"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/yuin/goldmark", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "yuin/goldmark",
		})),
		container.NewHBox(widget.NewHyperlink("MIT License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "yuin/goldmark/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2019 Yusuke Inuzuka"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("golang.org/x/image", &url.URL{
			Scheme: "https",
			Host:   "pkg.go.dev",
			Path:   "golang.org/x/image",
		})),
		container.NewHBox(widget.NewHyperlink("License", &url.URL{
			Scheme: "https",
			Host:   "cs.opensource.google",
			Path:   "go/x/image/+/master:LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("golang.org/x/mobile", &url.URL{
			Scheme: "https",
			Host:   "pkg.go.dev",
			Path:   "golang.org/x/mobile",
		})),
		container.NewHBox(widget.NewHyperlink("License", &url.URL{
			Scheme: "https",
			Host:   "cs.opensource.google",
			Path:   "go/x/mobile/+/master:LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("golang.org/x/net", &url.URL{
			Scheme: "https",
			Host:   "pkg.go.dev",
			Path:   "golang.org/x/net",
		})),
		container.NewHBox(widget.NewHyperlink("License", &url.URL{
			Scheme: "https",
			Host:   "cs.opensource.google",
			Path:   "go/x/net/+/master:LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("golang.org/x/sys", &url.URL{
			Scheme: "https",
			Host:   "pkg.go.dev",
			Path:   "golang.org/x/sys",
		})),
		container.NewHBox(widget.NewHyperlink("License", &url.URL{
			Scheme: "https",
			Host:   "cs.opensource.google",
			Path:   "go/x/sys/+/master:LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("golang.org/x/text", &url.URL{
			Scheme: "https",
			Host:   "pkg.go.dev",
			Path:   "golang.org/x/text",
		})),
		container.NewHBox(widget.NewHyperlink("License", &url.URL{
			Scheme: "https",
			Host:   "cs.opensource.google",
			Path:   "go/x/text/+/master:LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("gopkg.in/yaml.v3", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-yaml/yaml/tree/v3.0.1",
		})),
		container.NewHBox(widget.NewHyperlink("Licensed under the Apache License, Version 2.0", &url.URL{
			Scheme: "http",
			Host:   "www.apache.org",
			Path:   "licenses/LICENSE-2.0",
		})),
		widget.NewLabel("Copyright 2011-2016 Canonical Ltd."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("gorm.io/gorm", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-gorm/gorm",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "go-gorm/gorm/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2013-NOW  Jinzhu <wosmvp@gmail.com>"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("honnef.co/go/js/dom", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "dominikh/go-js-dom",
		})),
		container.NewHBox(widget.NewHyperlink("The MIT License (MIT)", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "dominikh/go-js-dom/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2014 Dominik Honnef"),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/golang/go", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "golang/go",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 3-Clause \"New\" or \"Revised\" License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "golang/go/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),

		container.NewHBox(widget.NewHyperlink("github.com/golang/go", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "golang/go",
		})),
		container.NewHBox(widget.NewHyperlink("BSD 3-Clause \"New\" or \"Revised\" License", &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "golang/go/blob/master/LICENSE",
		})),
		widget.NewLabel("Copyright (c) 2009 The Go Authors. All rights reserved."),
		canvas.NewLine(colornames.Darkgreen),
	)
}
